#ifndef _BASIC_H
#define _BASIC_H
 
#include <string>

typedef struct Element ElementType;

struct Element
{
    int ID;
    std::string Value;

    Element()
    {

    }

    Element(int id,std::string value)
    {
        ID = id;
        Value = value;
    }

    bool operator!=(Element const& b) const
    {
        return true;
    }
    /*
    bool operator==(Element const& b) const
    {
        return true;
    }
    bool operator<(Element const& b) const
    {
        return true;
    }
     */
};
#endif