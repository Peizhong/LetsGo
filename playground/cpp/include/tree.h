#ifndef _TREE_H
#define _TREE_H

#include "basic.h"

typedef struct TreeNode *PtrToTreeNode;
typedef struct TreeNode *Position;
typedef struct TreeNode *SearchTree;

// 实现树：每个节点除了数据本身，还有一些指针。
// 每个节点的子节点数量可以变化，不直接保存全部，而是将其放到链表
struct TreeNode
{
    ElementType Element;
    
    // 不固定数量的子节点
    PtrToTreeNode FirstChild;
    // 兄弟节点
    PtrToTreeNode NextSibling;
    
    // 二叉树的
    SearchTree Left;
    SearchTree Right;
    // AVL树，平衡的，深度logN
    // 左右子树高度差最多1
    // 插入、删除后，旋转树，保持平衡
    // 单旋转 左儿子的左树插入，右儿子的右树插入
    // 双旋转 左儿子的右树插入，右儿子的左树插入
};

// 遍历
void ListDir(PtrToTreeNode tree);

// 没有子节点
bool IsLeaf(PtrToTreeNode node);

// 造模拟数据
TreeNode BuildDemoTree(int depth);

// 初始化二叉树
SearchTree MakeEmpty(SearchTree t);

Position Find(ElementType x, SearchTree t);
Position FindMax(SearchTree t);
Position FindMin(SearchTree t);

SearchTree Insert(ElementType x, SearchTree t);
SearchTree Delete(ElementType x, SearchTree t);

// 带有平衡条件的树，深度是logN，左右高度差最多1
struct AVLTree
{

};

// 单旋转

// 双旋转

#endif