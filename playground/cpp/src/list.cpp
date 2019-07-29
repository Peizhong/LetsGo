#include <iostream>

#include "../include/list.h"

using namespace std;

namespace ADT::List{

List MakeEmpty(List l)
{
    try
    {
        cout<<l->Next<<endl;
        delete(l->Next);
    }
    catch(const std::exception& e)
    {
        std::cerr << e.what() << '\n';
    }
    cout<<"Make Empty OK"<<endl;
    return l;
}

int IsEmpty(List l)
{
    return l->Next == NULL;
}

int IsLast(Position p)
{
    return p->Next == NULL;
}

Position Find(ElementType x, List l)
{
    cout<<"In Find"<<endl;
    Position p = l->Next;
    cout<<"Find "<<p<<endl;
    while (p!=NULL && p->Item != x)
    {
        cout<<"Get Next"<<p->Next<<endl;
        p = p->Next;
    }
    return p;
}

}