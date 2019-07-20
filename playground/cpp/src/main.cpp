#include <iostream>
#include <vector>
#include <string>

#include "../include/list.h"

using namespace std;

// g++ -I../include/ *.cpp -g 
int main(){
    List l = new Node();
    cout<<"list"<<l<<"list.next"<<l->Next<<endl;
    Element i= Element(1,"hello");
    Position r = Find(i,l);
    cout<<"Find Result"<<r<<endl;
    delete(l);
    
    vector<string> msg {"Hello", "C++", "World", "from", "VS Code!"};
    for (const string& word : msg)
    {
        cout << word << " ";
    }
    cout << endl;
    return 0;
}