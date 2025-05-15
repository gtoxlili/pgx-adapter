# Pgx Adapter for Casbin

*   **Full Casbin Interface Support**: Implements `Adapter`, `ContextAdapter`, `FilteredAdapter`, `ContextFilteredAdapter`, `BatchAdapter`, `ContextBatchAdapter`, `UpdatableAdapter`, and `ContextUpdatableAdapter`.
*   **Powered by Pgx**: Utilizes the high-performance `pgx` driver for PostgreSQL.
*   **Customizable Column Length**: Allows configuration of the number of policy rule columns (v0, v1, ..., vn). Default is 6 (ptype, v0-v5).
*   **Customizable Table Name**: Allows configuration of the database table name. Default is `casbin_rule`.

## Usage

Here's a basic example of how to use the Pgx Adapter:

```go
// filepath: examples/main.go
package main

import (
    "context"
    "log"

    "github.com/casbin/casbin/v2"
    pgxadapter "github.com/gtoxlili/pgx-adapter" // Assuming this is the module path
    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    // Initialize pgx connection pool
    // Replace with your actual database connection string
    dbURL := "postgres://user:password@host:port/database"
    dbpool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer dbpool.Close()

    // Create an adapter instance
    // Default table name is "casbin_rule", default field count is 6
    // To customize, use options:
    // a, err := pgxadapter.NewAdapter(context.Background(), dbpool,
    //     pgxadapter.WithTableName("my_policy_rules"),
    //     pgxadapter.WithFieldCount(8), // For ptype, v0-v7
    // )
    a, err := pgxadapter.NewAdapter(context.Background(), dbpool)
    if err != nil {
        log.Fatalf("Failed to create adapter: %v", err)
    }

    // ... use the enforcer
}
```
The adapter will automatically create the policy table (default name: `casbin_rule`) if it doesn't exist. The table structure will adapt based on the `WithFieldCount` option.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.