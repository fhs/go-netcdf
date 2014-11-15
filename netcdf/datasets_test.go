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

func genData(i int) float64 {
	return float64(i - 10)
}

func (ft *FileTest) putInt(t *testing.T, v Var, n uint64) {
	data := make([]int32, n)
	for i := 0; i < int(n); i++ {
		data[i] = int32(genData(i))
	}
	if err := v.WriteInt(data); err != nil {
		t.Fatalf("WriteInt failed: %v\n", err)
	}
}

func (ft *FileTest) putFloat(t *testing.T, v Var, n uint64) {
	data := make([]float32, n)
	for i := 0; i < int(n); i++ {
		data[i] = float32(genData(i))
	}
	if err := v.WriteFloat(data); err != nil {
		t.Fatalf("WriteFloat failed: %v\n", err)
	}
}

func (ft *FileTest) getInt(t *testing.T, v Var, n uint64) {
	data := make([]int32, n)
	if err := v.ReadInt(data); err != nil {
		t.Fatalf("ReadInt failed: %v\n", err)
	}
	for i := 0; i < int(n); i++ {
		if val := int32(genData(i)); data[i] != val {
			t.Errorf("data at position %d is %d; expected %d\n", i, data[i], val)
		}
	}
}

func (ft *FileTest) getFloat(t *testing.T, v Var, n uint64) {
	data := make([]float32, n)
	if err := v.ReadFloat(data); err != nil {
		t.Fatalf("ReadFloat failed: %v\n", err)
	}
	for i := 0; i < int(n); i++ {
		if val := float32(genData(i)); data[i] != val {
			t.Errorf("data at position %d is %f; expected %f\n", i, data[i], val)
		}
	}
}

func (ft *FileTest) putAttrs(t *testing.T, v Var) {
	var err error
	for key, value := range ft.Attr {
		a := v.Attr(key)
		switch val := value.(type) {
		default:
			t.Fatalf("unexpected type %T\n", val)
		case []uint64:
			err = a.WriteUint64(val)
		case []int64:
			err = a.WriteInt64(val)
		case []float64:
			err = a.WriteDouble(val)
		case []uint32:
			err = a.WriteUint(val)
		case []int32:
			err = a.WriteInt(val)
		case []float32:
			err = a.WriteFloat(val)
		case []uint16:
			err = a.WriteUshort(val)
		case []int16:
			err = a.WriteShort(val)
		case []uint8:
			err = a.WriteUbyte(val)
		case []int8:
			err = a.WriteByte(val)
		case string:
			err = a.WriteChar([]byte(val))
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
		case NC_UINT64:
			q, err = GetUint64(a)
		case NC_INT64:
			q, err = GetInt64(a)
		case NC_DOUBLE:
			q, err = GetDouble(a)
		case NC_UINT:
			q, err = GetUint(a)
		case NC_INT:
			q, err = GetInt(a)
		case NC_FLOAT:
			q, err = GetFloat(a)
		case NC_USHORT:
			q, err = GetUshort(a)
		case NC_SHORT:
			q, err = GetShort(a)
		case NC_UBYTE:
			q, err = GetUbyte(a)
		case NC_BYTE:
			q, err = GetByte(a)
		case NC_CHAR:
			var b []byte
			b, err = GetChar(a)
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
		VarName:  "golang",
		DimNames: []string{"time", "growth"},
		DimLens:  []uint64{7, 3},
		DataType: NC_INT,
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
		VarName:  "golang",
		DimNames: []string{"time", "growth"},
		DimLens:  []uint64{7, 3},
		DataType: NC_FLOAT,
	},
}

func TestCreate(t *testing.T) {
	for _, ft := range fileTests {
		f, err := ioutil.TempFile("", "netcdf_test")
		if err != nil {
			t.Fatalf("createing temporary file failed: %v\n", err)
		}
		createFile(t, f.Name(), &ft)
		readFile(t, f.Name(), &ft)
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("removing temporary file failed: %v\n", err)
		}
	}
}

func createFile(t *testing.T, filename string, ft *FileTest) {
	f, err := Create(filename, NC_CLOBBER|NC_NETCDF4)
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
	case NC_INT:
		ft.putInt(t, v, n)
	case NC_FLOAT:
		ft.putFloat(t, v, n)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
}

func readFile(t *testing.T, filename string, ft *FileTest) {
	f, err := Open(filename, NC_NOWRITE)
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
	case NC_INT:
		ft.getInt(t, v, n)
	case NC_FLOAT:
		ft.getFloat(t, v, n)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
}

func TestError(t *testing.T) {
	_, err := Open("/non-existant.nc", NC_NOWRITE)
	if err == nil {
		t.Fatalf("Opened non-existant file\n")
	}
	if len(err.Error()) == 0 {
		t.Errorf("empty error\n")
	}
}
