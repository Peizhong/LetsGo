#include <iostream>
#include <vector>
#include <string>

#include "../include/list.h"
#include "../include/queue.h"
#include "../include/tree.h"

using namespace std;

void Test0()
{
    int a[10] = {1,2,3,4,5,6};
    // 指向数组的指针
    int *b = a;
    int c = b[1];
    cout<<c<<"and"<<*(b+1)<<endl;

    vector<string> msg {"Hello", "C++", "World", "from", "VS Code!"};
    for (const string& word : msg)
    {
        cout << word << " ";
    }
    cout << endl;
}

// abstract data type: 抽象数据类型
void TestADT()
{
    List l = new Node();
    cout<<"list"<<l<<"list.next"<<l->Next<<endl;
    Element i= Element(1,"hello");
    Position r = Find(i,l);
    cout<<"Find Result"<<r<<endl;
    delete l;
    
    Queue q = new QueueRecord();
    MakeEmpty(q,10);
    cout<<q->Array<<endl;
    cout<<"queue"<<IsEmpty(q)<<endl;
    for (int i=0;i<20;i++)
    {
        ElementType e = Element(i,"hello");
        int r = Enqueue(q,e);
        cout<<r<<endl;
    }
    for (int i=0;i<20;i++)
    {
        ElementType e = Dequeue(q);
        cout<<e.Value<<endl;
    }
    cout<<q->Array<<endl;
    DisposeQueue(q);

    TreeNode t = BuildDemoTree(3);
    ListDir(&t);
}

// g++ -I../include/ *.cpp -g 
int main(){
    TestADT();
    return 0;
}