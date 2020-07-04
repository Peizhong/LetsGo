#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <pthread.h>
#include <signal.h>
#include <errno.h>
#include <semaphore.h>
#include <sys/wait.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/msg.h>
#include <sys/shm.h>

#include <iostream>
#include <thread>

#include "tree.h"

#define FIFO_PATH "ff"
#define SEM_PATH "sm"
#define MSG_KEY 43
#define SHM_KEY 23

#define max(a,b) (a>b?a:b)

using namespace std;


void doPipe()
{
    int fd[2];
    pipe(fd);

    pid_t  pid = fork();
    if (pid==-1){
        perror("fork error");
        exit(EXIT_FAILURE);
    }
    int bf = open("test",O_RDWR);
    if (bf<1)
    {
        perror("open file error");
        exit(EXIT_FAILURE);
    }
    if (pid==0)
    {
        puts("child wait");
        for (int v =0;v<10000;v++)
        {
            for (int w=0;w<100000;w++)
            {
                
            }
        }
        puts("child read bf");
        char buf[50];
        int n = read(bf,buf,50);
        if (n<1)
        {
            puts("child can't read");
        }
        printf("child's bf: %s",buf);
        // 子进程关闭写
        close(fd[1]);
        n = read(fd[0],buf,50);
        printf("read something: %s",buf);
        // write(STDOUT_FILENO,buf,n);
    }
    else
    {
        close(bf);
        puts("parent close bf");
        char buf[50];
        int n = read(bf,buf,50);
        if (n<1)
        {
            puts("parent can't read");
        }
        printf("%x\n",fd[1]);
        // 父进程关闭读
        close(fd[0]);
        write(fd[1], "hello pipe\n", 12);

        wait(NULL);
        puts("bye bye");
    }
}

void* why()
{
    return NULL;
}

void doFifo(){
    if (access(FIFO_PATH, F_OK) == -1)
	{
        int ffid = mkfifo(FIFO_PATH,S_IRUSR|S_IWUSR);
        if (ffid<0 && errno!=EEXIST)
        {
            perror("make fifo failed:");
            return;
        }
        // ls -l: ff prw-------
        puts("create fifo");
    }
}

// posix的代码，编译时增加-pthread，动态库，选择thread-safe的实现
// gcc -pthread ipc.cpp -o ipc
void doSem()
{
    sem_t sem;
    sem_init(&sem, 0, 0);
    sem_post(&sem);
    sem_wait(&sem);
    sem_destroy(&sem);
    
    sem_t *psem;
    psem = sem_open(SEM_PATH, O_CREAT | O_EXCL, 0666, 0);
    sem_wait(psem);
    sem_close(psem);
}

typedef struct {
    // 第一个字段必须为long，msgrcv根据此接收消息
    long type;
    // 消息正文
    char buf[512];

} Message;

void doMsg()
{
    // 如果没有该块共享内存，则创建，并返回共享内存ID。若已有该块共享内存，则返回-1
    int queueid = msgget(MSG_KEY, IPC_CREAT| IPC_EXCL | 0666);
    Message *msg = (Message*)malloc(sizeof(Message));
    memset(msg,0x00,sizeof(Message));
    long msgtype = 43;
    msg->type=msgtype;
    strcpy(msg->buf,"hello");
    // IPC_NOWAIT: 当消息队列已满的时候，msgsnd函数不等待立即返回
    int sendlength = sizeof(Message)-sizeof(msg->type);
    int res = msgsnd(queueid,msg,sendlength,IPC_NOWAIT);
    // 指定要接收的msg.type
    msgrcv(queueid,msg,sendlength,msgtype,0);
    free(msg);
}

// 进程不安全，借助信号量等同步
void doShareMem()
{
    // 创建共享内存
    int shmid= shmget(SHM_KEY, 4096, IPC_CREAT | IPC_EXCL | 0666);
    // 共享内存连接到当前进程
    void* shm = shmat(shmid, 0, 0);
    Message *msg = (Message*)shm;
    // 分离共享内存
    shmdt(shm);
}

void sigusr1_handler(int signal)
{
    printf("signal %d received\n", signal);
}

void doSignal()
{
    int pid = getpid();
    // handler for the signal SIG to HANDLER
    signal(SIGUSR1,sigusr1_handler);
    // 发送信号
    kill(pid,SIGUSR1);
}

typedef class fClass {
    public:
    int a,b;
    private:
    int c,d;
    // 友元函数
    friend fClass duplicate(fClass);
    // 友元类
    friend class fClass2;
} fClass;

typedef class fClass2{
    public:
        fClass2(int v);
        void Dofriend();
        virtual void vv(); // 虚拟
        virtual void vv0()=0; // 抽象
    protected:
        int f2;
        void DoproFriend();
    private:
        void F2Private();
} fClass2;

void fClass2::Dofriend(){
    fClass f;
    int a = f.c+f.d;
}

typedef class fClass3 : public fClass2{
    public:
        int a;
        void DoClass3();
        fClass3();
        void vv0();
} fClass3;

typedef class fClass4 : protected fClass2{
    public:
        int a;
        void DoClass4();
        fClass4();
        void vv0();
} fClass4;


void fClass4::vv0(){
    this->DoproFriend();
}

typedef class fClass5 : private fClass2{
    public:
        int a;
        void DoClass4();
        fClass5();
        void vv0();
} fClass5;


void fClass5::vv0(){

}

fClass3::fClass3():fClass2(3){
    a = this->f2;
}

void fClass3::DoClass3(){
    this->DoproFriend();
    printf("%s:%d@%s\n",__func__,__LINE__,__FILE__);
}

void fClass3::vv0(){

}

void trun(int n){
    auto id = std::this_thread::get_id();
    std::this_thread::yield();
}

void doThread(){
    auto t = thread(trun,2);
    t.join();

    // 通过构造函数
    thread t2(trun,2);
    t2.join();

    thread t3(trun,3);
    t3.detach(); // main线程退出，进程不会退出
}

int main(){
    int a= 0;
    auto str = "ada\n";
    doMsg();
    a ++;
    fClass f;
    fClass3 f3;
    f3.Dofriend();
    duplicate(f);
    exit(EXIT_SUCCESS);
    fClass4 f4;
    fClass5 f5;
    f4.Dofriend();
}