[![Build Status](https://travis-ci.org/fhs/go-netcdf.png)](https://travis-ci.org/fhs/go-netcdf)
[![GoDoc](https://godoc.org/github.com/fhs/go-netcdf/netcdf?status.svg)](https://godoc.org/github.com/fhs/go-netcdf/netcdf)

## Overview

Package netcdf is a Go binding for the netCDF C library.
This package supports netCDF version 3, and 4 if
netCDF 4 support is enabled in the C library.

## Documentation

Documentation can be found here:
http://godoc.org/github.com/fhs/go-netcdf/netcdf

## Installation

How to install:

	$ go get github.com/fhs/go-netcdf/netcdf

## Example

### Open an existing netcdf file

	nc, err := netcdf.CreateFile("gopher.nc", netcdf.NOWRITE)
	if err != nil {
		return err
	}
	defer nc.Close()

### Variables

    var, err := nc.Var("time")
    if err == nil {
    	nattrs,_ := var.NAttrs()
    	for i:=0;i<nattrs;i++ {
    		attr, err1 := var.AttrN(i)
    		if err1 == nill {
    			aname := attr.Name()
    			atype, _ := attr.Type()
	    		fmt.Printf("Variable `time` has attribute %s of type %d [%s]\n",name,atype,atype.Name())
	    	}
 	
    	}
    }

### Attributes

    val, err := var.Attr("units").ValueString()
    if err == nil {
        fmt.Printf("The value of the attibute `units` is `%s`\n",val)
    }



