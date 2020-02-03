#include <iostream>
#include <cstring>

#include "../include/sort.h"
#include "../include/common.h"

using namespace std;

void printArray(const int *p, int len)
{
    for (int i = 0; i < len; i++)
    {
        cout << p[i] << " ";
    }
    cout << endl;
}

void randomArray(int *p, int len)
{
    for (int i = 0; i < len; i++)
    {
        p[i] = RandomInt(100);
    }
}

void isSorted(const int *p, int len)
{
    printArray(p, len);
    for (int i = 1; i < len; i++)
    {
        if (p[i] < p[i - 1])
        {
            cout << "wrong" << endl;
            return;
        }
    }
    cout << "ok" << endl;
}

void swap(int *a, int *b)
{
    (*a) = (*a) ^ (*b);

    (*b) = (*a) ^ (*b);
    (*a) = (*b) ^ (*a);
}

void RunSort()
{
    int arrayLen = 23;
    int *p = new int[arrayLen];
    randomArray(p, arrayLen);
    cout << "orignal" << endl;
    printArray(p, arrayLen);
    QuickSort(p, arrayLen);
    cout << "after sorted" << endl;
    isSorted(p, arrayLen);
    delete p;
}

void BubbleSort(int *p, int len)
{
    for (int i = 0; i < len; i++)
    {
        for (int j = 1; j < len - i; j++)
        {
            if (p[j] < p[j - 1])
            {
                swap(p[j], p[j - 1]);
            }
        }
        // 每轮过后，最大的会到后面
    }
}

void SelectionSort(int *p, int len)
{
    for (int i = 0; i < len; i++)
    {
        int minIndex = i;
        for (int j = i; j < len; j++)
        {
            if (p[j] < p[minIndex])
            {
                minIndex = j;
            }
        }
        if (minIndex != i)
        {
            swap(p[i], p[minIndex]);
        }
    }
}

void InsertionSort(int *p, int len)
{
    for (int i = 1; i < len; i++)
    {
        // 检查到前面已排序的
        for (int j = 0; j < i; j++)
        {
            if (p[i] < p[j])
            {
                int tmp = p[i];
                for (int k = i; k > j; k--)
                {
                    swap(p[k], p[k - 1]);
                }
                p[j] = tmp;
                break;
            }
        }
    }
}

void ShellSort(int *p, int len)
{
    int step = 1;
    while (step > (len / 3))
    {
        step = step * 3 + 1;
    }
    while (step > 0)
    {
        // 每个区间
        for (int s = 0; s < step; s++)
        {
            // 使用插入排序
            for (int i = s + step; i < len; i += step)
            {
                // 后面的值与前面的值比较
                for (int j = s; j < i; j++)
                {
                    // 如果后面的比前面的小
                    if (p[i] < p[j])
                    {
                        // 保存小的值，后面的值往后移动
                        int tmp = p[i];
                        for (int k = i; k > j; k -= step)
                        {
                            p[k] = p[k - step];
                        }
                        p[j] = tmp;
                        break;
                    }
                }
            }
        }

        step = step / 3;
    }
}

int *merge(int *left, int *right, int leftLen, int rightLen)
{
    int len = leftLen + rightLen;
    // 需要额外空间保存合并后的数据
    int *r = new int[len];
    int leftIndex = 0;
    int rightIndex = 0;
    for (int i = 0; i < len; i++)
    {
        // 从左右区块取数据比较
        if (leftIndex < leftLen && rightIndex < rightLen)
        {
            if (left[leftIndex] < right[rightIndex])
            {
                r[i] = left[leftIndex];
                leftIndex++;
            }
            else
            {
                r[i] = right[rightIndex];
                rightIndex++;
            }
        }
        // 只剩下右边
        else if (leftIndex == leftLen)
        {
            r[i] = right[rightIndex];
            rightIndex++;
        }
        // 只剩下左边
        else // if (rightIndex==rightLen)
        {
            r[i] = left[leftIndex];
            leftIndex++;
        }
    }
    delete left, right;
    return r;
}

int *divide(int *p, int len)
{
    if (len < 2)
    {
        return p;
    }
    int middle = len / 2;
    int leftLen = middle;
    int rightLen = len - middle;
    int *left = new int[leftLen];
    mempcpy(left, &p[0], leftLen * sizeof(int));
    int *right = new int[rightLen];
    mempcpy(right, &p[middle], rightLen * sizeof(int));
    // 递归将数据拆成小份排序后再合并
    return merge(divide(left, leftLen), divide(right, rightLen), leftLen, rightLen);
}

// 归并排序：拆分成小块排序后，再合并
void MergeSort(int *p, int len)
{
    int *p2 = divide(p, len);
    memcpy(p, p2, len * sizeof(len));
    delete p2;
}

void quickAdjust(int start,int end,int *p)
{
    if (start < end)
    {
        int i = start;
        int j = end;
        // 选取第一个作为基准?
        int x = p[i];
        //
        while (i < j)
        {
            // 从后往前
            while (j > i && p[j] > x)
            {
                j--;
            }
            if (j > i)
            {
                p[i] = p[j];
                i++;
            }
            // 从前往后
            while (i < j && p[i] < x)
            {
                i++;
            }
            if (i < j)
            {
                p[j] = p[i];
                j--;
            }
        }
        // 填写基准
        p[i] = x;
        // start -> i-1
        quickAdjust(start, i - 1, p);
        // i+1 -> end
        quickAdjust(start + 1, end, p);
    }
}

// 快速排序：大块排序分区后，再小块排序
void QuickSort(int *p, int len)
{
    quickAdjust(0,len,p);
}