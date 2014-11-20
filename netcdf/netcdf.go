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

import (
	"fmt"
)

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
		return fmt.Errorf("wrong data type %s; expected %s", typeNames[u], typeNames[t])
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
