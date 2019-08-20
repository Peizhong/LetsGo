#include <iostream>
#include "../include/basic.h"

using namespace std;

int& f(int& a)
{
    int b = 123;
    a = a+b;
    // return b; // 没有意义，局部变量在函数结束后被回收
    return a;
}

// 函数签名：由函数的形参类型和形参个数确定的
int f1(int a,int b)
{
    return a+b;
}
// f1(int& a,int& b)不行
int f1_1(int& a,int& b)
{
    return a+b;
}

void doExpt(int a)
{
    if (a>100)
    {
        // char*
        throw "more than 100";
    }
}

void doNew(int n)
{
    try
    {
        int* x = new int[n];
        delete x;
        
        char (*c)[5];
        c = new char[n][5];
        delete c;
    }
    catch(bad_alloc e)
    {
        std::cerr << e.what() << '\n';
    }
}

void Hello()
{
    int a=1;
    int& h2 = f(a);
    // 模板函数
    int b = 2;
    int c = 3;
    int v= abc<int>(1,b,c);

    int arr[] = {1,2,3,4,5};
    c = Count<int>(arr,sizeof(arr));
    
    doNew(100);
    
    try
    {
        doExpt(101);
    }
    catch (const char* e)
    {
        cout<<e<<endl;
    }
}