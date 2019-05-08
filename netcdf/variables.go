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

// Var represents a variable.
type Var struct {
	ds Dataset
	id C.int
}

// Dims returns the dimensions of variable v.
func (v Var) Dims() (dims []Dim, err error) {
	var ndims C.int
	err = newError(C.nc_inq_varndims(C.int(v.ds), C.int(v.id), &ndims))
	if err != nil {
		return
	}
	if ndims == 0 {
		return
	}
	dimids := make([]C.int, ndims)
	err = newError(C.nc_inq_vardimid(C.int(v.ds), C.int(v.id), &dimids[0]))
	if err != nil {
		return
	}
	dims = make([]Dim, ndims)
	for i, id := range dimids {
		dims[i] = Dim{v.ds, id}
	}
	return
}

// Type returns the data type of variable v.
func (v Var) Type() (t Type, err error) {
	var typ C.nc_type
	err = newError(C.nc_inq_vartype(C.int(v.ds), C.int(v.id), &typ))
	t = Type(typ)
	return
}

// Len returns the total number of values in the variable v.
func (v Var) Len() (uint64, error) {
	dims, err := v.Dims()
	if err != nil {
		return 0, err
	}
	n := uint64(1)
	for _, d := range dims {
		len, err := d.Len()
		if err != nil {
			return 0, err
		}
		n *= len
	}
	return n, nil
}

// LenDims returns the length of the dimensions of variable v.
func (v Var) LenDims() ([]uint64, error) {
	dims, err := v.Dims()
	if err != nil {
		return nil, err
	}
	ls := make([]uint64, len(dims))
	for i, d := range dims {
		ls[i], err = d.Len()
		if err != nil {
			return nil, err
		}
	}
	return ls, nil
}

// NAttrs returns the number of attributes assigned to variable v.
func (v Var) NAttrs() (n int, err error) {
	var cn C.int
	err = newError(C.nc_inq_varnatts(C.int(v.ds), C.int(v.id), &cn))
	n = int(cn)
	return
}

// Name returns the name of the variable.
func (v Var) Name() (name string, err error) {
	buf := C.CString(string(make([]byte, C.NC_MAX_NAME+1)))
	defer C.free(unsafe.Pointer(buf))
	err = newError(C.nc_inq_varname(C.int(v.ds), v.id, buf))
	name = C.GoString(buf)
	return
}

// SetCompression sets the deflate parameters for a variable in a NetCDF-4 file.
func (v Var) SetCompression(shuffle, deflate bool, deflateLevel int) error {
	sInt := C.int(0)
	if shuffle {
		sInt = 1
	}
	dInt := C.int(0)
	if deflate {
		dInt = 1
	}
	return newError(C.nc_def_var_deflate(C.int(v.ds), v.id, sInt, dInt, C.int(deflateLevel)))
}

// Compression returns the deflate settings for a variable in a NetCDF-4 file.
func (v Var) Compression() (shuffle, deflate bool, deflateLevel int, err error) {
	var (
		cShuffle      C.int
		cDeflate      C.int
		cDeflateLevel C.int
	)
	err = newError(C.nc_inq_var_deflate(C.int(v.ds), v.id, &cShuffle, &cDeflate, &cDeflateLevel))
	return cShuffle != 0, cDeflate != 0, int(cDeflateLevel), err
}

// AddVar adds a new a variable named name of type t and dimensions dims.
// The new variable v is returned.
func (ds Dataset) AddVar(name string, t Type, dims []Dim) (v Var, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var varid C.int
	var dimPtr *C.int
	if len(dims) > 0 {
		dimids := make([]C.int, len(dims))
		for i, d := range dims {
			dimids[i] = d.id
		}
		dimPtr = &dimids[0]
	}
	err = newError(C.nc_def_var(C.int(ds), cname, C.nc_type(t),
		C.int(len(dims)), dimPtr, &varid))
	v = Var{ds, varid}
	return
}

// VarN returns a new variable in File f with ID id.
func (ds Dataset) VarN(id int) Var {
	return Var{ds, C.int(id)}
}

// Var returns the Var for the variable named name.
func (ds Dataset) Var(name string) (v Var, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var id C.int
	err = newError(C.nc_inq_varid(C.int(ds), cname, &id))
	v = Var{ds, id}
	return
}
