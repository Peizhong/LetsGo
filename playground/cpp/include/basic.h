#ifndef _BASIC_H
#define _BASIC_H

#include <iostream>
#include <string>
#include <iterator>

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
void make2dArray(T ** &x, int row, int col)
{
    // 行指针
    // new[] 多分配4个直接保存数组大小，detele[] 配套
    x = new T*[row];
    // 为每行分配空间
    for (int i=0;i<row;i++)
    {
        x[i] = new T[col];
    }
}

template<class T>
void delete2dArray(T ** &x,int row)
{
    for (int i=0;i<row;i++)
    {
        // 删除行空间
        delete []x[i];
    }
    // 删除指针
    delete []x;
    x = NULL;
}

enum signType{sign_plus,sign_minus};

class currency
{
    public:
        // 构造
        currency(signType sign=sign_plus,unsigned long dollar=0,unsigned int cent=0);
        // 析构
        ~currency(){};

        void Set(signType,unsigned long,unsigned int);
        // 函数结尾的const，表示只能读取成员变量，不能修改。能被常量对象才能调用
        currency Add(const currency&) const;

    private:
        signType _sign;
        unsigned long _dollar;
        unsigned int _cent;
};

// 元素的所有排列
template<class T>
void permutations(T list[],int k,int m)
{
	if(k==m)
	{
		for(int i=0;i<=m;i++)
			cout<<list[i]<<" ";
		cout<<endl;
		cout<<endl;
	}
	else
    {
		for(int j=k;j<=m;j++)
		{
            cout<<"swap"<<k<<" "<<j<<endl;
			swap(list[k],list[j]);
            cout<<"recur"<<endl;
			permutations(list,k+1,m);
			swap(list[k],list[j]);
		}
    }
}

void Hello();

#endif