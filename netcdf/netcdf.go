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

type sliceableTypedArray interface {
	Type() (Type, error)
	Len() (uint64, error)
	LenDims() ([]uint64, error)
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

// okDataSlice checks if t agrees with a.Type() and n agrees with count.
func okDataSlice(a sliceableTypedArray, t Type, n int, start, count []uint64) error {
	u, err := a.Type()
	if err != nil {
		return err
	}
	if u != t {
		return fmt.Errorf("wrong data type %v; expected %v", u, t)
	}

	d, err := a.LenDims()
	if err != nil {
		return err
	}
	if len(start) != len(d) {
		return fmt.Errorf("incorrect number of dimensions in start: %d != %d", start, len(d))
	}
	if len(count) != len(d) {
		return fmt.Errorf("incorrect number of dimensions in count: %d != %d", count, len(d))
	}

	for i, id := range d {
		if start[i] >= id || start[i] < 0 {
			return fmt.Errorf("start of dimension %d of slice is out of range: 0 <= %d < %d", i, start, id)
		}
		v := start[i] + count[i]
		if v > id || v <= 0 {
			return fmt.Errorf("end of dimension %d of slice is out of range: 0 < %d <= %d", i, v, id)
		}
	}

	l := product(count)
	if n < int(l) {
		return fmt.Errorf("data length %d is smaller than %d", n, l)
	}
	return nil
}

// okDataStride checks if t agrees with a.Type() and n agrees with start, count and stride.
func okDataStride(a sliceableTypedArray, t Type, n int, start, count []uint64, stride []int64) error {
	u, err := a.Type()
	if err != nil {
		return err
	}
	if u != t {
		return fmt.Errorf("wrong data type %v; expected %v", u, t)
	}

	d, err := a.LenDims()
	if err != nil {
		return err
	}
	if len(start) != len(d) {
		return fmt.Errorf("incorrect number of dimensions in start: %d != %d", start, len(d))
	}
	if len(count) != len(d) {
		return fmt.Errorf("incorrect number of dimensions in count: %d != %d", count, len(d))
	}
	if len(stride) != len(d) {
		return fmt.Errorf("incorrect number of dimensions in stride: %d != %d", stride, len(d))
	}

	for i, id := range d {
		if start[i] >= id || start[i] < 0 {
			return fmt.Errorf("start of dimension %d of slice is out of range: 0 <= %d < %d", i, start, id)
		}
		v := int64(start[i]) + int64(count[i])*stride[i]
		if v > int64(id) || v <= 0 {
			return fmt.Errorf("end of dimension %d of slice is out of range: 0 < %d <= %d", i, v, id)
		}
	}

	l := product(count)
	if n < int(l) {
		return fmt.Errorf("data length %d is smaller than %d", n, l)
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
