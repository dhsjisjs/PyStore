#include "sql.h"
#include <sqlite3.h>
#include <stdio.h>
int insert(char* filename, char* command)
{
    sqlite3* db;
    char* zErrMsg = 0;
    int rc;

    rc = sqlite3_open(filename, &db);
    if (rc) {
        fprintf(stderr, "Can't open database: %s\n", sqlite3_errmsg(db));
        sqlite3_close(db);
        return (1);
    }
    rc = sqlite3_exec(db, command, NULL, 0, &zErrMsg);
    if (rc != SQLITE_OK) {
        fprintf(stderr, "SQL error: %s\n", zErrMsg);
        sqlite3_free(zErrMsg);
        sqlite3_close(db);
        return 1;
    }
    sqlite3_close(db);
    return 0;
}
