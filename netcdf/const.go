// Copyright 2014 The Go-NetCDF Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package netcdf

import "C"

const (
	_NC_NOERR    = 0
	_NC_MAX_NAME = 256
	_NC_GLOBAL   = -1 // Attribute id to put/get a global attribute
)

// FileMode represents a file's mode.
type FileMode C.int

// File modes for Open or Create
const (
	NC_DISKLESS FileMode = 0x0008 // use diskless file
	NC_MMAP     FileMode = 0x0010 // use diskless file with mmap
	NC_SHARE    FileMode = 0x0800 // share updates, limit cacheing
)

// File modes for Open
const (
	NC_NOWRITE FileMode = 0x0000 // set read-only access
	NC_WRITE   FileMode = 0x0001 // set read-write access
)

// File modes for Create
const (
	NC_CLOBBER       FileMode = 0x0000 // destroy existing file
	NC_NOCLOBBER     FileMode = 0x0004 // don't destroy existing file
	NC_CLASSIC_MODEL FileMode = 0x0100 // enforce classic model
	NC_NETCDF4       FileMode = 0x1000 // use netCDF-4/HDF5 format
	NC_64BIT_OFFSET  FileMode = 0x0200 // use large (64-bit) file offsets
)

// Type is a netCDF external data type.
type Type C.int

const (
	NC_BYTE   Type = 1      // signed 1 byte integer
	NC_CHAR   Type = 2      // ISO/ASCII character
	NC_SHORT  Type = 3      // signed 2 byte integer
	NC_INT    Type = 4      // signed 4 byte integer
	NC_LONG   Type = NC_INT // deprecated, but required for backward compatibility.
	NC_FLOAT  Type = 5      // single precision floating point number
	NC_DOUBLE Type = 6      // double precision floating point number
	NC_UBYTE  Type = 7      // unsigned 1 byte int
	NC_USHORT Type = 8      // unsigned 2-byte int
	NC_UINT   Type = 9      // unsigned 4-byte int
	NC_INT64  Type = 10     // signed 8-byte int
	NC_UINT64 Type = 11     // unsigned 8-byte int
	NC_STRING Type = 12     // string
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
