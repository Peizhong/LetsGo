#include <iostream>
#include "../include/sort.h"

using namespace std;

namespace ADT::Sort{
    void InsertionSort(int a[], int size)
    {
        // 取不到数组长度?
        size_t arraysize = sizeof(*a) / sizeof(int);
        if (arraysize>size)
        {
           // size = arraysize;
        }
        int j,p;
        int tmp;
        for (p=1;p<size;p++)
        {
            tmp = a[p];
            // 前面的比当前值，将前面的值移到后面
            for (j=p;j>0 && a[j-1]>tmp;j--)
            {
                a[j] = a[j-1];
            }
            // 最后的j，位置没有变，或往前移动了
            a[j] = tmp;
        }
    }
    
    void ShellSort(int a[],int size)
    {
        int i,j,increment;
        int tmp;
        // if size=6; // 3,1
        for (increment = size/2;increment>0;increment/=2)
        {
            // 3,4,5,6
            // 1,2,3,4,5,6
            for (i=increment;i<size;i++)
            {
                tmp = a[i];
                // 从后面到前面
                for (j=i;j>=increment;j-=increment)
                {
                    if (tmp<a[j-increment])
                    {
                        a[j] = a[j-increment];
                    }
                    else
                    {
                        break;
                    }
                }
                a[j]=tmp;
            }
        }
    }
}