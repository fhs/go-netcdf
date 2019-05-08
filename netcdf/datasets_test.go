// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package netcdf

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
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
		if name := a.Name(); name != key {
			t.Errorf("attribute name is %v; expected %v\n", name, key)
		}
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
	n, err := v.NAttrs()
	if err != nil {
		t.Fatalf("NAttrs failed: %v\n", err)
	}
	if n != len(ft.Attr) {
		t.Errorf("NAttrs is %v; expected %v\n", n, len(ft.Attr))
	}
	for key, val := range ft.Attr {
		a := v.Attr(key)
		if name := a.Name(); name != key {
			t.Errorf("attribute name is %v; expected %v\n", name, key)
		}
		typ, err := a.Type()
		if err != nil {
			t.Fatalf("getting data type of attribute %s failed: %v\n", key, err)
		}
		var q interface{}
		switch typ {
		default:
			t.Errorf("unexpected attribute type %v\n", typ)
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

func getFileTests() []FileTest {
	var tests []FileTest

	bases := []FileTest{
		FileTest{
			VarName:  "gopher",
			DimNames: []string{"height", "width"},
			DimLens:  []uint64{7, 3},
		},
		FileTest{
			VarName:  "gopher",
			DimNames: []string{"time", "height", "width"},
			DimLens:  []uint64{12, 7, 3},
		},
		FileTest{
			VarName:  "gopher",
			DimNames: []string{},
			DimLens:  []uint64{},
		},
	}
	attrs := []map[string]interface{}{
		nil,
		map[string]interface{}{
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
	}
	types := []Type{UINT64, INT64, DOUBLE, UINT, INT, FLOAT, USHORT, SHORT, UBYTE, BYTE, CHAR}
	for _, test := range bases {
		for _, attr := range attrs {
			test.Attr = attr
			for _, dataType := range types {
				test.DataType = dataType
				tests = append(tests, test)
			}
		}
	}
	return tests
}

func TestCreate(t *testing.T) {
	for _, ft := range getFileTests() {
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
	dims := make([]Dim, len(ft.DimNames))
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
		t.Fatalf("unexpected type %v\n", ft.DataType)
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
		t.Fatalf("unexpected type %v\n", ft.DataType)
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

func TestVersion(t *testing.T) {
	if ver := Version(); ver == "" {
		t.Errorf("Bad Version %q\n", ver)
	}
}

func TestProduct(t *testing.T) {
	results := []struct {
		expected uint64
		shape    []uint64
	}{
		{24, []uint64{4, 6}},
		{0, []uint64{0, 2, 6}},
		{8064, []uint64{4, 2, 4, 21, 12}},
	}

	for _, test := range results {
		if p := product(test.shape); p != test.expected {
			t.Errorf("Result of 'product(%v)' is %v, expected %v", test.shape, p, test.expected)
		}
	}

}

func TestValidUnravel(t *testing.T) {
	tests := []struct {
		idx      uint64
		shape    []uint64
		expected []uint64
	}{
		{4, []uint64{1, 5, 3}, []uint64{0, 1, 1}},
		{1021, []uint64{12, 421, 25}, []uint64{0, 40, 21}},
		{211, []uint64{24, 11, 5}, []uint64{3, 9, 1}},
		{12, []uint64{12, 5}, []uint64{2, 2}},
		{1412, []uint64{125231}, []uint64{1412}},
		{122, []uint64{1321, 123, 12}, []uint64{0, 10, 2}},
	}

	for _, test := range tests {
		result, err := UnravelIndex(test.idx, test.shape)
		if err != nil {
			t.Errorf("Got error while reading valid test, got %v", err)
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Result of 'UnravelIndex(%v, %v)' is %v, expected %v", test.idx, test.shape, result, test.expected)
		}
	}

}

func TestUnravelErrors(t *testing.T) {
	type IdxTests struct {
		idx      uint64
		shape    []uint64
		expected []uint64
	}

	var zeros = []IdxTests{
		{12, []uint64{3, 0}, nil},
		{2, []uint64{0, 3}, nil},
	}

	for _, test := range zeros {
		var _, err = UnravelIndex(test.idx, test.shape)
		if !strings.Contains(err.Error(), "0 encountered") {
			t.Errorf("Expected '0' error, got %v", err)
		}
	}

	var tooBig = []IdxTests{
		{241, []uint64{3, 5}, nil},
		{221, []uint64{2, 3}, nil},
	}

	for _, test := range tooBig {
		var _, err = UnravelIndex(test.idx, test.shape)
		if !strings.Contains(err.Error(), "> size") {
			t.Errorf("Expected '0' error, got %v", err)
		}
	}
}

func testWriteFileViaIdx(t *testing.T, filename string, ft *FileTest) {
	f, err := CreateFile(filename, CLOBBER|NETCDF4)
	if err != nil {
		t.Fatalf("Create failed: %v\n", err)
	}
	dims := make([]Dim, len(ft.DimNames))
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
		t.Fatalf("unexpected type %v\n", ft.DataType)
	case UINT64:
		err = testWriteUint64At(v, n)
	case INT64:
		err = testWriteInt64At(v, n)
	case DOUBLE:
		err = testWriteFloat64At(v, n)
	case UINT:
		err = testWriteUint32At(v, n)
	case INT:
		err = testWriteInt32At(v, n)
	case FLOAT:
		err = testWriteFloat32At(v, n)
	case USHORT:
		err = testWriteUint16At(v, n)
	case SHORT:
		err = testWriteInt16At(v, n)
	case UBYTE:
		err = testWriteUint8At(v, n)
	case BYTE:
		err = testWriteInt8At(v, n)
	case CHAR:
		err = testWriteBytesAt(v, n)
	}
	if err != nil {
		t.Errorf("%v: writing data failed: %v\n", ft.DataType, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
}

func testReadFileViaIdx(t *testing.T, filename string, ft *FileTest) {
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
		t.Fatalf("unexpected type %v\n", ft.DataType)
	case UINT64:
		err = testReadUint64At(v, n)
	case INT64:
		err = testReadInt64At(v, n)
	case DOUBLE:
		err = testReadFloat64At(v, n)
	case UINT:
		err = testReadUint32At(v, n)
	case INT:
		err = testReadInt32At(v, n)
	case FLOAT:
		err = testReadFloat32At(v, n)
	case USHORT:
		err = testReadUint16At(v, n)
	case SHORT:
		err = testReadInt16At(v, n)
	case UBYTE:
		err = testReadUint8At(v, n)
	case BYTE:
		err = testReadInt8At(v, n)
	case CHAR:
		err = testReadBytesAt(v, n)
	}
	if err != nil {
		t.Fatalf("reading data failed: %v\n", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
}

func TestAt(t *testing.T) {
	for _, ft := range getFileTests() {
		f, err := ioutil.TempFile("", "netcdf_test")
		if err != nil {
			t.Fatalf("creating temporary file failed: %v\n", err)
		}
		testWriteFileViaIdx(t, f.Name(), &ft)
		testReadFileViaIdx(t, f.Name(), &ft)
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("removing temporary file failed: %v\n", err)
		}
	}
}

func TestCompression(t *testing.T) {
	f, err := ioutil.TempFile("", "netcdf_test")
	if err != nil {
		t.Fatalf("creating temporary file failed: %v\n", err)
	}
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Errorf("removing temporary file failed: %v\n", err)
		}
	}()

	// Create a new NetCDF 4 file.
	ds, err := CreateFile(f.Name(), CLOBBER|NETCDF4)
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

	// Set and read compression parameters.
	tests := []struct {
		shuffle bool
		deflate bool
		level   int
	}{
		{true, true, 1},
		{true, true, 5},
		{true, true, 9},
		{false, true, 0},
		{true, false, 0},
	}

	for _, test := range tests {
		err = v.SetCompression(test.shuffle, test.deflate, test.level)
		if err != nil {
			t.Fatalf("setting compression failed: %v\n", err)
		}

		shuffle, deflate, level, err := v.Compression()
		if err != nil {
			t.Fatalf("reading compression failed: %v\n", err)
		}

		if shuffle != test.shuffle || deflate != test.deflate || level != test.level {
			t.Errorf("result of SetCompression(%t, %t, %d) is (%t, %t, %d), expected (%t, %t, %d)\n",
				test.shuffle, test.deflate, test.level,
				shuffle, deflate, level,
				test.shuffle, test.deflate, test.level,
			)
		}
	}
}
