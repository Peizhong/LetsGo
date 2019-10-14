#ifndef _TREE_H
#define _TREE_H

namespace ADT::Tree {
    void Hi();

    // 二叉树，插入、删除、搜索、遍历(前序、中序、后序)
    // AVL树
    // 红黑树
    // B树
    struct TreeNode
    {
        // Left
        // Right
        // Node
        // Parent.. 
        int Index;
        char Value[20];
    };

    template <class T>
    struct BinaryTreeNode
    {
        T element;

        BinaryTreeNode<T> *leftChild, *rightChild;

        BinaryTreeNode()
        {
            leftChild = rightChild = nullptr;
        }

        BinaryTreeNode(const T& e)
        {
            element = e;
            leftChild = rightChild = nullptr;
        }

        BinaryTreeNode(const T e, BinaryTreeNode *left, BinaryTreeNode *right)
        {
            element = e;
            leftChild = left;
            rightChild = rightChild;
        }
    };

    template <class T>
    void InsertBinaryTree(BinaryTreeNode<T> *t, const T& e);

    // 前序遍历：先访问节点，然后访问左右树
    template <class T>
    void PreOrderBinaryTree(BinaryTreeNode<T> *t);

    // 中序遍历，按大小输出
    template <class T>
    void InOrderBinaryTree(BinaryTreeNode<T> *t);
    
    // 后序遍历
    template <class T>
    void PostOrderBinaryTree(BinaryTreeNode<T> *t);
}
#endif