#include <stdio.h>
#include "../include/bridge.h"
#include "../include/client.h"

int Hello(int world)
{
    return world*world;
}

void PrintMessage(char* message) {
    printf("Go send me %s\n", message);
}

WrapClient WrapClientInit(void)
{
    Client *p = new Client();
    return (void*)p;
}

bool Connect(void* c,const char* addr,int port)
{
     Client *p = (Client *)c;
     bool res = p->Connect(addr,port);
     return res;
}

void WrapClientFree(void* c)
{
     Client *p = (Client *)c;
     delete p;
}

char* WrapClientSendMessage(void* c,char* message)
{
     Client *p = (Client *)c;
     p->SendMessage(message);
     return message;
}