#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <assert.h>
#include <string.h>
#include <stdint.h>
#include <errno.h>
#include <math.h>

#define shm_name "/my_shm"
// page size 4kb 倍数
#define shm_size 4096

#define BUFFER_B 8

// https://blog.csdn.net/bbzhaohui/article/details/81665370
void setShareMem()
{
    // shm_open创建posix共享内存对象，/dev/shm的tmpfs文件系统是基于内存的
    int shmfd = shm_open(shm_name,O_CREAT|O_RDWR,0666);
    assert(shmfd!=-1);
    // 设置共享大小
    int ret = ftruncate(shmfd,shm_size);
    assert(ret!=-1);
    // mmap也可以映射普通文件
    void *map = mmap(NULL,shm_size,PROT_READ|PROT_WRITE,MAP_SHARED,shmfd,0);
    assert(map!=MAP_FAILED);
    sprintf((char*)map,"hello, i will show you something");  
    printf("%s\n",(char*)map);
    // memset(map,43,shm_size);
    close(shmfd);
    // 删除/dev/shm的文件
    shm_unlink(shm_name);
    munmap(map, shm_size);
}

typedef struct  {
    int64_t col1,col2,col3;
    char v;
} MyData,*pMyData;

void readShareMem()
{
    // shm_open创建posix共享内存对象
    int shmfd = shm_open(shm_name,O_RDONLY,0666);
    assert(shmfd!=-1);
    void *map = mmap(NULL,shm_size,PROT_READ,MAP_SHARED,shmfd,0);
    assert(map!=MAP_FAILED);
    pMyData dd = (pMyData)map;
    printf("col1:%ld,col2:%ld,col3:%ld,v:%c\n",dd->col1,dd->col2,dd->col3,dd->v);
    int64_t *col = (int64_t*)map;
    for (int i=0;i<3;i++)
    {
        printf("%ld\n",*(col++));
    }
    close(shmfd);
    munmap(map, shm_size);
}

void enqueue(){
    int num = 10;
    int res = __sync_bool_compare_and_swap(&num,12,11);
    printf("cas:%d\n",res);
}

pMyData bufffer = NULL;

// gcc go.c -lrt
void main()
{
    int buffersize = pow(2,BUFFER_B);
    printf("%d\n",buffersize);
    bufffer = (pMyData)malloc(sizeof(MyData)*buffersize);
    for (int i=0;i<buffersize;i++)
    {
        bufffer[i].col1 = i;
    }
    for (int i=0;i<buffersize;i++)
    {
        printf("%d",bufffer[i].col1);
    }
    enqueue();
    free(bufffer);
    // readShareMem();
}