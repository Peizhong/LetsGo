package main

/*
#cgo linux LDFLAGS: -lrt
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <string.h>
#include <sys/mman.h>
#define FILE_MODE (S_IRUSR | S_IWUSR | S_IRGRP | S_IROTH)
int my_shm_new(char *name) {
    shm_unlink(name);
    return shm_open(name, O_RDWR|O_CREAT|O_EXCL, FILE_MODE);
}
int my_shm_open(char *name) {
    return shm_open(name, O_RDWR,0666);
}
char *str = "const string";

typedef struct aa {
	int a,b;
	char p[128];
} aa;

aa a;

void show(){
	strcpy(a.p,"hello, i will show you something\0");
}

*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const SHM_NAME = "my_shm"
const SHM_SIZE = 4096

type MyData struct {
	Col1 int64
	Col2 int64
	Col3 int64
	v    [512]C.char
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
	C.munmap(ptr, SHM_SIZE)
}

func read() {
	fd, err := os.OpenFile("/dev/shm/my_shm", os.O_RDWR, os.ModeAppend)
	if err != nil {
		log.Println("open", err.Error())
		return
	}
	log.Println("fd: ", fd.Fd())
	mmap, err := syscall.Mmap(int(fd.Fd()), 0, SHM_SIZE, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("mmap", err.Error())
		os.Exit(1)
	}
	shmid, _, err := syscall.Syscall(syscall.SYS_SHMGET, 2, 4, 0600)
	shmaddr, _, err := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	pdd := (*MyData)(unsafe.Pointer(uintptr(shmaddr)))
	log.Println(pdd.v)
	C.show()
	data := (*MyData)(unsafe.Pointer(&mmap[0]))
	log.Println("data.col1", data.Col1)
	log.Println("str", C.GoString(&(data.v[0])))
	err = syscall.Munmap(mmap)
	fd.Close()
}

func main() {
	//write()
	read()
}
