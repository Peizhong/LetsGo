#ifndef _CLIENT_H
#define _CLIENT_H

#define BUFFER_SIZE 1024
class Client
{
    private:
        char addr[16];
        int port;
        int fd;
        char sendbuf[BUFFER_SIZE];
        char recvbuf[BUFFER_SIZE];
    public:
        Client();
        
        bool Connect(const char* addr, int port);
        char* SendMessage(char* message);
};

int client();

#endif