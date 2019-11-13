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

        int height;

        BinaryTreeNode()
        {
            leftChild = rightChild = nullptr;
            height = 0;
        }

        BinaryTreeNode(const T& e)
        {
            const char *p0 = "Hello world!";
            //p[3] = '3';  //error C3892: “p”: 不能给常量赋值
	        p0 = "Hi!";

            char *const p1 = "Hello world!";
            p1[3] = '3';
	        //p = "Hi!";  error C3892: “p”: 不能给常量赋值

            const char *const p2 = "Hello world!";
            //p[3] = '3';  //error C3892: “p”: 不能给常量赋值
	        //p = "Hi!";   //error C3892: “p”: 不能给常量赋值

            // const T& e，不能在函数里修改
            element = e;
            leftChild = rightChild = nullptr;
            height = 1;
        }

        BinaryTreeNode(const T& e, BinaryTreeNode *left, BinaryTreeNode *right)
        {
            element = e;
            height = 1;
            leftChild = left;
            rightChild = right;
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
    
    // 后序遍历，表达式，操作符在操作数后
    template <class T>
    void PostOrderBinaryTree(BinaryTreeNode<T> *t);

    // 层级
    template <class T>
    void LevelOrderBinaryTree(BinaryTreeNode<T> *t);

    // AVL树
    template <class T>
    void InsertAVLTree(BinaryTreeNode<T> *t, const T& e);
}
#endif