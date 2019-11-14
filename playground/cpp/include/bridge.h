#ifdef __cplusplus
extern "C" {
#endif
// expose a C API implemented in C++ so that go can use it
// g++ -fPIC -Wall -shared -std=c++17 -O3 -o libbridge.so ./src/bridge.cpp
    int Hello(int world);
    int SocketClient();
#ifdef __cplusplus
}
#endif