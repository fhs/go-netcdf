package ncmem

// #cgo pkg-config: netcdf
// #include <stdlib.h>
// #include <netcdf.h>
// #include <netcdf_mem.h>
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"reflect"
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
// Must be closed via one of the added Close, CloseBytes, or CloseCopyBytes
// methods to properly release memory.
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

// OpenReader reads and opens an existing netCDF dataset from r.
// Path sets the dataset name.
// Mode is a bitwise-or of netcdf.FileMode values.
// Flags is a bitwise-or of Flags values.
//
// If r fulfills LenReader, it's Len method will be used to determine how much
// memory to allocate.
func OpenReader(path string, mode netcdf.FileMode, flags Flags, r io.Reader) (ds Dataset, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	size, ptr, err := readAll(r)
	if err != nil {
		return ds, err
	}

	memio := C.NC_memio{
		size:   C.size_t(size),
		memory: ptr,
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

// CloseCopyBytes closes the dataset and returns a copy of the in-memory data.
func (ds Dataset) CloseCopyBytes() (data []byte, err error) {
	var memio C.NC_memio
	err = newError(C.nc_close_memio(C.int(ds.Dataset), &memio))
	if memio.memory != nil {
		data = C.GoBytes(memio.memory, C.int(memio.size))
		C.free(memio.memory)
	}
	return
}

// CloseBytes closes the dataset and returns a reference to C allocated
// memory.
func (ds Dataset) CloseBytes() (*Bytes, error) {
	var memio C.NC_memio
	err := newError(C.nc_close_memio(C.int(ds.Dataset), &memio))
	if err != nil {
		return nil, err
	}

	data, err := unsafeByteSlice(memio.memory, uint64(memio.size))
	if err != nil {
		C.free(memio.memory)
		return nil, err
	}

	return &Bytes{
		Data: data,
		ptr:  memio.memory,
	}, nil
}

// Bytes provides a view into C allocated data. It's Free method
// method must be called to release the memory.
type Bytes struct {
	Data []byte
	ptr  unsafe.Pointer
}

// Free releases the underlying data. All references to Data are
// invalid after calling.
func (b *Bytes) Free() {
	if b.ptr != nil {
		C.free(b.ptr)
		b.Data = nil
		b.ptr = nil
	}
}

// LenReader is an io.Reader which can provide it's unread length.
type LenReader interface {
	io.Reader
	Len() int
}

func readAll(r io.Reader) (C.size_t, unsafe.Pointer, error) {
	if lr, ok := r.(LenReader); ok {
		// allocate the full amount of memory needed
		length := lr.Len()
		ptr := C.malloc(C.size_t(length))

		// create a slice view
		data, err := unsafeByteSlice(ptr, uint64(length))
		if err != nil {
			C.free(ptr)
			return 0, nil, err
		}

		// read out the data
		_, err = io.ReadFull(r, data)
		if err != nil {
			C.free(ptr)
			return 0, nil, err
		}
		return C.size_t(length), ptr, nil
	}

	// The amount of memory required is unknown. Use a growth strategy
	// similar to ioutil.ReadAll. Start with bytes.MinRead and double
	// when filled.
	var (
		ptr    = C.malloc(bytes.MinRead)
		length = C.size_t(0)
		cap    = C.size_t(bytes.MinRead)
	)

	for {
		// if cap has been reached, realloc at twice the size.
		if length >= cap {
			cap *= 2
			newPtr, err := C.realloc(ptr, cap)
			if err != nil {
				C.free(ptr) // if realloc fails, existing ptr must be freed
				return 0, nil, fmt.Errorf("failed realloc to %d: %v", cap, err)
			}
			ptr = newPtr
		}

		// create slice view
		s, err := unsafeByteSlice(ptr, uint64(cap))
		if err != nil {
			C.free(ptr)
			return 0, nil, err
		}

		// read additional data
		n, err := r.Read(s[length:])
		length += C.size_t(n)
		if err == io.EOF {
			// io.EOF indicates r has been read to completion
			return length, ptr, nil
		}
		if err != nil {
			C.free(ptr) // reading failed, free the data
			return 0, nil, err
		}
	}
}

// unsafeByteSlice creates a []byte view into ptr.
func unsafeByteSlice(ptr unsafe.Pointer, size uint64) ([]byte, error) {
	const maxInt = uint64(^uint(0) >> 1)
	if size > maxInt {
		return nil, fmt.Errorf("unsafeByteSlice: size %d larger than max %d", size, maxInt)
	}

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr),
		Len:  int(size),
		Cap:  int(size),
	})), nil
}

func newError(n C.int) error {
	if n == C.NC_NOERR {
		return nil
	}
	return netcdf.Error(n)
}
