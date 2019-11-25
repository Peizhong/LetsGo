#ifdef __cplusplus
extern "C" {
#endif
// expose a C API implemented in C++ so that go can use it
// g++ -fPIC -Wall -shared -std=c++17 -O3 -o libbridge.so ./src/*
    typedef void* WrapClient;
    WrapClient WrapClientInit(void);
    void WrapClientFree(void*);
    bool Connect(void*,const char*,int);
    char*  WrapClientSendMessage(void*,char*);
    
    int Hello(int world);
    void PrintMessage(char* message);
#ifdef __cplusplus
}
#endif