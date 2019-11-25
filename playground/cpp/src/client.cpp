#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h> 
#include <arpa/inet.h>
#include <netdb.h>
#include "../include/client.h"

int client()
{
    int sockfd;
    struct sockaddr_in serv_addr;
    // socket: 描述通讯方式; socket_stream: 面向连接,流式
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd<0)
    {
        return -1;
    } 
    bzero(&serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = inet_addr("127.0.0.1");
    serv_addr.sin_port = htons(8081);
    if (connect(sockfd, (struct sockaddr *)&serv_addr, sizeof(serv_addr))<1)
    {
        return -1;
    }
    char sendbuf[BUFFER_SIZE];
    char recvbuf[BUFFER_SIZE];
    memset(sendbuf, 0, sizeof(sendbuf));
    memset(recvbuf, 0, sizeof(recvbuf));
    while (1)
    {
        send(sockfd,sendbuf,strlen(sendbuf),0);
    }
    close(sockfd);
    return sockfd;
}

Client::Client()
{

}

bool Client::Connect(const char* addr, int port)
{
    this->port = port;
    // strncpy
    strcpy(this->addr,addr);
    this->fd = socket(AF_INET, SOCK_STREAM, 0);
    if (this->fd<0)
    {
        return false;
    }
    struct sockaddr_in serv_addr;
    bzero(&serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = inet_addr(this->addr);
    serv_addr.sin_port = htons(this->port);
    if (connect(this->fd, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0)
    {
        return false;
    }
    while (true)
    {
        // while(true){ send();recv(); }只能按顺序执行
        // close(fd);
        // 异步方式
        fd_set rfds;
        struct timeval tv;
        int retval, maxfd;
        // 把可读文件描述符的集合清空
        FD_ZERO(&rfds);
        // 把标准输入的文件描述符加入到集合中
        FD_SET(0,&rfds);
        // 把当前连接的文件描述符加入到集合中
        FD_SET(this->fd, &rfds);
        // 找出文件描述符集合中最大的文件描述符
        maxfd = this->fd;
        // 设置超时时间，等待select
        tv.tv_sec = 5;
        tv.tv_usec = 0;
        retval = select(maxfd+1, &rfds, NULL, NULL, &tv);
        if (retval==-1)
        {
            puts("select -1");
        }
        else if (retval==0)
        {
            puts("select 0");
            continue;
        }
        else
        {
            // socket数据
            if (FD_ISSET(this->fd,&rfds))
            {
                // memset(this->recvbuf, 0, BUFFER_SIZE);
                memset(this->recvbuf, 0, BUFFER_SIZE);
                int len = recv(this->fd, this->recvbuf, BUFFER_SIZE,0);
                puts(this->recvbuf);
            }
            // 用户输入了信息
            if (FD_ISSET(0,&rfds))
            {
                memset(this->sendbuf, 0, BUFFER_SIZE);
                fgets(this->sendbuf, BUFFER_SIZE, stdin);
                send(this->fd,this->sendbuf, strlen(this->sendbuf),0);
            }
        }
    }
    
    return true;
}

char* Client::SendMessage(char* message)
{
    char str[80];
    sprintf(str, "mock send = %s", message);
    return str;
}