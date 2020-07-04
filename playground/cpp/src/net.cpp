#include <stdio.h>  
#include <stdlib.h>
#include <string.h>  
#include <fcntl.h>
#include <event2/event.h>  
#include <event2/bufferevent.h>  

int eserver(){
    int server_socketfd; //服务端socket    
    struct sockaddr_in server_addr; //服务器网络地址结构体      
    memset(&server_addr,0,sizeof(server_addr)); //数据初始化--清零      
    server_addr.sin_family = AF_INET; //设置为IP通信      
    server_addr.sin_addr.s_addr = INADDR_ANY;//服务器IP地址--允许连接到所有本地地址上      
    server_addr.sin_port = htons(8001); //服务器端口号      
    
    //创建服务端套接字    
    server_socketfd = socket(PF_INET,SOCK_STREAM,0);
    if (server_socketfd < 0) {
        puts("socket error");
        return 0;
    }
    evutil_make_listen_socket_reuseable(server_socketfd); //设置端口重用

    struct event_base *base_ev = event_base_new();
    const char *x = event_base_get_method(base_ev); //获取IO多路复用的模型，linux一般为epoll  
    printf("METHOD:%s\n", x);

    return 0;
}