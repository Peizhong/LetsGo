#include <iostream>

#include "../include/common.h"
#include "../include/tree.h"

using namespace std;

void ListDir(PtrToTreeNode tree)
{
    cout<<"Id:"<<tree->Element.ID<<"value:"<<tree->Element.Value<<endl;
    if (!IsLeaf(tree))
    {
        ListDir(tree->FirstChild);
    }
    ListDir(tree->NextSibling);
}

bool IsLeaf(PtrToTreeNode node)
{
    return node->FirstChild == NULL;
}

int repeatCount = 0;
void buildBranches(PtrToTreeNode parent, int count, int level)
{
    repeatCount++;
    if (parent==NULL || count<1 || level<0)
    {
        return;
    }
    cout<<"build branch at"<<level<<"with child"<<count<<endl;
    level--;
    PtrToTreeNode previous = NULL;
    for (int c=0;c<count;c++)
    {
        PtrToTreeNode current = new TreeNode();
        current->Element.ID = c+1;
        current->Element.Value ="Demo";
        if (c==0)
        {
            parent->FirstChild = current;
        }
        else
        {
            previous->NextSibling = current;
        }
        previous = current;
        int childCount = RandomInt(4);
        buildBranches(current,childCount,level);
    }
}

TreeNode BuildDemoTree(int depth)
{
    auto hello = [] { cout << "BuildDemoTree" << endl; };
    hello();

    PtrToTreeNode root = new TreeNode();
    int childCount = RandomInt(4);
    buildBranches(root,childCount,depth);
    cout<<"repeat count"<<repeatCount<<endl;

    return *root;
}