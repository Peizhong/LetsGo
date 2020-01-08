#include <iostream>
#include <vector>
#include <string>
#include <sstream>

#include "../include/bridge.h"

using namespace std;

// g++ -I./include/ ./src/*.cpp -g 
int main(){
    WrapClient client = WrapClientInit();
    Connect(client,"127.0.0.1",8081);
    return 0;
}