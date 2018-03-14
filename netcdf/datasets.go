// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

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
	if n == C.NC_NOERR {
		return nil
	}
	return Error(n)
}

// Error returns a string representation of Error e.
func (e Error) Error() string {
	return C.GoString(C.nc_strerror(C.int(e)))
}

// Dataset represents a netCDF dataset.
type Dataset C.int

// CreateFile creates a new netCDF dataset.
// Mode is a bitwise-or of FileMode values.
func CreateFile(path string, mode FileMode) (ds Dataset, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var id C.int
	err = newError(C.nc_create(cpath, C.int(mode), &id))
	ds = Dataset(id)
	return
}

// OpenFile opens an existing netCDF dataset file at path.
// Mode is a bitwise-or of FileMode values.
func OpenFile(path string, mode FileMode) (ds Dataset, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var id C.int
	err = newError(C.nc_open(cpath, C.int(mode), &id))
	ds = Dataset(id)
	return
}

// Close closes an open netCDF dataset.
func (ds Dataset) Close() (err error) {
	return newError(C.nc_close(C.int(ds)))
}

// EndDef leaves define mode and enters data mode, so variable data
// can be read or written. Calling this method is not required
// for netCDF-4 files.
func (ds Dataset) EndDef() (err error) {
	return newError(C.nc_enddef(C.int(ds)))
}

// NVars returns the number of variables defined for dataset f.
func (ds Dataset) NVars() (n int, err error) {
	var cn C.int
	err = newError(C.nc_inq_nvars(C.int(ds), &cn))
	n = int(cn)
	return
}

// NAttrs returns the number of global attributes defined for dataset f.
func (ds Dataset) NAttrs() (n int, err error) {
	var cn C.int
	err = newError(C.nc_inq_natts(C.int(ds), &cn))
	n = int(cn)
	return
}

// Version returns a string identifying the version of the netCDF library,
// and when it was built.
func Version() string {
	return C.GoString(C.nc_inq_libvers())
}
