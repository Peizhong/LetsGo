package main

/*
#cgo linux LDFLAGS: -lrt
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>
#define FILE_MODE (S_IRUSR | S_IWUSR | S_IRGRP | S_IROTH)
int my_shm_new(char *name) {
    shm_unlink(name);
    return shm_open(name, O_RDWR|O_CREAT|O_EXCL, FILE_MODE);
}
int my_shm_open(char *name) {
    return shm_open(name, O_RDWR,0666);
}
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

const SHM_NAME = "my_shm"
const SHM_SIZE = 4096

type MyData struct {
	Col1 int64
	Col2 int64
	Col3 int64
	Ch   C.char
}

func write() {
	fd, err := C.my_shm_new(C.CString(SHM_NAME))
	if err != nil {
		return
	}
	C.ftruncate(fd, SHM_SIZE)
	ptr, err := C.mmap(nil, SHM_SIZE, C.PROT_READ|C.PROT_WRITE, C.MAP_SHARED, fd, 0)
	if err != nil {
		return
	}
	C.close(fd)
	data := (*MyData)(unsafe.Pointer(ptr))
	data.Col1 = 43
	data.Col2 = 86
	data.Col3 = 108
	data.Ch = 'x'
	C.munmap(ptr, SHM_SIZE)
}

func read() {
	fd, err := C.my_shm_open(C.CString(SHM_NAME))
	if err != nil {
		return
	}
	log.Println("fd: ", fd)
	ptr, err := C.mmap(nil, SHM_SIZE, C.PROT_READ, C.MAP_SHARED, fd, 0)
	if err != nil {
		return
	}
	log.Println("Ptr: ", ptr)
	C.close(fd)
	data := (*MyData)(unsafe.Pointer(ptr))
	fmt.Println(data)
	C.munmap(ptr, SHM_SIZE)
}

func main() {
	write()
	// read()
}
