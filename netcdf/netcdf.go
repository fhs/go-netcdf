// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//go:generate go run generate.go

// Package netcdf is a Go binding for the netCDF C library.
//
// This package supports netCDF version 3, and 4 if
// netCDF 4 support is enabled in the C library.
// The C library interface used is documented here:
// http://www.unidata.ucar.edu/software/netcdf/docs/netcdf-c/
package netcdf

import "fmt"

type typedArray interface {
	Type() (Type, error)
	Len() (uint64, error)
}

// okData checks if t agrees with a.Type() and n agrees with a.Len().
func okData(a typedArray, t Type, n int) error {
	u, err := a.Type()
	if err != nil {
		return err
	}
	if u != t {
		return fmt.Errorf("wrong data type %v; expected %v", u, t)
	}
	m, err := a.Len()
	if err != nil {
		return err
	}
	if n < int(m) {
		return fmt.Errorf("data length %d is smaller than %d", n, m)
	}
	return nil
}

func product(nums []uint64) (prod uint64) {
	prod = 1
	for _, i := range nums {
		prod *= i
	}
	return
}

// UnravelIndex calculates coordinate position based on index
func UnravelIndex(idx uint64, shape []uint64) ([]uint64, error) {
	for _, v := range shape {
		if v == 0 {
			return nil, fmt.Errorf("invalid shape, 0 encountered in shape %v", shape)
		}
	}

	if idx > product(shape) {
		return nil, fmt.Errorf("index %v > size %v of shape", idx, shape)
	}

	var maxval = product(shape)
	var ndim = len(shape)
	var coord = make([]uint64, ndim)

	for i := 0; i < ndim; i++ {
		shape[i] = 1
		maxval = product(shape)
		coord[i] = uint64(idx / maxval)
		idx -= coord[i] * maxval
	}
	return coord, nil
}
