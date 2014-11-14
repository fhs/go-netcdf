// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Some design decisions:
//
// The Get* methods below do their own allocations unlike the Var Get*
// methods because attributes are expected to be much smaller. Thus, we
// free the user from the burden of finding the length of the attribute
// and allocating the buffer for it.
//
// Also, we return []byte (i.e. []uint8) for NC_CHAR Type because:
//	- Returning string would not be very flexible, since '\0' characters
//	  may or may not require trimming.
//	- Returning []rune (i.e. []int32) takes up more space and we know
//	  we're limited to ASCII.
//	- Any other types can't be easily converted to a string
//	  (e.g. string([]int8) does not work)

package netcdf

// #include <stdlib.h>
// #include <netcdf.h>
import "C"

import (
	"unsafe"
)

// Attr represents an attribute associated with a variable.
type Attr struct {
	v    Var
	name string
}

// Type returns the data type of attribute a.
func (a Attr) Type() (t Type, err error) {
	// TODO: convert a.name to CString only once instead of in
	// each method of a.
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	var ct C.nc_type
	err = newError(C.nc_inq_atttype(C.int(a.v.f), C.int(a.v.id), cname, &ct))
	t = Type(ct)
	return
}

// Len returns the length of the attribute value.
func (a Attr) Len() (n uint64, err error) {
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	var cn C.size_t
	err = newError(C.nc_inq_attlen(C.int(a.v.f), C.int(a.v.id), cname, &cn))
	n = uint64(cn)
	return
}

// PutChar sets the value of attribute a to val.
func (a Attr) PutChar(val []byte) error {
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	return newError(C.nc_put_att_text(C.int(a.v.f), C.int(a.v.id), cname,
		C.size_t(len(val)), (*C.char)(unsafe.Pointer(&val[0]))))
}

// GetChar returns the attribute value.
func (a Attr) GetChar() (val []byte, err error) {
	n, err := a.Len()
	if err != nil {
		return nil, err
	}
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	val = make([]byte, n)
	err = newError(C.nc_get_att_text(C.int(a.v.f), C.int(a.v.id), cname,
		(*C.char)(unsafe.Pointer(&val[0]))))
	return
}

// Attr returns attribute named name.
func (v Var) Attr(name string) (a Attr) {
	return Attr{v: v, name: name}
}

// AttrN returns attribute for attribute number n.
func (v Var) AttrN(n int) (a Attr, err error) {
	buf := C.CString(string(make([]byte, _NC_MAX_NAME+1)))
	defer C.free(unsafe.Pointer(buf))
	err = newError(C.nc_inq_attname(C.int(v.f), C.int(v.id), C.int(n), buf))
	a = Attr{v: v, name: C.GoString(buf)}
	return
}

// Attr returns attribute named name.
func (f File) Attr(name string) (a Attr) {
	return Var{f, _NC_GLOBAL}.Attr(name)
}

// AttrN returns attribute for attribute number n.
func (f File) AttrN(n int) (a Attr, err error) {
	return Var{f, _NC_GLOBAL}.AttrN(n)
}
