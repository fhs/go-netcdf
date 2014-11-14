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
}

func genData(i int) float64 {
	return float64(i - 10)
}

func (ft *FileTest) putInt(t *testing.T, v Var, n uint64) {
	data := make([]int32, n)
	for i := 0; i < int(n); i++ {
		data[i] = int32(genData(i))
	}
	if err := v.PutInt(data); err != nil {
		t.Fatalf("PutInt failed: %v\n", err)
	}
}

func (ft *FileTest) putFloat(t *testing.T, v Var, n uint64) {
	data := make([]float32, n)
	for i := 0; i < int(n); i++ {
		data[i] = float32(genData(i))
	}
	if err := v.PutFloat(data); err != nil {
		t.Fatalf("PutFloat failed: %v\n", err)
	}
}

func (ft *FileTest) getInt(t *testing.T, v Var, n uint64) {
	data := make([]int32, n)
	if err := v.GetInt(data); err != nil {
		t.Fatalf("GetInt failed: %v\n", err)
	}
	for i := 0; i < int(n); i++ {
		if val := int32(genData(i)); data[i] != val {
			t.Errorf("data at position %d is %d; expected %d\n", i, data[i], val)
		}
	}
}

func (ft *FileTest) getFloat(t *testing.T, v Var, n uint64) {
	data := make([]float32, n)
	if err := v.GetFloat(data); err != nil {
		t.Fatalf("GetInt failed: %v\n", err)
	}
	for i := 0; i < int(n); i++ {
		if val := float32(genData(i)); data[i] != val {
			t.Errorf("data at position %d is %d; expected %d\n", i, data[i], val)
		}
	}
}

var fileTests = []FileTest{
	{
		VarName:  "golang",
		DimNames: []string{"time", "growth"},
		DimLens:  []uint64{7, 3},
		DataType: NC_INT,
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
		if dims[i], err = f.PutDim(name, ft.DimLens[i]); err != nil {
			t.Fatalf("PutDim failed: %v\n", err)
		}
	}
	v, err := f.PutVar(ft.VarName, ft.DataType, dims)
	if err != nil {
		t.Fatalf("PutVar failed: %v\n", err)
	}
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
		d, err := f.GetDim(name)
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
	v, err := f.GetVar(ft.VarName)
	if err != nil {
		t.Errorf("GetVar failed: %v\n", err)
	}
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
