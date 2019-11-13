#ifndef _STACK_H
#define _STACK_H

namespace ADT::Stack {

    // 抽象类栈（基于数组、基于链表
    template<class T>
    class stack
    {
        public:
            virtual ~stack() {}
            virtual bool empty() const = 0;
            virtual int size() const = 0;
            // 返回栈顶的引用
            virtual T& top() = 0;
            // 删除栈顶
            virtual pop() = 0;
            // 
            virtual void push(const T& element) = 0;
    };

    template<class T>
    class linkedStack : public stack<T>
    {
        public:
            linkedStack(){}

        private:
            int stackSize;

    };
}

#endif