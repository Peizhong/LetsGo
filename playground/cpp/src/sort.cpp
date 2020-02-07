#include <iostream>
#include <cstring>
#include <set> 
#include <map> 

#include "../include/sort.h"
#include "../include/common.h"

#define GOODS 5

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

void quickAdjust(int start, int end, int *p)
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
    quickAdjust(0, len, p);
}

void printTable(int **p, int rows, int columns)
{
    for (int c = 0; c < columns; c++)
    {
        cout << c << '\t';
    }
    cout << endl;
    cout << "----------" << endl;
    for (int r = 0; r < rows; r++)
    {
        for (int c = 0; c < columns; c++)
        {
            cout << p[r][c] << '\t';
        }
        cout << endl;
    }
}

int goodsWeight[GOODS] = {1, 2, 9, 10, 11};
int goodsValue[GOODS] = {6, 7, 10, 11, 12};

void resolve(int **p, int *r ,int cap,int n)
{
    if (n>0)
    {
        // 如果前后2行的差等于n的价格，说明使用了
        if ((cap>=goodsWeight[n-1]) && (p[n][cap]==p[n-1][cap-goodsWeight[n-1]]+goodsValue[n-1]))
        {
            r[n-1] = 1;
            resolve(p,r,cap-goodsWeight[n-1],n-1);
        }
        // 如果值相等，表示没有使用第n个
        else // if (p[n][cap]==p[n-1][cap])
        {
            // 没有增加
            r[n-1] = 0;
            resolve(p,r,cap,n-1);
        }
    }
}

// 最背包解
void resolveKnap(int **p, int cap,int n)
{
    int *r = new int[n];
    resolve(p,r,cap,n);
    printArray(r,n);
    delete r;
}

// 动态规划
void Knapsack(int capacity)
{
    int colums = capacity + 1;
    int rows = GOODS + 1;
    // 初始化数据
    int **table = new int *[rows];
    for (int i = 0; i < rows; i++)
    {
        table[i] = new int[colums];
        memset(table[i], 0, colums * sizeof(int));
    }
    // 填充数据，从1开始
    for (int r = 1; r < rows; r++)
    {
        for (int c = 1; c < colums; c++)
        {
            // 没有容量或不使用
            int value = table[r - 1][c];
            // 如果每个商品可以选择多个
            if (c >= goodsWeight[r - 1])
            {
                int maxCap = c / goodsWeight[r-1];
                maxCap =1;
                // 用n个
                for (int n = 1; n <= maxCap; n++)
                {
                    int value1 = goodsValue[r - 1] * n + table[r - 1][c - goodsWeight[r - 1] * n];
                    value = max(value, value1);
                }
            }
            table[r][c] = value;
        }
    }
    printTable(table, rows, colums);
    resolveKnap(table,capacity,GOODS);
    // 释放资源
    for (int i = 0; i < rows; i++)
    {
        delete table[i];
    }
    delete table;

    int *dp = new int[colums];
    memset(dp, 0, colums * sizeof(int));
    for (int i = 0; i < GOODS; i++)
    {
        for (int v = capacity; v >= goodsWeight[i]; v--)
        {
            int maxCap = v / goodsWeight[i];
            for (int n = 1; n <= maxCap; n++)
            {
                dp[v] = max(dp[v], dp[v - goodsWeight[i] * n] + goodsValue[i] * n);
            }
        }
    }
    printArray(dp, colums);
}

// 贪心算法
void Broadcasts()
{
    string s[8] = {"ID","NV","UT","WA","MT","OR","CA","AZ"};
    set<string> st;
    for (int i =0;i<8;i++)
    {
        st.insert(s[i]);
    }
    map<string,string*> station0;
    string k0[3] = {"ID","NV","UT"};
    string k1[3] = {"WA","ID","MT"};
    string k2[3] = {"OR","NV","CA"};
    string k3[2] = {"NV","UT"};
    string k4[2] = {"CA","AZ"};
    station0["k0"] = k0;
    station0["k1"] = k1;
    station0["k2"] = k2;
    station0["k3"] = k3;
    station0["k4"] = k4;
}