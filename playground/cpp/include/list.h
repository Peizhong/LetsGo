#ifndef _LIST_H
#define _LIST_H

#include "basic.h"

typedef struct Node *PtrToNode;
typedef PtrToNode List;
typedef PtrToNode Position;

struct Node
{
    ElementType Item;
    Position Next;
};

List MakeEmpty(List l);
int IsEmpty(List l);
int IsLast(Position p);
Position Find(ElementType x, List l);

#endif