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

// Dim represents a dimension.
type Dim struct {
	f  File
	id C.int
}

// Name returns the name of dimension d.
func (d Dim) Name() (name string, err error) {
	buf := C.CString(string(make([]byte, _NC_MAX_NAME+1)))
	defer C.free(unsafe.Pointer(buf))
	err = newError(C.nc_inq_dimname(C.int(d.f), d.id, buf))
	name = C.GoString(buf)
	return
}

// Len returns the length of dimension d.
func (d Dim) Len() (n uint64, err error) {
	var len C.size_t
	err = newError(C.nc_inq_dimlen(C.int(d.f), d.id, &len))
	n = uint64(len)
	return
}

// PutDim adds a new dimension named name of length len.
// The new dimension d is returned.
func (f File) PutDim(name string, len uint64) (d Dim, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var dimid C.int
	err = newError(C.nc_def_dim(C.int(f), cname, C.size_t(len), &dimid))
	d = Dim{f, dimid}
	return
}

// GetDim returns the Dim for the dimension named name.
func (f File) GetDim(name string) (d Dim, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var id C.int
	err = newError(C.nc_inq_dimid(C.int(f), cname, &id))
	d = Dim{f, id}
	return
}
