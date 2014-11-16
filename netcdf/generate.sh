#!/bin/sh

cat nc_double.go |
	gofmt -r 'float64 -> uint64' |
	gofmt -r 'C.double -> C.ulonglong' |
	gofmt -r 'DOUBLE -> UINT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ulonglong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ulonglong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ulonglong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ulonglong' |
	sed 's/Float64sReader/Uint64sReader/' |
	sed 's/GetFloat64s/GetUint64s/' |
	sed 's/WriteFloat64s/WriteUint64s/' |
	sed 's/ReadFloat64s/ReadUint64s/' \
	> nc_uint64.go

cat nc_double.go |
	gofmt -r 'float64 -> int64' |
	gofmt -r 'C.double -> C.longlong' |
	gofmt -r 'DOUBLE -> INT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_longlong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_longlong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_longlong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_longlong' |
	sed 's/Float64sReader/Int64sReader/' |
	sed 's/GetFloat64s/GetInt64s/' |
	sed 's/WriteFloat64s/WriteInt64s/' |
	sed 's/ReadFloat64s/ReadInt64s/' \
	> nc_int64.go

cat nc_double.go |
	gofmt -r 'float64 -> uint32' |
	gofmt -r 'C.double -> C.uint' |
	gofmt -r 'DOUBLE -> UINT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uint' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uint' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uint' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uint' |
	sed 's/Float64sReader/Uint32sReader/' |
	sed 's/GetFloat64s/GetUint32s/' |
	sed 's/WriteFloat64s/WriteUint32s/' |
	sed 's/ReadFloat64s/ReadUint32s/' \
	> nc_uint.go

cat nc_double.go |
	gofmt -r 'float64 -> int32' |
	gofmt -r 'C.double -> C.int' |
	gofmt -r 'DOUBLE -> INT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_int' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_int' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_int' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_int' |
	sed 's/Float64sReader/Int32sReader/' |
	sed 's/GetFloat64s/GetInt32s/' |
	sed 's/WriteFloat64s/WriteInt32s/' |
	sed 's/ReadFloat64s/ReadInt32s/' \
	> nc_int.go

cat nc_double.go |
	gofmt -r 'float64 -> float32' |
	gofmt -r 'C.double -> C.float' |
	gofmt -r 'DOUBLE -> FLOAT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_float' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_float' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_float' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_float' |
	sed 's/Float64sReader/Float32sReader/' |
	sed 's/GetFloat64s/GetFloat32s/' |
	sed 's/WriteFloat64s/WriteFloat32s/' |
	sed 's/ReadFloat64s/ReadFloat32s/' \
	> nc_float.go

cat nc_double.go |
	gofmt -r 'float64 -> uint16' |
	gofmt -r 'C.double -> C.ushort' |
	gofmt -r 'DOUBLE -> USHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ushort' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ushort' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ushort' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ushort' |
	sed 's/Float64sReader/Uint16sReader/' |
	sed 's/GetFloat64s/GetUint16s/' |
	sed 's/WriteFloat64s/WriteUint16s/' |
	sed 's/ReadFloat64s/ReadUint16s/' \
	> nc_ushort.go

cat nc_double.go |
	gofmt -r 'float64 -> int16' |
	gofmt -r 'C.double -> C.short' |
	gofmt -r 'DOUBLE -> SHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_short' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_short' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_short' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_short' |
	sed 's/Float64sReader/Int16sReader/' |
	sed 's/GetFloat64s/GetInt16s/' |
	sed 's/WriteFloat64s/WriteInt16s/' |
	sed 's/ReadFloat64s/ReadInt16s/' \
	> nc_short.go

cat nc_double.go |
	gofmt -r 'float64 -> uint8' |
	gofmt -r 'C.double -> C.uchar' |
	gofmt -r 'DOUBLE -> UBYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uchar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uchar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uchar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uchar' |
	sed 's/Float64sReader/Uint8sReader/' |
	sed 's/GetFloat64s/GetUint8s/' |
	sed 's/WriteFloat64s/WriteUint8s/' |
	sed 's/ReadFloat64s/ReadUint8s/' \
	> nc_ubyte.go

cat nc_double.go |
	gofmt -r 'float64 -> int8' |
	gofmt -r 'C.double -> C.schar' |
	gofmt -r 'DOUBLE -> BYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_schar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_schar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_schar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_schar' |
	sed 's/Float64sReader/Int8sReader/' |
	sed 's/GetFloat64s/GetInt8s/' |
	sed 's/WriteFloat64s/WriteInt8s/' |
	sed 's/ReadFloat64s/ReadInt8s/' \
	> nc_byte.go


# We return []byte (i.e. []uint8) for CHAR Type because:
#	- Returning string would not be very flexible, since '\0' characters
#	  may or may not require trimming.
#	- Returning []rune (i.e. []int32) takes up more space and we know
#	  we're limited to ASCII.
#	- Any other types can't be easily converted to a string
#	  (e.g. string([]int8) does not work)
cat nc_double.go |
	gofmt -r 'float64 -> byte' |
	gofmt -r 'C.double -> C.char' |
	gofmt -r 'DOUBLE -> CHAR' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_text' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_text' |
	gofmt -r 'C.nc_put_att_double(a, b, c, d, e, f) -> C.nc_put_att_text(a, b, c, e, f)' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_text' |
	sed 's/Float64sReader/BytesReader/' |
	sed 's/GetFloat64s/GetBytes/' |
	sed 's/WriteFloat64s/WriteBytes/' |
	sed 's/ReadFloat64s/ReadBytes/' \
	> nc_char.go
