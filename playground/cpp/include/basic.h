#ifndef _BASIC_H
#define _BASIC_H

#include <iostream>
#include <string>

using namespace std;

typedef struct Element ElementType;

struct Element
{
    int ID;
    std::string Value;

    Element()
    {

    }

    Element(int id,std::string value)
    {
        ID = id;
        Value = value;
    }
    
    bool operator==(Element const& b) const
    {
        return this->ID==b.ID && this->Value==b.Value;
    }

    bool operator!=(Element const& b) const
    {
        return this->ID!=b.ID || this->Value!=b.Value;
    }

    bool operator<(Element const& b) const
    {
        return this->ID < b.ID;
    }

    bool operator>(Element const& b) const
    {
        return this->ID > b.ID;
    }
    
    /*
    bool operator==(Element const& b) const
    {
        return true;
    }
    bool operator<(Element const& b) const
    {
        return true;
    }
     */
};

// 模板函数
template<class T>
T abc(T a, T& b, const T& c)
{
    // 形参：引用参数、常量引用参数
    // 实参
    // 调用时，引用参数不用复制；返回时没有析构
    return a+b+c;
}

template<class T>
int Count(T arr[],int n)
{
    // int total = sizeof(arr); // 传递的是数组第一个元素
    int total = n;
    int item = sizeof(T);
    return total/item;
}

template<class T>
bool make2dArray(T ** &x, int row, int col)
{
    try
    {
        // 行指针
        x = new T*[row];
        // 为每行分配空间
        for (int i=0;i<row;i++)
        {
            x[i] = new int[col];
        }
    }
    catch(bad_alloc e)
    {
        std::cerr << e.what() << '\n';
    }
}

void Hello();

#endif