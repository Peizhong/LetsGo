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
        int **i2arr;
        make2dArray<int>(i2arr,3,4);
        delete2dArray<int>(i2arr,3);

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
    int arr[] = {6,7,8,9};
    permutations<int>(arr,0,2);

    int a=1;
    int& h2 = f(a);
    // 模板函数
    int b = 2;
    int c = 3;
    int v= abc<int>(1,b,c);

    c = Count<int>(arr,sizeof(arr));
    
    currency cc = currency(sign_minus,100,1);
    currency c2 = cc.Add(cc);

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

currency::currency(signType sign,unsigned long dollar,unsigned int cent)
{
    Set(sign,dollar,cent);
}

void currency::Set(signType sign,unsigned long dollar,unsigned int cent)
{
    _sign = sign;
    _dollar = dollar;
    _cent = cent;
}

currency currency::Add(const currency& x) const
{
    currency result;
    unsigned long d = _dollar + x._dollar;
    unsigned int c = _cent +x._cent;
    result.Set(sign_plus,d,c);
    return result;
}