// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package netcdf is a Go binding for the netCDF C library.
//
// This package supports netCDF version 3, and 4 if
// netCDF 4 support is enabled in the C library.
// The C library interface used is documented here:
// http://www.unidata.ucar.edu/software/netcdf/docs/netcdf-c/
package netcdf

// #cgo pkg-config: netcdf
// #include <stdlib.h>
// #include <netcdf.h>
import "C"

import (
	"unsafe"
)

// Error represents an error returned by netCDF C library.
type Error C.int

func newError(n C.int) error {
	if n == _NC_NOERR {
		return nil
	}
	return Error(n)
}

// Error returns a string representation of Error e.
func (e Error) Error() string {
	return C.GoString(C.nc_strerror(C.int(e)))
}

// File represents an open netCDF file.
type File C.int

// Create creates a new netCDF dataset.
// Mode is a bitwise-or of FileMode values.
func Create(path string, mode FileMode) (f File, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var id C.int
	err = newError(C.nc_create(cpath, C.int(mode), &id))
	f = File(id)
	return
}

// Open opens an existing nefCDF dataset file at path.
// Mode is a bitwise-or of FileMode values.
func Open(path string, mode FileMode) (f File, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var id C.int
	err = newError(C.nc_open(cpath, C.int(mode), &id))
	f = File(id)
	return
}

// Close closes an open netCDF dataset.
func (f File) Close() (err error) {
	return newError(C.nc_close(C.int(f)))
}
