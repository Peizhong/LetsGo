#include "../include/server.h"
#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <unistd.h>
#include <strings.h>
#include <fcntl.h> 
#include <sys/socket.h>
#include <arpa/inet.h>
#include <sys/epoll.h>

#define PORT 8080
#define BUFFER_SIZE 1024
#define MAX_EVENTS 10

int setnonblocking(int sock)
{
   int opts;
   opts = fcntl(sock, F_GETFL);
   if(opts < 0)
   {
       return -1;
   }
   opts = opts | O_NONBLOCK;
   if(fcntl(sock, F_SETFL, opts) < 0)
   {
       return -1;
   }
   return 0;
}

void readData(int fd)
{
    char buf[BUFFER_SIZE];
    size_t bytes_read = recv(fd,buf,BUFFER_SIZE,0);
    buf[bytes_read]='\0';
    std::cout<<buf<<std::endl;
}

int EpollServer()
{
    unsigned int evenx = EPOLLIN | EPOLLET;

    struct epoll_event ev, events[MAX_EVENTS];
    int listen_sock, conn_sock, nfds, epollfd;

    listen_sock = socket(PF_INET,SOCK_STREAM,0);
    if (listen_sock == -1)
    {
        return -1;
    }
    struct sockaddr_in addr;
    socklen_t addrlen;
    bzero(&addr, sizeof(addr));     
    addr.sin_family = PF_INET;
    addr.sin_port = (in_port_t)htons(PORT);
    addr.sin_addr.s_addr = htonl(INADDR_ANY);
    if (bind(listen_sock, (struct sockaddr *)&addr, sizeof(addr)) == -1)
    {
        return -1;
    }
    int res = listen(listen_sock,MAX_EVENTS);
    if (res==-1)
    {
        return -1;
    }
    // 创建epoll句柄
    epollfd = epoll_create1(0);
    if (epollfd == -1)
    {
        return -1;
    }
    // 用于 注册感兴趣的事件 和回传时发生的事件
    ev.data.fd = listen_sock;
    ev.events = EPOLLIN; // 文件可读
    if (epoll_ctl(epollfd, EPOLL_CTL_ADD, listen_sock, &ev) == -1)
    {
        return -1;
    }
    for (;;)
    {
        nfds = epoll_wait(epollfd, events, MAX_EVENTS, -1);
        if (nfds == -1)
        {
            return -1;
        }
        for (int n=0;n<nfds;n++)
        {
            // 监听节点
            if (events[n].data.fd == listen_sock)
            {
                conn_sock = accept(listen_sock,(struct sockaddr *) &addr, &addrlen);
                if (conn_sock == -1)
                {
                    break;
                }
                setnonblocking(conn_sock);
                ev.data.fd = conn_sock;
                // EPOLLET 边缘触发
                // EPOLLONESHOT 最多触发一次，用完要再epoll_ctl一次
                ev.events = EPOLLIN | EPOLLET | EPOLLONESHOT;
                // 添加连接进入epoll监听列表
                int res = epoll_ctl(epollfd, EPOLL_CTL_ADD, conn_sock,&ev);
                if (res==-1)
                {
                    break;
                }
            }
            else
            {
                readData(events[n].data.fd);
            }
        }
    }
    close(epollfd);
}