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
    
    // 希尔排序：递减增量排序，插入排序的改进版本，单非稳定
    // 插入排序对几乎排序好的数据，效率高
    // 希尔排序的思想：将待排列的数据分割成子序列进行插入排序，待基本有序后，进行插入排序
    void ShellSort(int a[],int size)
    {
        // 12个 3,1,21,4,8,16,7,14,6,13,41,11
        // 增量为6,3,1
        int i,j,increment;
        int tmp;
        // if size=6; // 3,1
        for (increment = size/2;increment>0;increment/=2)
        {
            for (i=increment;i<size;i++)
            {
                tmp = a[i];
                // 从后面到前面，j,j-increment,j-increment*2,j-increment*3
                for (j=i;j>=increment;j-=increment)
                {
                    // 前面都排好序了，只要当前值比前面的小，一直往前进
                    if (tmp<a[j-increment])
                    {
                        a[j] = a[j-increment];
                    }
                    // 顺序正常了就可以结束了
                    else
                    {
                        break;
                    }
                }
                // 当前数据的摆放位置
                a[j]=tmp;
            }
        }
    }
}