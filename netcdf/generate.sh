#!/bin/sh

cat nc_double.go |
	gofmt -r 'float64 -> float32' |
	gofmt -r 'C.double -> C.float' |
	gofmt -r 'NC_DOUBLE -> NC_FLOAT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_float' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_float' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_float' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_float' |
	gofmt -r 'WriteDouble -> WriteFloat' |
	gofmt -r 'ReadDouble -> ReadFloat' |
	sed 's_^// WriteDouble _// WriteFloat _' |
	sed 's_^// ReadDouble _// ReadFloat _' \
	> nc_float.go

cat nc_double.go |
	gofmt -r 'float64 -> uint64' |
	gofmt -r 'C.double -> C.ulonglong' |
	gofmt -r 'NC_DOUBLE -> NC_UINT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ulonglong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ulonglong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ulonglong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ulonglong' |
	gofmt -r 'WriteDouble -> WriteUint64' |
	gofmt -r 'ReadDouble -> ReadUint64' |
	sed 's_^// WriteDouble _// WriteUint64 _' |
	sed 's_^// ReadDouble _// ReadUint64 _' \
	> nc_uint64.go

cat nc_double.go |
	gofmt -r 'float64 -> int64' |
	gofmt -r 'C.double -> C.longlong' |
	gofmt -r 'NC_DOUBLE -> NC_INT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_longlong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_longlong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_longlong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_longlong' |
	gofmt -r 'WriteDouble -> WriteInt64' |
	gofmt -r 'ReadDouble -> ReadInt64' |
	sed 's_^// WriteDouble _// WriteInt64 _' |
	sed 's_^// ReadDouble _// ReadInt64 _' \
	> nc_int64.go

cat nc_double.go |
	gofmt -r 'float64 -> uint32' |
	gofmt -r 'C.double -> C.uint' |
	gofmt -r 'NC_DOUBLE -> NC_UINT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uint' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uint' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uint' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uint' |
	gofmt -r 'WriteDouble -> WriteUint' |
	gofmt -r 'ReadDouble -> ReadUint' |
	sed 's_^// WriteDouble _// WriteUint _' |
	sed 's_^// ReadDouble _// ReadUint _' \
	> nc_uint.go

cat nc_double.go |
	gofmt -r 'float64 -> int32' |
	gofmt -r 'C.double -> C.int' |
	gofmt -r 'NC_DOUBLE -> NC_INT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_int' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_int' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_int' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_int' |
	gofmt -r 'WriteDouble -> WriteInt' |
	gofmt -r 'ReadDouble -> ReadInt' |
	sed 's_^// WriteDouble _// WriteInt _' |
	sed 's_^// ReadDouble _// ReadInt _' \
	> nc_int.go

cat nc_double.go |
	gofmt -r 'float64 -> uint16' |
	gofmt -r 'C.double -> C.ushort' |
	gofmt -r 'NC_DOUBLE -> NC_USHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ushort' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ushort' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ushort' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ushort' |
	gofmt -r 'WriteDouble -> WriteUshort' |
	gofmt -r 'ReadDouble -> ReadUshort' |
	sed 's_^// WriteDouble _// WriteUshort _' |
	sed 's_^// ReadDouble _// ReadUshort _' \
	> nc_ushort.go

cat nc_double.go |
	gofmt -r 'float64 -> int16' |
	gofmt -r 'C.double -> C.short' |
	gofmt -r 'NC_DOUBLE -> NC_SHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_short' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_short' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_short' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_short' |
	gofmt -r 'WriteDouble -> WriteShort' |
	gofmt -r 'ReadDouble -> ReadShort' |
	sed 's_^// WriteDouble _// WriteShort _' |
	sed 's_^// ReadDouble _// ReadShort _' \
	> nc_short.go

cat nc_double.go |
	gofmt -r 'float64 -> uint8' |
	gofmt -r 'C.double -> C.uchar' |
	gofmt -r 'NC_DOUBLE -> NC_UBYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uchar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uchar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uchar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uchar' |
	gofmt -r 'WriteDouble -> WriteUbyte' |
	gofmt -r 'ReadDouble -> ReadUbyte' |
	sed 's_^// WriteDouble _// WriteUbyte _' |
	sed 's_^// ReadDouble _// ReadUbyte _' \
	> nc_ubyte.go

cat nc_double.go |
	gofmt -r 'float64 -> int8' |
	gofmt -r 'C.double -> C.schar' |
	gofmt -r 'NC_DOUBLE -> NC_BYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_schar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_schar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_schar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_schar' |
	gofmt -r 'WriteDouble -> WriteByte' |
	gofmt -r 'ReadDouble -> ReadByte' |
	sed 's_^// WriteDouble _// WriteByte _' |
	sed 's_^// ReadDouble _// ReadByte _' \
	> nc_byte.go

cat nc_double.go |
	gofmt -r 'float64 -> byte' |
	gofmt -r 'C.double -> C.char' |
	gofmt -r 'NC_DOUBLE -> NC_CHAR' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_text' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_text' |
	gofmt -r 'C.nc_put_att_double(a, b, c, d, e, f) -> C.nc_put_att_text(a, b, c, e, f)' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_text' |
	gofmt -r 'WriteDouble -> WriteChar' |
	gofmt -r 'ReadDouble -> ReadChar' |
	sed 's_^// WriteDouble _// WriteChar _' |
	sed 's_^// ReadDouble _// ReadChar _' \
	> nc_char.go
