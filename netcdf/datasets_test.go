// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package netcdf

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type FileTest struct {
	VarName  string
	DimNames []string
	DimLens  []uint64
	DataType Type
	Attr     map[string]interface{}
}

func (ft *FileTest) putAttrs(t *testing.T, v Var) {
	var err error
	for key, value := range ft.Attr {
		a := v.Attr(key)
		switch val := value.(type) {
		default:
			t.Fatalf("unexpected type %T\n", val)
		case []uint64:
			err = a.WriteUint64s(val)
		case []int64:
			err = a.WriteInt64s(val)
		case []float64:
			err = a.WriteFloat64s(val)
		case []uint32:
			err = a.WriteUint32s(val)
		case []int32:
			err = a.WriteInt32s(val)
		case []float32:
			err = a.WriteFloat32s(val)
		case []uint16:
			err = a.WriteUint16s(val)
		case []int16:
			err = a.WriteInt16s(val)
		case []uint8:
			err = a.WriteUint8s(val)
		case []int8:
			err = a.WriteInt8s(val)
		case string:
			err = a.WriteBytes([]byte(val))
		}
		if err != nil {
			t.Fatalf("writing attribute %s failed: %v\n", key, err)
		}
	}
}

func (ft *FileTest) getAttrs(t *testing.T, v Var) {
	for key, val := range ft.Attr {
		a := v.Attr(key)
		typ, err := a.Type()
		if err != nil {
			t.Fatalf("getting data type of attribute %s failed: %v\n", key, err)
		}
		var q interface{}
		switch typ {
		default:
			t.Errorf("unexpected attribute type %s\n", typeNames[typ])
		case UINT64:
			q, err = GetUint64s(a)
		case INT64:
			q, err = GetInt64s(a)
		case DOUBLE:
			q, err = GetFloat64s(a)
		case UINT:
			q, err = GetUint32s(a)
		case INT:
			q, err = GetInt32s(a)
		case FLOAT:
			q, err = GetFloat32s(a)
		case USHORT:
			q, err = GetUint16s(a)
		case SHORT:
			q, err = GetInt16s(a)
		case UBYTE:
			q, err = GetUint8s(a)
		case BYTE:
			q, err = GetInt8s(a)
		case CHAR:
			var b []byte
			b, err = GetBytes(a)
			q = string(b)
		}
		if err != nil {
			t.Fatalf("reading attribute %s failed: %v\n", key, err)
		}
		if !reflect.DeepEqual(q, val) {
			t.Errorf("attribute %s is %v; expected %v\n", key, q, val)
		}
	}
}

var fileTests = []FileTest{
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: UINT64,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: INT64,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: DOUBLE,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: UINT,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: INT,
		Attr: map[string]interface{}{
			"uint64_test": []uint64{0xFABCFABCFABCFABC, 999, 0, 222},
			"int64_test":  []int64{0x7ABC7ABC7ABC7ABC, -999, 0, 222},
			"double_test": []float64{3.14, 1.23, -5.7, 0, 7},
			"uint_test":   []uint32{2, 1, 600, 7},
			"int_test":    []int32{2, 1, -5, 7},
			"float_test":  []float32{3.14, 1.23, -5.7, 0, 7},
			"ushort_test": []uint16{2, 1, 600, 7},
			"short_test":  []int16{2, 1, -5, 7},
			"ubyte_test":  []uint8{2, 1, 255, 0, 7},
			"byte_test":   []int8{2, 100, -128, 127, -17},
			"birthday":    "2009-11-10",
		},
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: FLOAT,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: USHORT,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: SHORT,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: UBYTE,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: BYTE,
	},
	{
		VarName:  "gopher",
		DimNames: []string{"height", "width"},
		DimLens:  []uint64{7, 3},
		DataType: CHAR,
	},
}

func TestCreate(t *testing.T) {
	for _, ft := range fileTests {
		f, err := ioutil.TempFile("", "netcdf_test")
		if err != nil {
			t.Fatalf("creating temporary file failed: %v\n", err)
		}
		createFile(t, f.Name(), &ft)
		readFile(t, f.Name(), &ft)
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("removing temporary file failed: %v\n", err)
		}
	}
}

func createFile(t *testing.T, filename string, ft *FileTest) {
	f, err := CreateFile(filename, CLOBBER|NETCDF4)
	if err != nil {
		t.Fatalf("Create failed: %v\n", err)
	}
	dims := make([]Dim, 2)
	for i, name := range ft.DimNames {
		if dims[i], err = f.AddDim(name, ft.DimLens[i]); err != nil {
			t.Fatalf("PutDim failed: %v\n", err)
		}
	}
	v, err := f.AddVar(ft.VarName, ft.DataType, dims)
	if err != nil {
		t.Fatalf("PutVar failed: %v\n", err)
	}
	ft.putAttrs(t, v)

	n, err := v.Len()
	if err != nil {
		t.Fatalf("Var.Len failed: %v\n", err)
	}
	switch ft.DataType {
	default:
		t.Fatalf("unexpected type %s\n", typeNames[ft.DataType])
	case UINT64:
		err = testWriteUint64s(v, n)
	case INT64:
		err = testWriteInt64s(v, n)
	case DOUBLE:
		err = testWriteFloat64s(v, n)
	case UINT:
		err = testWriteUint32s(v, n)
	case INT:
		err = testWriteInt32s(v, n)
	case FLOAT:
		err = testWriteFloat32s(v, n)
	case USHORT:
		err = testWriteUint16s(v, n)
	case SHORT:
		err = testWriteInt16s(v, n)
	case UBYTE:
		err = testWriteUint8s(v, n)
	case BYTE:
		err = testWriteInt8s(v, n)
	case CHAR:
		err = testWriteBytes(v, n)
	}
	if err != nil {
		t.Errorf("writing data failed: %v\n", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
}

func readFile(t *testing.T, filename string, ft *FileTest) {
	f, err := OpenFile(filename, NOWRITE)
	if err != nil {
		t.Fatalf("Open failed: %v\n", err)
	}
	for i, name := range ft.DimNames {
		d, err := f.Dim(name)
		if err != nil {
			t.Fatalf("GetDim failed: %v\n", err)
		}
		s, err := d.Name()
		if err != nil {
			t.Fatalf("Dim.Name failed: %v\n", err)
		}
		if err == nil && s != name {
			t.Fatalf("Dim name is %q; expected %q\n", s, name)
		}
		n, err := d.Len()
		if err != nil {
			t.Fatalf("Dim.Len failed: %v\n", err)
		}
		if err == nil && n != ft.DimLens[i] {
			t.Fatalf("Dim length is %d; expected %d\n", n, ft.DimLens[i])
		}
	}
	v, err := f.Var(ft.VarName)
	if err != nil {
		t.Errorf("GetVar failed: %v\n", err)
	}
	ft.getAttrs(t, v)

	n, err := v.Len()
	if err != nil {
		t.Fatalf("Var.Len failed: %v\n", err)
	}
	switch ft.DataType {
	default:
		t.Fatalf("unexpected type %s\n", typeNames[ft.DataType])
	case UINT64:
		err = testReadUint64s(v, n)
	case INT64:
		err = testReadInt64s(v, n)
	case DOUBLE:
		err = testReadFloat64s(v, n)
	case UINT:
		err = testReadUint32s(v, n)
	case INT:
		err = testReadInt32s(v, n)
	case FLOAT:
		err = testReadFloat32s(v, n)
	case USHORT:
		err = testReadUint16s(v, n)
	case SHORT:
		err = testReadInt16s(v, n)
	case UBYTE:
		err = testReadUint8s(v, n)
	case BYTE:
		err = testReadInt8s(v, n)
	case CHAR:
		err = testReadBytes(v, n)
	}
	if err != nil {
		t.Fatalf("reading data failed: %v\n", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
}

func TestError(t *testing.T) {
	_, err := OpenFile("/non-existant.nc", NOWRITE)
	if err == nil {
		t.Fatalf("Opened non-existant file\n")
	}
	if len(err.Error()) == 0 {
		t.Errorf("empty error\n")
	}
}

func TestEndDef(t *testing.T) {
	f, err := ioutil.TempFile("", "netcdf_test")
	if err != nil {
		t.Fatalf("creating temporary file failed: %v\n", err)
	}
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("removing temporary file failed: %v\n", err)
		}
	}()

	// Create a new NetCDF 3 file.
	ds, err := CreateFile(f.Name(), CLOBBER)
	if err != nil {
		t.Fatalf("creating file failed: %v\n", err)
	}
	defer ds.Close()

	size, err := ds.AddDim("size", 5)
	if err != nil {
		t.Fatalf("adding dimension failed: %v\n", err)
	}
	v, err := ds.AddVar("gopher", INT, []Dim{size})
	if err != nil {
		t.Fatalf("adding variable failed: %v\n", err)
	}

	// writing data will fail unless we leave define mode
	if err := ds.EndDef(); err != nil {
		t.Fatalf("failed to end define mode: %v\n", err)
	}
	if err := v.WriteInt32s([]int32{1, 2, 3, 4, 5}); err != nil {
		t.Fatalf("writing data failed: %v\n", err)
	}
}
