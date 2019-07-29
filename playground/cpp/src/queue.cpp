#include <iostream>

#include "../include/queue.h"

namespace ADT::Queue {

bool IsEmpty(Queue q)
{
    return q->Size==0;
}

bool IsFull(Queue q)
{
    return q->Size==q->Capacity;
}

void MakeEmpty(Queue q, int capacity)
{
    q->Size = 0;
    q->Capacity = capacity;
    q->Front = 1;
    q->Rear = 0;
    q->Array = new ElementType[capacity];
}

void DisposeQueue(Queue q)
{
    if (q->Array!=NULL)
    {
        // 删除数组
        delete []q->Array;
    }
    delete q;
}

int Succ(int index, Queue q)
{
    if (++index==q->Capacity)
    {
        index = 0;
    }
    return index;
}

ElementType Dequeue(Queue q)
{
    ElementType e;
    if (!IsEmpty(q))
    {
        q->Size--;
        q->Front = Succ(q->Front,q);
        e = q->Array[q->Front];
    }
    return e;
}

int Enqueue(Queue q, ElementType e)
{
    if(IsFull(q))
    {
        return -1;
    }
    q->Size++;
    q->Rear = Succ(q->Rear,q);
    q->Array[q->Rear] = e;
    return q->Size;
}

}