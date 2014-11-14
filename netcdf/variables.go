// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package netcdf

// #include <stdlib.h>
// #include <netcdf.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// Var represents a variable.
type Var struct {
	f  File
	id C.int
}

// Dims returns the dimensions of variable v.
func (v Var) Dims() (dims []Dim, err error) {
	var ndims C.int
	err = newError(C.nc_inq_varndims(C.int(v.f), C.int(v.id), &ndims))
	if err != nil {
		return
	}
	dimids := make([]C.int, ndims)
	err = newError(C.nc_inq_vardimid(C.int(v.f), C.int(v.id), &dimids[0]))
	if err != nil {
		return
	}
	dims = make([]Dim, ndims)
	for i, id := range dimids {
		dims[i] = Dim{f: v.f, id: id}
	}
	return
}

// Type returns the data type of variable v.
func (v Var) Type() (t Type, err error) {
	var typ C.nc_type
	err = newError(C.nc_inq_vartype(C.int(v.f), C.int(v.id), &typ))
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

// okData checks if t agrees with v.Type() and n agrees with v.Len().
func (v Var) okData(t Type, n int) error {
	u, err := v.Type()
	if err != nil {
		return err
	}
	if u != t {
		return fmt.Errorf("wrong data type %s; expected %s", typeNames[u], typeNames[t])
	}

	m, err := v.Len()
	if err != nil {
		return err
	}
	if n < int(m) {
		return fmt.Errorf("data length %d is smaller than %d", n, m)
	}
	return nil
}

// PutInt writes data as the entire data for variable v.
func (v Var) PutInt(data []int32) error {
	if err := v.okData(NC_INT, len(data)); err != nil {
		return err
	}
	return newError(C.nc_put_var_int(C.int(v.f), C.int(v.id), (*C.int)(unsafe.Pointer(&data[0]))))
}

// PutInt writes data as the entire data for variable v.
func (v Var) PutFloat(data []float32) error {
	if err := v.okData(NC_FLOAT, len(data)); err != nil {
		return err
	}
	return newError(C.nc_put_var_float(C.int(v.f), C.int(v.id), (*C.float)(unsafe.Pointer(&data[0]))))
}

// GetInt reads the entire variable v into data, which must have enough
// space for all the values (i.e. len(data) must be at least v.Len()).
func (v Var) GetInt(data []int32) error {
	if err := v.okData(NC_INT, len(data)); err != nil {
		return err
	}
	return newError(C.nc_get_var_int(C.int(v.f), C.int(v.id), (*C.int)(unsafe.Pointer(&data[0]))))
}

// GetFloat reads the entire variable v into data, which must have enough
// space for all the values (i.e. len(data) must be at least v.Len()).
func (v Var) GetFloat(data []float32) error {
	if err := v.okData(NC_FLOAT, len(data)); err != nil {
		return err
	}
	return newError(C.nc_get_var_float(C.int(v.f), C.int(v.id), (*C.float)(unsafe.Pointer(&data[0]))))
}

// PutVar adds a new a variable named name of type t and dimensions dims.
// The new variable v is returned.
func (f File) PutVar(name string, t Type, dims []Dim) (v Var, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var varid C.int
	dimids := make([]C.int, len(dims))
	for i, d := range dims {
		dimids[i] = d.id
	}
	err = newError(C.nc_def_var(C.int(f), cname, C.nc_type(t),
		C.int(len(dimids)), &dimids[0], &varid))
	v = Var{f: f, id: varid}
	return
}

// GetVar returns the Var for the variable named name.
func (f File) GetVar(name string) (d Var, err error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var id C.int
	err = newError(C.nc_inq_varid(C.int(f), cname, &id))
	d = Var{f, id}
	return
}
