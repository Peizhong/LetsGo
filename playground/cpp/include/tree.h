#ifndef _TREE_H
#define _TREE_H

#include "basic.h"

typedef struct TreeNode *PtrToTreeNode;

// 实现树：每个节点除了数据本身，还有一些指针。
// 每个节点的子节点数量可以变化，不直接保存全部，而是将其放到链表
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

// 造模拟数据
TreeNode BuildDemoTree(int depth);

#endif