#include "sql.h"

int main()
{
    return insert("test.db", "select * from test");
}
