package netcdfmem

// #cgo pkg-config: netcdf
// #include <stdlib.h>
// #include <netcdf.h>
// #include <netcdf_mem.h>
import "C"
import (
	"unsafe"

	"github.com/fhs/go-netcdf/netcdf"
)

// Flags define how the netCDF library should manage memory.
type Flags C.int

// Flags for Open.
const (
	MEMIO_LOCKED Flags = C.NC_MEMIO_LOCKED
)

// Dataset wraps netcdf.Dataset, adding methods specific to in-memory datasets.
//
// Must be closed via the added Close or CloseMem method to properly release memory.
type Dataset struct {
	netcdf.Dataset
}

// Open opens an existing netCDF dataset from a copy of the provided data.
// Path sets the dataset name.
// Mode is a bitwise-or of netcdf.FileMode values.
// Flags is a bitwise-or of Flags values.
func Open(path string, mode netcdf.FileMode, flags Flags, data []byte) (ds Dataset, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	memio := C.NC_memio{
		size:   C.size_t(len(data)),
		memory: C.CBytes(data),
		flags:  C.int(flags),
	}

	var id C.int
	err = newError(C.nc_open_memio(cpath, C.int(mode), &memio, &id))
	if err != nil {
		C.free(memio.memory)
	}

	ds.Dataset = netcdf.Dataset(id)
	return
}

// Create creates a new in-memory dataset.
// Path sets the dataset name.
// Mode is a bitwise-or of netcdf.FileMode values.
// InitialSize is a hint for the initial amount of memory to allocate for the dataset.
func Create(path string, mode netcdf.FileMode, initialSize int) (ds Dataset, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var id C.int
	err = newError(C.nc_create_mem(cpath, C.int(mode), C.size_t(initialSize), &id))
	ds.Dataset = netcdf.Dataset(id)

	return
}

/*
nc_close can be used to close and in-memory data set, but using nc_close_memio simplifies memory
management. With nc_close_memio the returned data is always owned by the caller and should be
freed when no longer needed. Using nc_close would require keeping track of the data pointer and
whether NC_MEMIO_LOCKED was specified, then conditionally freeing the memory if it was.


> * If the NC_MEMIO_LOCKED flag is set, then the netcdf library will make no attempt to reallocate
>   or free the provided memory. If the caller invokes the nc_close_memio() function to retrieve
>   the final memory block, it should be the same memory block as was provided when nc_open_memio
>   was called. Note that it is still possible to modify the in-memory file if the NC_WRITE mode
>   flag was set. However, failures can occur if an operation cannot complete because the memory
>   needs to be expanded.
> * If the NC_MEMIO_LOCKED flag is not set, then the netcdf library will take control of the
>   incoming memory. This means that the user should not make any attempt to free or even read
>   the incoming memory block in this case. The newcdf library is free to reallocate the incoming
>   memory block to obtain a larger block when an attempt to modify the in-memory file requires
>   more space. Note that implicit in this is that the old block – the one originally provided –
>   may be free'd as a side effect of re-allocating the memory using the realloc() function. The
>   caller may invoke the nc_close_memio() function to retrieve the final memory block, which may
>   not be the same as the originally block provided by the caller. In any case, the returned
>   block must always be freed by the caller and the original block should not be freed.
https://www.unidata.ucar.edu/software/netcdf/docs/md__Volumes_Workspace_releases_netcdf-c-4_87_84_netcdf-c_docs_inmemory.html
*/

// Close closes and releases the memory of the dataset.
//
// Use CloseMem to retrieve the in-memory data.
func (ds Dataset) Close() (err error) {
	var memio C.NC_memio
	err = newError(C.nc_close_memio(C.int(ds.Dataset), &memio))
	if memio.memory != nil {
		C.free(memio.memory)
	}
	return
}

// CloseMem closes the dataset and returns a copy of the in-memory data.
func (ds Dataset) CloseMem() (data []byte, err error) {
	var memio C.NC_memio
	err = newError(C.nc_close_memio(C.int(ds.Dataset), &memio))
	if memio.memory != nil {
		data = C.GoBytes(memio.memory, C.int(memio.size))
		C.free(memio.memory)
	}
	return
}

func newError(n C.int) error {
	if n == C.NC_NOERR {
		return nil
	}
	return netcdf.Error(n)
}
