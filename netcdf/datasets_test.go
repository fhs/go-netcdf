// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package netcdf

import (
	"io/ioutil"
	"os"
	"testing"
)

type FileTest struct {
	VarName  string
	DimNames []string
	DimLens  []uint64
	DataType Type
	Attr     map[string]string
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
	for key, val := range ft.Attr {
		if err := v.Attr(key).WriteChar([]byte(val)); err != nil {
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
		switch typ {
		case NC_CHAR:
			b, err := GetChar(a)
			if err != nil {
				t.Fatalf("read attribute %s failed: %v\n", key, err)
			}
			if string(b) != val {
				t.Errorf("attribute %s is %s; expected %s\n", key, string(b), val)
			}
		default:
			t.Errorf("unexpected attribute type %s\n", typeNames[typ])
		}
	}
}

var fileTests = []FileTest{
	{
		VarName:  "golang",
		DimNames: []string{"time", "growth"},
		DimLens:  []uint64{7, 3},
		DataType: NC_INT,
		Attr:     map[string]string{"birthday": "2009-11-10"},
	},
	{
		VarName:  "golang",
		DimNames: []string{"time", "growth"},
		DimLens:  []uint64{7, 3},
		DataType: NC_FLOAT,
		Attr:     map[string]string{"birthday": "2009-11-10"},
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
	case NC_INT:
		ft.putInt(t, v, n)
	case NC_FLOAT:
		ft.putFloat(t, v, n)
	default:
		t.Fatalf("unexpected type %s\n", typeNames[ft.DataType])
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
	case NC_INT:
		ft.getInt(t, v, n)
	case NC_FLOAT:
		ft.getFloat(t, v, n)
	default:
		t.Fatalf("unexpected type %s\n", typeNames[ft.DataType])
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
