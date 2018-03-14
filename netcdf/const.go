// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package netcdf

// #include <netcdf.h>
import "C"

// FileMode represents a file's mode.
type FileMode C.int

// File modes for Open or Create
const (
	SHARE FileMode = C.NC_SHARE // share updates, limit cacheing
)

// File modes for Open
const (
	NOWRITE FileMode = C.NC_NOWRITE // set read-only access
	WRITE   FileMode = C.NC_WRITE   // set read-write access
)

// File modes for Create
const (
	CLOBBER       FileMode = C.NC_CLOBBER       // destroy existing file
	NOCLOBBER     FileMode = C.NC_NOCLOBBER     // don't destroy existing file
	CLASSIC_MODEL FileMode = C.NC_CLASSIC_MODEL // enforce classic model
	NETCDF4       FileMode = C.NC_NETCDF4       // use netCDF-4/HDF5 format
	OFFSET_64BIT  FileMode = C.NC_64BIT_OFFSET  // use large (64-bit) file offsets
)

// Type is a netCDF external data type.
type Type C.nc_type

// Type declarations according to C standards
const (
	BYTE   Type = C.NC_BYTE   // signed 1 byte integer
	CHAR   Type = C.NC_CHAR   // ISO/ASCII character
	SHORT  Type = C.NC_SHORT  // signed 2 byte integer
	INT    Type = C.NC_INT    // signed 4 byte integer
	LONG   Type = C.NC_LONG   // deprecated, but required for backward compatibility.
	FLOAT  Type = C.NC_FLOAT  // single precision floating point number
	DOUBLE Type = C.NC_DOUBLE // double precision floating point number
	UBYTE  Type = C.NC_UBYTE  // unsigned 1 byte int
	USHORT Type = C.NC_USHORT // unsigned 2-byte int
	UINT   Type = C.NC_UINT   // unsigned 4-byte int
	INT64  Type = C.NC_INT64  // signed 8-byte int
	UINT64 Type = C.NC_UINT64 // unsigned 8-byte int
	STRING Type = C.NC_STRING // string
)

var typeNames = map[Type]string{
	BYTE:   "BYTE",
	CHAR:   "CHAR",
	SHORT:  "SHORT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	DOUBLE: "DOUBLE",
	UBYTE:  "UBYTE",
	USHORT: "USHORT",
	UINT:   "UINT",
	INT64:  "INT64",
	UINT64: "UINT64",
	STRING: "STRING",
}

// String converts a Type to its string representation.
func (t Type) String() string {
	return typeNames[t]
}
