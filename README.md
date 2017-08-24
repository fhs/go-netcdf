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

First, make sure you have the netCDF C library is installed. Most Linux distributions have a package for it: `libnetcdf-dev` in Ubuntu/Debian, `netcdf` in ArchLinux, etc. You can also download the source from [Unidata](https://www.unidata.ucar.edu/downloads/netcdf/index.jsp), compile and install it.

Then, to install go-netcdf, run:

	$ go get github.com/fhs/go-netcdf/netcdf
