// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

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

// Name returns the name of attribute a.
func (a Attr) Name() string {
	return a.name
}

// Type returns the data type of attribute a.
func (a Attr) Type() (t Type, err error) {
	// TODO: convert a.name to CString only once instead of in
	// each method of a.
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	var ct C.nc_type
	err = newError(C.nc_inq_atttype(C.int(a.v.ds), C.int(a.v.id), cname, &ct))
	t = Type(ct)
	return
}

// Len returns the length of the attribute value.
func (a Attr) Len() (n uint64, err error) {
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	var cn C.size_t
	err = newError(C.nc_inq_attlen(C.int(a.v.ds), C.int(a.v.id), cname, &cn))
	n = uint64(cn)
	return
}

// Attr returns attribute named name. If the attribute does not yet exist,
// it'll be created once it's written.
func (v Var) Attr(name string) (a Attr) {
	return Attr{v: v, name: name}
}

// AttrN returns attribute for attribute number n.
func (v Var) AttrN(n int) (a Attr, err error) {
	buf := C.CString(string(make([]byte, C.NC_MAX_NAME+1)))
	defer C.free(unsafe.Pointer(buf))
	err = newError(C.nc_inq_attname(C.int(v.ds), C.int(v.id), C.int(n), buf))
	a = Attr{v: v, name: C.GoString(buf)}
	return
}

// Attr returns global attribute named name. If the attribute does not yet
// exist, it'll be created once it's written.
func (ds Dataset) Attr(name string) (a Attr) {
	return Var{ds, C.NC_GLOBAL}.Attr(name)
}

// AttrN returns global attribute for attribute number n.
func (ds Dataset) AttrN(n int) (a Attr, err error) {
	return Var{ds, C.NC_GLOBAL}.AttrN(n)
}
