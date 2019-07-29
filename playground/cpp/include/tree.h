#ifndef _TREE_H
#define _TREE_H

#include "basic.h"

namespace ADT::Tree {
    typedef struct TreeNode *PtrToTreeNode;
    typedef struct TreeNode *Position;
    typedef struct TreeNode *SearchTree;
    typedef struct TreeNode *AVLTree;

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
        // 带有平衡条件的树，深度是logN，左右高度差最多1
        // 插入节点时，导致不平衡的4种情况：
        // 1. 左儿子的左子树多了
        // 2. 左儿子的右树多了
        // 3. 右儿子的左树多了
        // 4. 右儿子的右树多了
        // 1,4 单旋转
        // 2,3 双 旋转
        AVLTree ALeft;
        AVLTree ARight;
        int AHeight;
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

    // 计算树的高度
    int Height(Position p);
    AVLTree AInsert(ElementType x, AVLTree t);
    AVLTree ADelete(ElementType x, AVLTree t);

    Position SingleRotateWithLeft(Position p);
    Position DoubleRotateWithLeft(Position p);
}
#endif