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
	NC_SHARE FileMode = C.NC_SHARE // share updates, limit cacheing
)

// File modes for Open
const (
	NC_NOWRITE FileMode = C.NC_NOWRITE // set read-only access
	NC_WRITE   FileMode = C.NC_WRITE   // set read-write access
)

// File modes for Create
const (
	NC_CLOBBER       FileMode = C.NC_CLOBBER       // destroy existing file
	NC_NOCLOBBER     FileMode = C.NC_NOCLOBBER     // don't destroy existing file
	NC_CLASSIC_MODEL FileMode = C.NC_CLASSIC_MODEL // enforce classic model
	NC_NETCDF4       FileMode = C.NC_NETCDF4       // use netCDF-4/HDF5 format
	NC_64BIT_OFFSET  FileMode = C.NC_64BIT_OFFSET  // use large (64-bit) file offsets
)

// Type is a netCDF external data type.
type Type C.nc_type

const (
	NC_BYTE   Type = C.NC_BYTE   // signed 1 byte integer
	NC_CHAR   Type = C.NC_CHAR   // ISO/ASCII character
	NC_SHORT  Type = C.NC_SHORT  // signed 2 byte integer
	NC_INT    Type = C.NC_INT    // signed 4 byte integer
	NC_LONG   Type = C.NC_LONG   // deprecated, but required for backward compatibility.
	NC_FLOAT  Type = C.NC_FLOAT  // single precision floating point number
	NC_DOUBLE Type = C.NC_DOUBLE // double precision floating point number
	NC_UBYTE  Type = C.NC_UBYTE  // unsigned 1 byte int
	NC_USHORT Type = C.NC_USHORT // unsigned 2-byte int
	NC_UINT   Type = C.NC_UINT   // unsigned 4-byte int
	NC_INT64  Type = C.NC_INT64  // signed 8-byte int
	NC_UINT64 Type = C.NC_UINT64 // unsigned 8-byte int
	NC_STRING Type = C.NC_STRING // string
)

var typeNames map[Type]string = map[Type]string{
	NC_BYTE:   "NC_BYTE",
	NC_CHAR:   "NC_CHAR",
	NC_SHORT:  "NC_SHORT",
	NC_INT:    "NC_INT",
	NC_FLOAT:  "NC_FLOAT",
	NC_DOUBLE: "NC_DOUBLE",
	NC_UBYTE:  "NC_UBYTE",
	NC_USHORT: "NC_USHORT",
	NC_UINT:   "NC_UINT",
	NC_INT64:  "NC_INT64",
	NC_UINT64: "NC_UINT64",
	NC_STRING: "NC_STRING",
}
