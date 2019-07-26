#include <iostream>
#include <string>

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
    if (tree->NextSibling!=NULL)
    {
        ListDir(tree->NextSibling);
    }
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
    int nextLevel = level-1;
    PtrToTreeNode previous = NULL;
    for (int c=0;c<count;c++)
    {
        PtrToTreeNode current = new TreeNode();
        current->Element.ID = level*10+c;
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
        buildBranches(current,childCount,nextLevel);
    }
    level = nextLevel;
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


SearchTree MakeEmpty(SearchTree t)
{
    if (t!=NULL)
    {
        MakeEmpty(t->Left);
        MakeEmpty(t->Right);
        delete t;
    }
    return NULL;
}

Position Find(ElementType x, SearchTree t)
{
    if (t==NULL)
    {
        return NULL;
    }
    if (x== t->Element)
    {
        return t;
    }
    else if (x<t->Element)
    {
        return Find(x,t->Left);
    }
    // else if (x>t->Element)
    {
        return Find(x,t->Right);
    }
}

Position FindMin(SearchTree t)
{
    if (t==NULL)
    {
        return NULL;
    }
    else if (t->Left==NULL)
    {
        return t->Left;
    }
    else
    {
        return FindMin(t->Left);
    }
}

Position FindMax(SearchTree t)
{
    if (t==NULL)
    {
        return NULL;
    }
    while (t->Right!=NULL)
    {
        t=t->Right;
    }
    return t;
}

// 不平衡的
SearchTree Insert(ElementType x, SearchTree t)
{
    if (t==NULL)
    {
        t = new TreeNode();
        t->Element = x;
        t->Left = t->Right = NULL;
    }
    else if (x<t->Element)
    {
        t->Left = Insert(x,t->Left);
    }
    else if (x>t->Element)
    {
        t->Right = Insert(x,t->Right);
    }
    else
    {
        // already in it, do nothing
    }
    return t;
}

SearchTree Delete(ElementType x, SearchTree t)
{
    if (t==NULL)
    {
        return NULL;
    }
    Position tempNode = NULL;
    if (x<t->Element)
    {
        t->Left = Delete(x,t->Left);
    }
    else if (x>t->Element)
    {
        t->Right = Delete(x,t->Right);
    }
    // 找到要删除的值
    else if (t->Left && t->Right)
    {
        // 如果有两个子节点：用右树最小数据代替该节点，然后删除原最小节点
        tempNode = FindMin(t->Right);
        t->Element = tempNode->Element;
        t->Right = Delete(t->Element,t->Right);
    }
    else
    {
        // 如果是叶子节点，直接删除
        // 如果有一个子节点，用子节点代替
        tempNode = t;
        if (t->Left==NULL)
        {
            t = t->Right;
        }
        else if (t->Right==NULL)
        {
            t = t->Left;
        }
        delete tempNode;
    }
    return t;
}