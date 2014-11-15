#!/bin/sh

cat nc_double.go |
	gofmt -r 'float64 -> uint64' |
	gofmt -r 'C.double -> C.ulonglong' |
	gofmt -r 'NC_DOUBLE -> NC_UINT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ulonglong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ulonglong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ulonglong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ulonglong' |
	sed 's/DoubleReader/Uint64Reader/' |
	sed 's/GetDouble/GetUint64/' |
	sed 's/WriteDouble/WriteUint64/' |
	sed 's/ReadDouble/ReadUint64/' \
	> nc_uint64.go

cat nc_double.go |
	gofmt -r 'float64 -> int64' |
	gofmt -r 'C.double -> C.longlong' |
	gofmt -r 'NC_DOUBLE -> NC_INT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_longlong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_longlong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_longlong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_longlong' |
	sed 's/DoubleReader/Int64Reader/' |
	sed 's/GetDouble/GetInt64/' |
	sed 's/WriteDouble/WriteInt64/' |
	sed 's/ReadDouble/ReadInt64/' \
	> nc_int64.go

cat nc_double.go |
	gofmt -r 'float64 -> uint32' |
	gofmt -r 'C.double -> C.uint' |
	gofmt -r 'NC_DOUBLE -> NC_UINT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uint' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uint' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uint' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uint' |
	sed 's/DoubleReader/UintReader/' |
	sed 's/GetDouble/GetUint/' |
	sed 's/WriteDouble/WriteUint/' |
	sed 's/ReadDouble/ReadUint/' \
	> nc_uint.go

cat nc_double.go |
	gofmt -r 'float64 -> int32' |
	gofmt -r 'C.double -> C.int' |
	gofmt -r 'NC_DOUBLE -> NC_INT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_int' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_int' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_int' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_int' |
	sed 's/DoubleReader/IntReader/' |
	sed 's/GetDouble/GetInt/' |
	sed 's/WriteDouble/WriteInt/' |
	sed 's/ReadDouble/ReadInt/' \
	> nc_int.go

cat nc_double.go |
	gofmt -r 'float64 -> float32' |
	gofmt -r 'C.double -> C.float' |
	gofmt -r 'NC_DOUBLE -> NC_FLOAT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_float' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_float' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_float' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_float' |
	sed 's/DoubleReader/FloatReader/' |
	sed 's/GetDouble/GetFloat/' |
	sed 's/WriteDouble/WriteFloat/' |
	sed 's/ReadDouble/ReadFloat/' \
	> nc_float.go

cat nc_double.go |
	gofmt -r 'float64 -> uint16' |
	gofmt -r 'C.double -> C.ushort' |
	gofmt -r 'NC_DOUBLE -> NC_USHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ushort' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ushort' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ushort' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ushort' |
	sed 's/DoubleReader/UshortReader/' |
	sed 's/GetDouble/GetUshort/' |
	sed 's/WriteDouble/WriteUshort/' |
	sed 's/ReadDouble/ReadUshort/' \
	> nc_ushort.go

cat nc_double.go |
	gofmt -r 'float64 -> int16' |
	gofmt -r 'C.double -> C.short' |
	gofmt -r 'NC_DOUBLE -> NC_SHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_short' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_short' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_short' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_short' |
	sed 's/DoubleReader/ShortReader/' |
	sed 's/GetDouble/GetShort/' |
	sed 's/WriteDouble/WriteShort/' |
	sed 's/ReadDouble/ReadShort/' \
	> nc_short.go

cat nc_double.go |
	gofmt -r 'float64 -> uint8' |
	gofmt -r 'C.double -> C.uchar' |
	gofmt -r 'NC_DOUBLE -> NC_UBYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uchar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uchar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uchar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uchar' |
	sed 's/DoubleReader/UbyteReader/' |
	sed 's/GetDouble/GetUbyte/' |
	sed 's/WriteDouble/WriteUbyte/' |
	sed 's/ReadDouble/ReadUbyte/' \
	> nc_ubyte.go

cat nc_double.go |
	gofmt -r 'float64 -> int8' |
	gofmt -r 'C.double -> C.schar' |
	gofmt -r 'NC_DOUBLE -> NC_BYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_schar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_schar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_schar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_schar' |
	sed 's/DoubleReader/ByteReader/' |
	sed 's/GetDouble/GetByte/' |
	sed 's/WriteDouble/WriteByte/' |
	sed 's/ReadDouble/ReadByte/' \
	> nc_byte.go


# We return []byte (i.e. []uint8) for NC_CHAR Type because:
#	- Returning string would not be very flexible, since '\0' characters
#	  may or may not require trimming.
#	- Returning []rune (i.e. []int32) takes up more space and we know
#	  we're limited to ASCII.
#	- Any other types can't be easily converted to a string
#	  (e.g. string([]int8) does not work)
cat nc_double.go |
	gofmt -r 'float64 -> byte' |
	gofmt -r 'C.double -> C.char' |
	gofmt -r 'NC_DOUBLE -> NC_CHAR' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_text' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_text' |
	gofmt -r 'C.nc_put_att_double(a, b, c, d, e, f) -> C.nc_put_att_text(a, b, c, e, f)' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_text' |
	sed 's/DoubleReader/CharReader/' |
	sed 's/GetDouble/GetChar/' |
	sed 's/WriteDouble/WriteChar/' |
	sed 's/ReadDouble/ReadChar/' \
	> nc_char.go
