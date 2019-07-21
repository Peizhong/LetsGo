#include<ctime>
#include<cstdlib>

#include "../include/common.h"

bool bSeed;

int RandomInt(int max)
{
    if (!bSeed)
    {
        srand((unsigned)time(0));
        bSeed = true;
    }
    int r = rand();
    if (max>0)
    {
        r = r%max;
    }
    return r;
}