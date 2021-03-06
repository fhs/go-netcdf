// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// These files are autogenerated from nc_double.go using generate.go
// DO NOT EDIT (except nc_double.go).

package netcdf

import (
	"fmt"
	"unsafe"
)

// #include <stdlib.h>
// #include <netcdf.h>
import "C"

// WriteInt8s writes data as the entire data for variable v.
func (v Var) WriteInt8s(data []int8) error {
	if err := okData(v, BYTE, len(data)); err != nil {
		return err
	}
	return newError(C.nc_put_var_schar(C.int(v.ds), C.int(v.id), (*C.schar)(unsafe.Pointer(&data[0]))))
}

// ReadInt8s reads the entire variable v into data, which must have enough
// space for all the values (i.e. len(data) must be at least v.Len()).
func (v Var) ReadInt8s(data []int8) error {
	if err := okData(v, BYTE, len(data)); err != nil {
		return err
	}
	return newError(C.nc_get_var_schar(C.int(v.ds), C.int(v.id), (*C.schar)(unsafe.Pointer(&data[0]))))
}

// WriteInt8s sets the value of attribute a to val.
func (a Attr) WriteInt8s(val []int8) error {
	// We don't need okData here because netcdf library doesn't know
	// the length or type of the attribute yet.
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	return newError(C.nc_put_att_schar(C.int(a.v.ds), C.int(a.v.id), cname,
		C.nc_type(BYTE), C.size_t(len(val)), (*C.schar)(unsafe.Pointer(&val[0]))))
}

// ReadInt8s reads the entire attribute value into val.
func (a Attr) ReadInt8s(val []int8) (err error) {
	if err := okData(a, BYTE, len(val)); err != nil {
		return err
	}
	cname := C.CString(a.name)
	defer C.free(unsafe.Pointer(cname))
	err = newError(C.nc_get_att_schar(C.int(a.v.ds), C.int(a.v.id), cname,
		(*C.schar)(unsafe.Pointer(&val[0]))))
	return
}

// ReadInt8At returns a value via index position
func (v Var) ReadInt8At(idx []uint64) (val int8, err error) {
	var dimPtr *C.size_t
	if len(idx) > 0 {
		dimPtr = (*C.size_t)(unsafe.Pointer(&idx[0]))
	}
	err = newError(C.nc_get_var1_schar(C.int(v.ds), C.int(v.id),
		dimPtr, (*C.schar)(unsafe.Pointer(&val))))
	return
}

// WriteInt8At sets a value via its index position
func (v Var) WriteInt8At(idx []uint64, val int8) (err error) {
	var dimPtr *C.size_t
	if len(idx) > 0 {
		dimPtr = (*C.size_t)(unsafe.Pointer(&idx[0]))
	}
	err = newError(C.nc_put_var1_schar(C.int(v.ds), C.int(v.id),
		dimPtr, (*C.schar)(unsafe.Pointer(&val))))
	return
}

// WriteInt8Slice writes data as a slice of variable v. The slice is specified by start and count:
// https://www.unidata.ucar.edu/software/netcdf/docs/programming_notes.html#specify_hyperslab.
func (v Var) WriteInt8Slice(data []int8, start, count []uint64) error {
	if err := okDataSlice(v, BYTE, len(data), start, count); err != nil {
		return err
	}
	return newError(C.nc_put_vara_schar(C.int(v.ds), C.int(v.id),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.schar)(unsafe.Pointer(&data[0])),
	))
}

// ReadInt8Slice reads a slice of variable v into data, which must have enough
// space for all the values. The slice is specified by start and count:
// https://www.unidata.ucar.edu/software/netcdf/docs/programming_notes.html#specify_hyperslab.
func (v Var) ReadInt8Slice(data []int8, start, count []uint64) error {
	if err := okDataSlice(v, BYTE, len(data), start, count); err != nil {
		return err
	}
	return newError(C.nc_get_vara_schar(C.int(v.ds), C.int(v.id),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.schar)(unsafe.Pointer(&data[0])),
	))
}

// WriteInt8StridedSlice writes data as a slice of variable v. The slice is specified by start, count and stride:
// https://www.unidata.ucar.edu/software/netcdf/docs/programming_notes.html#specify_hyperslab.
func (v Var) WriteInt8StridedSlice(data []int8, start, count []uint64, stride []int64) error {
	if err := okDataStride(v, BYTE, len(data), start, count, stride); err != nil {
		return err
	}
	return newError(C.nc_put_vars_schar(C.int(v.ds), C.int(v.id),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.ptrdiff_t)(unsafe.Pointer(&stride[0])),
		(*C.schar)(unsafe.Pointer(&data[0])),
	))
}

// ReadInt8StridedSlice reads a strided slice of variable v into data, which must have enough
// space for all the values. The slice is specified by start, count and stride:
// https://www.unidata.ucar.edu/software/netcdf/docs/programming_notes.html#specify_hyperslab.
func (v Var) ReadInt8StridedSlice(data []int8, start, count []uint64, stride []int64) error {
	if err := okDataStride(v, BYTE, len(data), start, count, stride); err != nil {
		return err
	}
	return newError(C.nc_get_vars_schar(C.int(v.ds), C.int(v.id),
		(*C.size_t)(unsafe.Pointer(&start[0])),
		(*C.size_t)(unsafe.Pointer(&count[0])),
		(*C.ptrdiff_t)(unsafe.Pointer(&stride[0])),
		(*C.schar)(unsafe.Pointer(&data[0])),
	))
}

// Int8sReader is a interface that allows reading a sequence of values of fixed length.
type Int8sReader interface {
	Len() (n uint64, err error)
	ReadInt8s(val []int8) (err error)
}

// GetInt8s reads the entire data in r and returns it.
func GetInt8s(r Int8sReader) (data []int8, err error) {
	n, err := r.Len()
	if err != nil {
		return
	}
	data = make([]int8, n)
	err = r.ReadInt8s(data)
	return
}

// testReadInt8s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteInt8s(v Var, n uint64) error {
	data := make([]int8, n)
	for i := 0; i < int(n); i++ {
		data[i] = int8(i + 10)
	}
	return v.WriteInt8s(data)
}

// testReadInt8s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadInt8s(v Var, n uint64) error {
	data := make([]int8, n)
	if err := v.ReadInt8s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := int8(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

// testWriteInt8Slice writes somes data to v. N is v.LenDim().
// This function is only used for testing.
func testWriteInt8Slice(v Var, n []uint64) error {
	if len(n) == 0 {
		return nil // Don't test empty data.
	}
	start, count := make([]uint64, len(n)), make([]uint64, len(n))
	for i, v := range n {
		start[i] = v / 2
		count[i] = v / 2
	}
	data := make([]int8, product(count))
	for i := 0; i < int(product(count)); i++ {
		data[i] = int8(i + 10)
	}
	return v.WriteInt8Slice(data, start, count)
}

// testReadInt8Slice reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.LenDim().
// This function is only used for testing.
func testReadInt8Slice(v Var, n []uint64) error {
	if len(n) == 0 {
		return nil // Don't test empty data.
	}
	start, count := make([]uint64, len(n)), make([]uint64, len(n))
	for i, v := range n {
		start[i] = v / 2
		count[i] = v / 2
	}
	data := make([]int8, product(count))
	if err := v.ReadInt8Slice(data, start, count); err != nil {
		return err
	}
	for i := 0; i < int(product(count)); i++ {
		if val := int8(i + 10); data[i] != val {
			return fmt.Errorf("strided slice data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

// testWriteInt8StridedSlice writes somes data to v. N is v.LenDim().
// This function is only used for testing.
func testWriteInt8StridedSlice(v Var, n []uint64) error {
	if len(n) == 0 {
		return nil // Don't test empty data.
	}
	start, count, stride := make([]uint64, len(n)), make([]uint64, len(n)), make([]int64, len(n))
	for i, v := range n {
		start[i] = 1
		count[i] = (v - 1) / 2
		stride[i] = 2
	}
	data := make([]int8, product(count))
	for i := 0; i < int(product(count)); i++ {
		data[i] = int8(i + 10)
	}
	return v.WriteInt8StridedSlice(data, start, count, stride)
}

// testReadInt8StridedSlice reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.LenDim().
// This function is only used for testing.
func testReadInt8StridedSlice(v Var, n []uint64) error {
	if len(n) == 0 {
		return nil // Don't test empty data.
	}
	start, count, stride := make([]uint64, len(n)), make([]uint64, len(n)), make([]int64, len(n))
	for i, v := range n {
		start[i] = 1
		count[i] = (v - 1) / 2
		stride[i] = 2
	}
	data := make([]int8, product(count))
	if err := v.ReadInt8StridedSlice(data, start, count, stride); err != nil {
		return err
	}
	for i := 0; i < int(product(count)); i++ {
		if val := int8(i + 10); data[i] != val {
			return fmt.Errorf("strided slice data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadInt8At(v Var, n uint64) error {
	data := make([]int8, n)
	if err := v.ReadInt8s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadInt8At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteInt8At(v Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteInt8At(coord, int8(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadInt8At(coord)
		if val != int8(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}
