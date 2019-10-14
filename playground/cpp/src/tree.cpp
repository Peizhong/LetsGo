#include <iostream>
#include <string>

#include "../include/common.h"
#include "../include/tree.h"

using namespace std;

namespace ADT::Tree{
    void Hi()
    {
        BinaryTreeNode<int> t = BinaryTreeNode<int>(12);
        InsertBinaryTree(&t,5);
        InsertBinaryTree(&t,20);
        InsertBinaryTree(&t,2);
        InsertBinaryTree(&t,8);
        InsertBinaryTree(&t,7);
        InsertBinaryTree(&t,1);
        InsertBinaryTree(&t,13);
        InsertBinaryTree(&t,24);
        cout<<"pre"<<endl;
        PreOrderBinaryTree(&t);
        cout<<"in"<<endl;
        InOrderBinaryTree(&t);
        cout<<"post"<<endl;
        PostOrderBinaryTree(&t);
    }
    
    template <class T>
    void visit(BinaryTreeNode<T> *t)
    {
        cout<<t->element<<endl;
    }

    // 如何记录高度?
    template <class T>
    void InsertBinaryTree(BinaryTreeNode<T> *t, const T& e)
    {
        if (e<t->element)
        {
            if (t->leftChild==nullptr)
            {
                t->leftChild = new BinaryTreeNode<T>(e);
            }
            else
            {
                InsertBinaryTree(t->leftChild,e);
            }
        }
        else if (e>t->element)
        {
            if (t->rightChild==nullptr)
            {
                t->rightChild = new BinaryTreeNode<T>(e);
            }
            else
            {
                InsertBinaryTree(t->rightChild,e);
            }
        }
    }

    // 二叉树排好序了
    template <class T>
    void PreOrderBinaryTree(BinaryTreeNode<T> *t)
    {
        visit(t);
        if (t->leftChild!=nullptr)
        {
            PreOrderBinaryTree(t->leftChild);
        }
        if (t->rightChild!=nullptr)
        {
            PreOrderBinaryTree(t->rightChild);
        }
    }
    
    template <class T>
    void InOrderBinaryTree(BinaryTreeNode<T> *t)
    {
        if (t->leftChild!=nullptr)
        {
            InOrderBinaryTree(t->leftChild);
        }
        visit(t);
        if (t->rightChild!=nullptr)
        {
            InOrderBinaryTree(t->rightChild);
        }
    }

    template <class T>
    void PostOrderBinaryTree(BinaryTreeNode<T> *t)
    {
        if (t->leftChild!=nullptr)
        {
            PostOrderBinaryTree(t->leftChild);
        }
        if (t->rightChild!=nullptr)
        {
            PostOrderBinaryTree(t->rightChild);
        }
        visit(t);
    }
}