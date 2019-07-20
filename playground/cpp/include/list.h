#ifndef _LIST_H

#include <string>

using namespace std;

typedef struct Node *PtrToNode;
typedef PtrToNode List;
typedef PtrToNode Position;

typedef struct Element ElementType;

struct Element
{
    int ID;
    string Value;

    Element()
    {

    }

    Element(int id,string value)
    {
        ID = id;
        Value = value;
    }

    bool operator!=(Element const& b) const
    {
        return true;
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