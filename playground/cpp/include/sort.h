#ifndef _SORT_H
#define _SORT_H

#include "basic.h"

namespace ADT::Sort{
    
    // 未定义大小的数组，指定数组长度
    void InsertionSort(int a[], int size);

    // 希尔排序，增量逐渐减少
    void ShellSort(int a[],int size);
}

#endif