#!/bin/sh

cat nc_double.go |
	gofmt -r 'float64 -> float32' |
	gofmt -r 'C.double -> C.float' |
	gofmt -r 'NC_DOUBLE -> NC_FLOAT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_float' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_float' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_float' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_float' |
	gofmt -r 'PutDouble -> PutFloat' |
	gofmt -r 'GetDouble -> GetFloat' |
	sed 's_^// PutDouble _// PutFloat _' |
	sed 's_^// GetDouble _// GetFloat _' \
	> nc_float.go

cat nc_double.go |
	gofmt -r 'float64 -> uint64' |
	gofmt -r 'C.double -> C.ulonglong' |
	gofmt -r 'NC_DOUBLE -> NC_UINT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ulonglong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ulonglong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ulonglong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ulonglong' |
	gofmt -r 'PutDouble -> PutUint64' |
	gofmt -r 'GetDouble -> GetUint64' |
	sed 's_^// PutDouble _// PutUint64 _' |
	sed 's_^// GetDouble _// GetUint64 _' \
	> nc_uint64.go

cat nc_double.go |
	gofmt -r 'float64 -> int64' |
	gofmt -r 'C.double -> C.longlong' |
	gofmt -r 'NC_DOUBLE -> NC_INT64' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_longlong' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_longlong' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_longlong' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_longlong' |
	gofmt -r 'PutDouble -> PutInt64' |
	gofmt -r 'GetDouble -> GetInt64' |
	sed 's_^// PutDouble _// PutInt64 _' |
	sed 's_^// GetDouble _// GetInt64 _' \
	> nc_int64.go

cat nc_double.go |
	gofmt -r 'float64 -> uint32' |
	gofmt -r 'C.double -> C.uint' |
	gofmt -r 'NC_DOUBLE -> NC_UINT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uint' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uint' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uint' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uint' |
	gofmt -r 'PutDouble -> PutUint' |
	gofmt -r 'GetDouble -> GetUint' |
	sed 's_^// PutDouble _// PutUint _' |
	sed 's_^// GetDouble _// GetUint _' \
	> nc_uint.go

cat nc_double.go |
	gofmt -r 'float64 -> int32' |
	gofmt -r 'C.double -> C.int' |
	gofmt -r 'NC_DOUBLE -> NC_INT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_int' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_int' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_int' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_int' |
	gofmt -r 'PutDouble -> PutInt' |
	gofmt -r 'GetDouble -> GetInt' |
	sed 's_^// PutDouble _// PutInt _' |
	sed 's_^// GetDouble _// GetInt _' \
	> nc_int.go

cat nc_double.go |
	gofmt -r 'float64 -> uint16' |
	gofmt -r 'C.double -> C.ushort' |
	gofmt -r 'NC_DOUBLE -> NC_USHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_ushort' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_ushort' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_ushort' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_ushort' |
	gofmt -r 'PutDouble -> PutUshort' |
	gofmt -r 'GetDouble -> GetUshort' |
	sed 's_^// PutDouble _// PutUshort _' |
	sed 's_^// GetDouble _// GetUshort _' \
	> nc_ushort.go

cat nc_double.go |
	gofmt -r 'float64 -> int16' |
	gofmt -r 'C.double -> C.short' |
	gofmt -r 'NC_DOUBLE -> NC_SHORT' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_short' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_short' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_short' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_short' |
	gofmt -r 'PutDouble -> PutShort' |
	gofmt -r 'GetDouble -> GetShort' |
	sed 's_^// PutDouble _// PutShort _' |
	sed 's_^// GetDouble _// GetShort _' \
	> nc_short.go

cat nc_double.go |
	gofmt -r 'float64 -> uint8' |
	gofmt -r 'C.double -> C.uchar' |
	gofmt -r 'NC_DOUBLE -> NC_UBYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_uchar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_uchar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_uchar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_uchar' |
	gofmt -r 'PutDouble -> PutUbyte' |
	gofmt -r 'GetDouble -> GetUbyte' |
	sed 's_^// PutDouble _// PutUbyte _' |
	sed 's_^// GetDouble _// GetUbyte _' \
	> nc_ubyte.go

cat nc_double.go |
	gofmt -r 'float64 -> int8' |
	gofmt -r 'C.double -> C.schar' |
	gofmt -r 'NC_DOUBLE -> NC_BYTE' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_schar' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_schar' |
	gofmt -r 'C.nc_put_att_double -> C.nc_put_att_schar' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_schar' |
	gofmt -r 'PutDouble -> PutByte' |
	gofmt -r 'GetDouble -> GetByte' |
	sed 's_^// PutDouble _// PutByte _' |
	sed 's_^// GetDouble _// GetByte _' \
	> nc_byte.go

cat nc_double.go |
	gofmt -r 'float64 -> byte' |
	gofmt -r 'C.double -> C.char' |
	gofmt -r 'NC_DOUBLE -> NC_CHAR' |
	gofmt -r 'C.nc_put_var_double -> C.nc_put_var_text' |
	gofmt -r 'C.nc_get_var_double -> C.nc_get_var_text' |
	gofmt -r 'C.nc_put_att_double(a, b, c, d, e, f) -> C.nc_put_att_text(a, b, c, e, f)' |
	gofmt -r 'C.nc_get_att_double -> C.nc_get_att_text' |
	gofmt -r 'PutDouble -> PutChar' |
	gofmt -r 'GetDouble -> GetChar' |
	sed 's_^// PutDouble _// PutChar _' |
	sed 's_^// GetDouble _// GetChar _' \
	> nc_char.go
