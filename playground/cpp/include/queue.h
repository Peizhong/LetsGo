#ifndef _QUEUE_H
#define _QUEUE_H

#include "basic.h"

typedef struct QueueRecord *Queue;

struct QueueRecord
{
    int Capacity,Size;
    // 头，尾
    int Front,Rear;
    ElementType *Array;
};

void MakeEmpty(Queue q, int capacity);
void DisposeQueue(Queue q);

bool IsEmpty(Queue q);
bool IsFull(Queue q);
int Succ(int index, Queue q);

// Front增加
int Enqueue(Queue q, ElementType e);

// Rear增加
ElementType Dequeue(Queue q);

#endif