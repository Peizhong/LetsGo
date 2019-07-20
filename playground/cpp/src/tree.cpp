#include <iostream>

#include "../include/tree.h"

using namespace std;

void ListDir(PtrToTreeNode tree)
{
    cout<<"Id:"<<tree->Element.ID<<"value:"<<tree->Element.Value<<endl;
    if(!IsLeaf(tree))
    {
        ListDir(tree->FirstChild);
    }
    ListDir(tree->NextSibling);
}

bool IsLeaf(PtrToTreeNode node)
{
    return node->FirstChild == NULL;
}