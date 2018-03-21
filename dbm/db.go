package dbm // import "sour.is/x/toolbox/dbm"

/*
Include the driver in your main package.

```
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)
```

*/

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

// GetDB eturns a database connection.
// Depricated: Use Transaction instead.
func GetDB() (*sql.Tx, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	return db.Begin()
}
