#ifndef _TREE_H
#define _TREE_H

#include "basic.h"

typedef struct TreeNode *PtrToTreeNode;

struct TreeNode
{
    ElementType Element;
    PtrToTreeNode FirstChild;
    // 兄弟节点
    PtrToTreeNode NextSibling;
};

// 遍历
void ListDir(PtrToTreeNode tree);

// 没有子节点
bool IsLeaf(PtrToTreeNode node);

#endif