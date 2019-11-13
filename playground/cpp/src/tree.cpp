#include <iostream>
#include <string>
#include <queue>

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
        LevelOrderBinaryTree(&t);
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
    
    template <class T>
    void LevelOrderBinaryTree(BinaryTreeNode<T> *t)
    {
        queue<BinaryTreeNode<T> *> q;
        while (t!=nullptr)
        {
            // 访问根节点，左节点，右节点
            visit(t);
            // 将左，右节点入队列
            if (t->leftChild!=nullptr)
            {
                q.push(t->leftChild);
            }
            if (t->rightChild!=nullptr)
            {
                q.push(t->rightChild);
            }
            if (!q.empty())
            {
                // 取最前面的
                t = q.front();
                // 移除第一个元素
                q.pop();
            }
            else
            {
                break;
            }
        }
    }

    // AVL树高度
    template <class T>
    int height(BinaryTreeNode<T> *t)
    {
        if (t!=nullptr)
        {
            return t->height;
        }
        return 0;
    }


    template <class T>
    int max(BinaryTreeNode<T> *a,BinaryTreeNode<T> *b)
    {
        return a->height>b->height?a->height:b->height;
    }
    

    template <class T>
    BinaryTreeNode<T>* RotateLLAVLTree(BinaryTreeNode<T> *t, const T& e)
    {

    }

    template <class T>
    BinaryTreeNode<T>* RotateLRAVLTree(BinaryTreeNode<T> *t, const T& e)
    {
        
    }

    template <class T>
    BinaryTreeNode<T>* RotateRLAVLTree(BinaryTreeNode<T> *t, const T& e)
    {
        
    }

    template <class T>
    BinaryTreeNode<T>* RotateRRAVLTree(BinaryTreeNode<T> *t, const T& e)
    {
        
    }

    // AVL树
    template <class T>
    BinaryTreeNode<T>* InsertAVLTree(BinaryTreeNode<T> *t, const T& e)
    {
        if (t==nullptr)
        {
            t = new BinaryTreeNode<T>(e);
            return t;
        }
        // 不平衡的情况 LL, LR, RL, RR
        if (t->element<e)
        {
            t->leftChild = InsertAVLTree(t->leftChild,e);   
            // 插入后失去平衡
            if (height(t->leftChild)-height(t->rightChild)==2)
            {

            }
        }
        else if (t->element>e)
        {
            t->rightChild = InsertAVLTree(t->rightChild,e);
        }
        else // t->element == e
        {
            return t;
        }
        t->heigh = max(t->leftChild,t->rightChild);
    }
}