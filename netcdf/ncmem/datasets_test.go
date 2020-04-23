// Copyright 2020 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package ncmem

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/fhs/go-netcdf/netcdf"
)

func TestOpenReader(t *testing.T) {
	const expectedYear = 2012

	f, err := os.Open("../gopher.nc")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	ds, err := OpenReader("gopher.nc", netcdf.NOWRITE, 0, f)
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()

	v, err := ds.Var("gopher")
	if err != nil {
		t.Fatal(err)
	}

	year, err := netcdf.GetInt32s(v.Attr("year"))
	if err != nil {
		t.Fatal(err)
	}
	if year[0] != expectedYear {
		t.Errorf("year is %d; expected %d", year[0], expectedYear)
	}
}

type FileTest struct {
	VarName  string
	DimNames []string
	DimLens  []uint64
	DataType netcdf.Type
	Attr     map[string]interface{}
}

func (ft *FileTest) putAttrs(t *testing.T, v netcdf.Var) {
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

func (ft *FileTest) getAttrs(t *testing.T, v netcdf.Var) {
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
		case netcdf.UINT64:
			q, err = netcdf.GetUint64s(a)
		case netcdf.INT64:
			q, err = netcdf.GetInt64s(a)
		case netcdf.DOUBLE:
			q, err = netcdf.GetFloat64s(a)
		case netcdf.UINT:
			q, err = netcdf.GetUint32s(a)
		case netcdf.INT:
			q, err = netcdf.GetInt32s(a)
		case netcdf.FLOAT:
			q, err = netcdf.GetFloat32s(a)
		case netcdf.USHORT:
			q, err = netcdf.GetUint16s(a)
		case netcdf.SHORT:
			q, err = netcdf.GetInt16s(a)
		case netcdf.UBYTE:
			q, err = netcdf.GetUint8s(a)
		case netcdf.BYTE:
			q, err = netcdf.GetInt8s(a)
		case netcdf.CHAR:
			var b []byte
			b, err = netcdf.GetBytes(a)
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
	types := []netcdf.Type{netcdf.UINT64, netcdf.INT64, netcdf.DOUBLE, netcdf.UINT, netcdf.INT, netcdf.FLOAT, netcdf.USHORT, netcdf.SHORT, netcdf.UBYTE, netcdf.BYTE, netcdf.CHAR}
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
		data := testCreate(t, &ft)
		testReadOpen(t, &ft, data)
		testReadOpenReader(t, &ft, data)
		testReadOpenLenReader(t, &ft, data)
	}
}

func TestCreateCBytes(t *testing.T) {
	for _, ft := range getFileTests() {
		data := testCreateBytes(t, &ft)
		testReadOpen(t, &ft, data.Data)
		testReadOpenReader(t, &ft, data.Data)
		testReadOpenLenReader(t, &ft, data.Data)
		data.Free()
	}
}

func testCreate(t *testing.T, ft *FileTest) []byte {
	f, err := Create("netcdf_test", netcdf.CLOBBER|netcdf.NETCDF4, 0)
	if err != nil {
		t.Fatalf("Create failed: %v\n", err)
	}

	testCreateData(t, ft, f)

	data, err := f.CloseCopyBytes()
	if err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
	return data
}

func testCreateBytes(t *testing.T, ft *FileTest) *Bytes {
	f, err := Create("netcdf_test", netcdf.CLOBBER|netcdf.NETCDF4, 0)
	if err != nil {
		t.Fatalf("Create failed: %v\n", err)
	}

	testCreateData(t, ft, f)

	data, err := f.CloseBytes()
	if err != nil {
		data.Free()
		t.Fatalf("Close failed: %v\n", err)
	}
	return data
}

func testCreateData(t *testing.T, ft *FileTest, f Dataset) {
	dims := make([]netcdf.Dim, len(ft.DimNames))
	var err error
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
	case netcdf.UINT64:
		err = testWriteUint64s(v, n)
	case netcdf.INT64:
		err = testWriteInt64s(v, n)
	case netcdf.DOUBLE:
		err = testWriteFloat64s(v, n)
	case netcdf.UINT:
		err = testWriteUint32s(v, n)
	case netcdf.INT:
		err = testWriteInt32s(v, n)
	case netcdf.FLOAT:
		err = testWriteFloat32s(v, n)
	case netcdf.USHORT:
		err = testWriteUint16s(v, n)
	case netcdf.SHORT:
		err = testWriteInt16s(v, n)
	case netcdf.UBYTE:
		err = testWriteUint8s(v, n)
	case netcdf.BYTE:
		err = testWriteInt8s(v, n)
	case netcdf.CHAR:
		err = testWriteBytes(v, n)
	}
	if err != nil {
		t.Errorf("writing data failed: %v\n", err)
	}
}

func testReadOpen(t *testing.T, ft *FileTest, data []byte) {
	f, err := Open("netdf_test", netcdf.NOWRITE, 0, data)
	if err != nil {
		t.Fatalf("Open failed: %v\n", err)
	}

	testReadData(t, ft, f)
}

// readerOnly is an io.Reader which does not fulfill other
// interfaces, such as LenReader.
type readerOnly struct {
	r io.Reader
}

func (ro readerOnly) Read(p []byte) (int, error) {
	return ro.r.Read(p)
}

func testReadOpenReader(t *testing.T, ft *FileTest, data []byte) {
	r := readerOnly{bytes.NewReader(data)}

	f, err := OpenReader("netdf_test", netcdf.NOWRITE, 0, r)
	if err != nil {
		t.Fatalf("Open failed: %v\n", err)
	}

	testReadData(t, ft, f)
}

func testReadOpenLenReader(t *testing.T, ft *FileTest, data []byte) {
	r := bytes.NewReader(data)

	f, err := OpenReader("netdf_test", netcdf.NOWRITE, 0, r)
	if err != nil {
		t.Fatalf("Open failed: %v\n", err)
	}

	testReadData(t, ft, f)
}

func testReadData(t *testing.T, ft *FileTest, f Dataset) {
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
	case netcdf.UINT64:
		err = testReadUint64s(v, n)
	case netcdf.INT64:
		err = testReadInt64s(v, n)
	case netcdf.DOUBLE:
		err = testReadFloat64s(v, n)
	case netcdf.UINT:
		err = testReadUint32s(v, n)
	case netcdf.INT:
		err = testReadInt32s(v, n)
	case netcdf.FLOAT:
		err = testReadFloat32s(v, n)
	case netcdf.USHORT:
		err = testReadUint16s(v, n)
	case netcdf.SHORT:
		err = testReadInt16s(v, n)
	case netcdf.UBYTE:
		err = testReadUint8s(v, n)
	case netcdf.BYTE:
		err = testReadInt8s(v, n)
	case netcdf.CHAR:
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
	_, err := Open("netcdf_test", netcdf.NOWRITE, 0, nil)
	if err == nil {
		t.Fatalf("Opened nil data\n")
	}
	if len(err.Error()) == 0 {
		t.Errorf("empty error\n")
	}
}

func TestEndDef(t *testing.T) {
	// Create a new NetCDF 3 file.
	ds, err := Create("netcdf_test", 0, 0)
	if err != nil {
		t.Fatalf("creating file failed: %v\n", err)
	}
	defer ds.Close()

	size, err := ds.AddDim("size", 5)
	if err != nil {
		t.Fatalf("adding dimension failed: %v\n", err)
	}
	v, err := ds.AddVar("gopher", netcdf.INT, []netcdf.Dim{size})
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

func testWriteFileViaIdx(t *testing.T, ft *FileTest) []byte {
	f, err := Create("netcdf_test", netcdf.NETCDF4, 0)
	if err != nil {
		t.Fatalf("Create failed: %v\n", err)
	}
	dims := make([]netcdf.Dim, len(ft.DimNames))
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
	case netcdf.UINT64:
		err = testWriteUint64At(v, n)
	case netcdf.INT64:
		err = testWriteInt64At(v, n)
	case netcdf.DOUBLE:
		err = testWriteFloat64At(v, n)
	case netcdf.UINT:
		err = testWriteUint32At(v, n)
	case netcdf.INT:
		err = testWriteInt32At(v, n)
	case netcdf.FLOAT:
		err = testWriteFloat32At(v, n)
	case netcdf.USHORT:
		err = testWriteUint16At(v, n)
	case netcdf.SHORT:
		err = testWriteInt16At(v, n)
	case netcdf.UBYTE:
		err = testWriteUint8At(v, n)
	case netcdf.BYTE:
		err = testWriteInt8At(v, n)
	case netcdf.CHAR:
		err = testWriteBytesAt(v, n)
	}
	if err != nil {
		t.Errorf("%v: writing data failed: %v\n", ft.DataType, err)
	}
	data, err := f.CloseCopyBytes()
	if err != nil {
		t.Fatalf("Close failed: %v\n", err)
	}
	return data
}

func testReadFileViaIdx(t *testing.T, ft *FileTest, data []byte) {
	f, err := Open("netcdf_test", netcdf.NOWRITE, 0, data)
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
	case netcdf.UINT64:
		err = testReadUint64At(v, n)
	case netcdf.INT64:
		err = testReadInt64At(v, n)
	case netcdf.DOUBLE:
		err = testReadFloat64At(v, n)
	case netcdf.UINT:
		err = testReadUint32At(v, n)
	case netcdf.INT:
		err = testReadInt32At(v, n)
	case netcdf.FLOAT:
		err = testReadFloat32At(v, n)
	case netcdf.USHORT:
		err = testReadUint16At(v, n)
	case netcdf.SHORT:
		err = testReadInt16At(v, n)
	case netcdf.UBYTE:
		err = testReadUint8At(v, n)
	case netcdf.BYTE:
		err = testReadInt8At(v, n)
	case netcdf.CHAR:
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
		data := testWriteFileViaIdx(t, &ft)
		testReadFileViaIdx(t, &ft, data)
	}
}

func TestCompression(t *testing.T) {
	// Create a new NetCDF 4 file.
	ds, err := Create("netcdf_test", netcdf.CLOBBER|netcdf.NETCDF4, 0)
	if err != nil {
		t.Fatalf("creating file failed: %v\n", err)
	}
	defer ds.Close()

	size, err := ds.AddDim("size", 5)
	if err != nil {
		t.Fatalf("adding dimension failed: %v\n", err)
	}
	v, err := ds.AddVar("gopher", netcdf.INT, []netcdf.Dim{size})
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

// testReadUint64s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteUint64s(v netcdf.Var, n uint64) error {
	data := make([]uint64, n)
	for i := 0; i < int(n); i++ {
		data[i] = uint64(i + 10)
	}
	return v.WriteUint64s(data)
}

// testReadUint64s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadUint64s(v netcdf.Var, n uint64) error {
	data := make([]uint64, n)
	if err := v.ReadUint64s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := uint64(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadUint64At(v netcdf.Var, n uint64) error {
	data := make([]uint64, n)
	if err := v.ReadUint64s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadUint64At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteUint64At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteUint64At(coord, uint64(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadUint64At(coord)
		if val != uint64(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadInt64s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteInt64s(v netcdf.Var, n uint64) error {
	data := make([]int64, n)
	for i := 0; i < int(n); i++ {
		data[i] = int64(i + 10)
	}
	return v.WriteInt64s(data)
}

// testReadInt64s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadInt64s(v netcdf.Var, n uint64) error {
	data := make([]int64, n)
	if err := v.ReadInt64s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := int64(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadInt64At(v netcdf.Var, n uint64) error {
	data := make([]int64, n)
	if err := v.ReadInt64s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadInt64At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteInt64At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteInt64At(coord, int64(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadInt64At(coord)
		if val != int64(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadFloat64s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteFloat64s(v netcdf.Var, n uint64) error {
	data := make([]float64, n)
	for i := 0; i < int(n); i++ {
		data[i] = float64(i + 10)
	}
	return v.WriteFloat64s(data)
}

// testReadFloat64s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadFloat64s(v netcdf.Var, n uint64) error {
	data := make([]float64, n)
	if err := v.ReadFloat64s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := float64(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadFloat64At(v netcdf.Var, n uint64) error {
	data := make([]float64, n)
	if err := v.ReadFloat64s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadFloat64At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteFloat64At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteFloat64At(coord, float64(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadFloat64At(coord)
		if val != float64(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadUint32s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteUint32s(v netcdf.Var, n uint64) error {
	data := make([]uint32, n)
	for i := 0; i < int(n); i++ {
		data[i] = uint32(i + 10)
	}
	return v.WriteUint32s(data)
}

// testReadUint32s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadUint32s(v netcdf.Var, n uint64) error {
	data := make([]uint32, n)
	if err := v.ReadUint32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := uint32(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadUint32At(v netcdf.Var, n uint64) error {
	data := make([]uint32, n)
	if err := v.ReadUint32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadUint32At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteUint32At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteUint32At(coord, uint32(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadUint32At(coord)
		if val != uint32(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadInt32s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteInt32s(v netcdf.Var, n uint64) error {
	data := make([]int32, n)
	for i := 0; i < int(n); i++ {
		data[i] = int32(i + 10)
	}
	return v.WriteInt32s(data)
}

// testReadInt32s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadInt32s(v netcdf.Var, n uint64) error {
	data := make([]int32, n)
	if err := v.ReadInt32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := int32(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadInt32At(v netcdf.Var, n uint64) error {
	data := make([]int32, n)
	if err := v.ReadInt32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadInt32At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteInt32At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteInt32At(coord, int32(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadInt32At(coord)
		if val != int32(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadFloat32s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteFloat32s(v netcdf.Var, n uint64) error {
	data := make([]float32, n)
	for i := 0; i < int(n); i++ {
		data[i] = float32(i + 10)
	}
	return v.WriteFloat32s(data)
}

// testReadFloat32s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadFloat32s(v netcdf.Var, n uint64) error {
	data := make([]float32, n)
	if err := v.ReadFloat32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := float32(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadFloat32At(v netcdf.Var, n uint64) error {
	data := make([]float32, n)
	if err := v.ReadFloat32s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadFloat32At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteFloat32At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteFloat32At(coord, float32(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadFloat32At(coord)
		if val != float32(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadUint16s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteUint16s(v netcdf.Var, n uint64) error {
	data := make([]uint16, n)
	for i := 0; i < int(n); i++ {
		data[i] = uint16(i + 10)
	}
	return v.WriteUint16s(data)
}

// testReadUint16s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadUint16s(v netcdf.Var, n uint64) error {
	data := make([]uint16, n)
	if err := v.ReadUint16s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := uint16(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadUint16At(v netcdf.Var, n uint64) error {
	data := make([]uint16, n)
	if err := v.ReadUint16s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadUint16At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteUint16At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteUint16At(coord, uint16(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadUint16At(coord)
		if val != uint16(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadInt16s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteInt16s(v netcdf.Var, n uint64) error {
	data := make([]int16, n)
	for i := 0; i < int(n); i++ {
		data[i] = int16(i + 10)
	}
	return v.WriteInt16s(data)
}

// testReadInt16s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadInt16s(v netcdf.Var, n uint64) error {
	data := make([]int16, n)
	if err := v.ReadInt16s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := int16(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadInt16At(v netcdf.Var, n uint64) error {
	data := make([]int16, n)
	if err := v.ReadInt16s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadInt16At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteInt16At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteInt16At(coord, int16(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadInt16At(coord)
		if val != int16(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadUint8s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteUint8s(v netcdf.Var, n uint64) error {
	data := make([]uint8, n)
	for i := 0; i < int(n); i++ {
		data[i] = uint8(i + 10)
	}
	return v.WriteUint8s(data)
}

// testReadUint8s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadUint8s(v netcdf.Var, n uint64) error {
	data := make([]uint8, n)
	if err := v.ReadUint8s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := uint8(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadUint8At(v netcdf.Var, n uint64) error {
	data := make([]uint8, n)
	if err := v.ReadUint8s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadUint8At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteUint8At(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteUint8At(coord, uint8(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadUint8At(coord)
		if val != uint8(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}

// testReadInt8s writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteInt8s(v netcdf.Var, n uint64) error {
	data := make([]int8, n)
	for i := 0; i < int(n); i++ {
		data[i] = int8(i + 10)
	}
	return v.WriteInt8s(data)
}

// testReadInt8s reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadInt8s(v netcdf.Var, n uint64) error {
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

func testReadInt8At(v netcdf.Var, n uint64) error {
	data := make([]int8, n)
	if err := v.ReadInt8s(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadInt8At(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteInt8At(v netcdf.Var, n uint64) error {
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

// testReadBytes writes somes data to v. N is v.Len().
// This function is only used for testing.
func testWriteBytes(v netcdf.Var, n uint64) error {
	data := make([]byte, n)
	for i := 0; i < int(n); i++ {
		data[i] = byte(i + 10)
	}
	return v.WriteBytes(data)
}

// testReadBytes reads data from v and checks that it's the same as what
// was written by testWriteDouble. N is v.Len().
// This function is only used for testing.
func testReadBytes(v netcdf.Var, n uint64) error {
	data := make([]byte, n)
	if err := v.ReadBytes(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if val := byte(i + 10); data[i] != val {
			return fmt.Errorf("data at position %d is %v; expected %v", i, data[i], val)
		}
	}
	return nil
}

func testReadBytesAt(v netcdf.Var, n uint64) error {
	data := make([]byte, n)
	if err := v.ReadBytes(data); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		shape, _ := v.LenDims()
		coords, _ := netcdf.UnravelIndex(uint64(i), shape)
		expected := data[i]
		val, _ := v.ReadBytesAt(coords)
		if val != expected {
			return fmt.Errorf("data at position %v is %v; expected %v", i, val, expected)
		}
	}
	return nil
}

func testWriteBytesAt(v netcdf.Var, n uint64) error {
	shape, _ := v.LenDims()
	ndim := len(shape)
	coord := make([]uint64, ndim)
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		v.WriteBytesAt(coord, byte(i))
	}
	for i := 0; i < ndim; i++ {
		for k := 0; k < ndim; k++ {
			coord[k] = uint64(i)
		}
		val, _ := v.ReadBytesAt(coord)
		if val != byte(i) {
			return fmt.Errorf("data at position %v is %v; expected %v", coord, val, int(i))
		}
	}
	return nil
}
