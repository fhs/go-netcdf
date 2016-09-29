package netcdf_test

import (
	"fmt"
	"log"

	"github.com/fhs/go-netcdf/netcdf"
)

// CreateExampleFile creates an example NetCDF file containing only one variable.
func CreateExampleFile(filename string) error {
	// Create a new NetCDF 4 file. The dataset is returned.
	ds, err := netcdf.CreateFile("gopher.nc", netcdf.CLOBBER|netcdf.NETCDF4)
	if err != nil {
		return err
	}
	defer ds.Close()

	// Add the dimensions for our data to the dataset
	nhours := uint64(96)
	dimension, err := ds.AddDim("hour", nhours)
	if err != nil {
		return err
	}

	// Add the variable to the dataset that will store our data
	v, err := ds.AddVar("time", netcdf.DOUBLE, []netcdf.Dim{dimension})
	if err != nil {
		return err
	}

	attr, err := v.AddAttr("units")
	if err != nil {
		return err
	}
	err = attr.WriteBytes([]byte("hours since 2016-09-10"))
	if err != nil {
		return err
	}

	attr, err = v.AddAttr("standard_name")
	if err != nil {
		return err
	}
	err = attr.WriteBytes([]byte("time"))
	if err != nil {
		return err
	}

	hours := make([]float64, nhours)
	for h := uint64(0); h < nhours; h++ {
		hours[h] = float64(h * 3)
	}
	err = v.WriteFloat64s(hours)
	if err != nil {
		return err
	}

	// Add the dimensions for our data to the dataset
	dims := make([]netcdf.Dim, 2)
	ht, wd := 5, 4
	dims[0], err = ds.AddDim("height", uint64(ht))
	if err != nil {
		return err
	}
	dims[1], err = ds.AddDim("width", uint64(wd))
	if err != nil {
		return err
	}

	// Add the variable to the dataset that will store our data
	v, err = ds.AddVar("gopher", netcdf.UBYTE, dims)
	if err != nil {
		return err
	}

	// Create the data with the above dimensions and write it to the file.
	gopher := make([]uint8, ht*wd)
	i := 0
	for y := 0; y < ht; y++ {
		for x := 0; x < wd; x++ {
			gopher[i] = uint8(x + y)
			i++
		}
	}
	return v.WriteUint8s(gopher)
}

// ReadExampleFile reads the data in NetCDF file at filename and prints it out.
func ReadExampleFile(filename string) error {
	// Open example file in read-only mode. The dataset is returned.
	ds, err := netcdf.OpenFile(filename, netcdf.NOWRITE)
	if err != nil {
		return err
	}
	defer ds.Close()

	// Get the variable containing our data and read the data from the variable.
	v, err := ds.Var("gopher")
	if err != nil {
		return err
	}
	gopher, err := netcdf.GetUint8s(v)
	if err != nil {
		return err
	}

	// Get the length of the dimensions of the data.
	dims, err := v.LenDims()
	if err != nil {
		return err
	}

	// Print out the data
	i := 0
	for y := 0; y < int(dims[0]); y++ {
		for x := 0; x < int(dims[1]); x++ {
			fmt.Printf(" %d", gopher[i])
			i++
		}
		fmt.Printf("\n")
	}
	return nil
}

func DiscoverExampleFile(filename string) error {

	// Open example file in read-only mode. The dataset is returned.
	ds, err := netcdf.OpenFile(filename, netcdf.NOWRITE)
	if err != nil {
		return err
	}
	defer ds.Close()

	nvars, err := ds.NVars()
	if err != nil {
		return err
	}

	for i := 0; i < nvars; i++ {

		ncVar := ds.VarN(i)

		nvals, err := ncVar.Len()
		if err != nil {
			return err
		}

		nattrs, err := ncVar.NAttrs()
		if err != nil {
			return err
		}

		varname, err := ncVar.Name()
		if err != nil {
			return err
		}

		vartype, err := ncVar.Type()
		if err != nil {
			return err
		}

		fmt.Printf("Var %s [%s]\n", varname, vartype)

		// discover available attributes
		fmt.Printf("-- #attributes:%d\n", nattrs)
		for i := 0; i < nattrs; i++ {
			attr, err := ncVar.AttrN(i)
			if err != nil {
				return err
			}

			attrtype, err := attr.Type()
			if err != nil {
				return err
			}

			attrvalue, err := attr.GetBytes()
			if err != nil {
				return err
			}

			fmt.Printf("-- |attribute %d: %s Type:%s Value:%s\n", i, attr.Name(), attrtype, attrvalue)
		}

		// discover available dimensions
		dims, err := ncVar.Dims()
		if err != nil {
			return err
		}

		fmt.Printf("-- #dimensions: %d\n", len(dims))
		for i, dim := range dims {
			name, err := dim.Name()
			if err != nil {
				return err
			}

			dimlen, err := dim.Len()
			if err != nil {
				return err
			}

			fmt.Printf("-- |dimension %d: %s (len %d)\n", i, name, dimlen)
		}

		// discover available values
		fmt.Printf("-- #values:%d\n", nvals)

	}

	return nil

}

func Example() {
	// Create example file
	filename := "gopher.nc"
	if err := CreateExampleFile(filename); err != nil {
		log.Fatalf("creating example file failed: %v\n", err)
	}

	// Open and read example file
	if err := ReadExampleFile(filename); err != nil {
		log.Fatalf("reading example file failed: %v\n", err)
	}

	// Example files for netCDF: http://www.unidata.ucar.edu/software/netcdf/examples/files.html
	if err := DiscoverExampleFile(filename); err != nil {
		log.Fatalf("discovering example file %s failed: %v\n", filename, err)
	}

	// Output:
	//  0 1 2 3
	//  1 2 3 4
	//  2 3 4 5
	//  3 4 5 6
	//  4 5 6 7
	// Var time [DOUBLE]
	// -- #attributes:2
	// -- |attribute 0: units Type:CHAR Value:hours since 2016-09-10
	// -- |attribute 1: standard_name Type:CHAR Value:time
	// -- #dimensions: 1
	// -- |dimension 0: hour (len 96)
	// -- #values:96
	// Var gopher [UBYTE]
	// -- #attributes:0
	// -- #dimensions: 2
	// -- |dimension 0: height (len 5)
	// -- |dimension 1: width (len 4)
	// -- #values:20

}
