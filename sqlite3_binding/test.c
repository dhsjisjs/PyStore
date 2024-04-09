#include "sql.h"

int main()
{
    return insert("test.db", "insert into test values('test', 0)");
}
