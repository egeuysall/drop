To use these SQLC queries in your Go code, you'll need to follow these steps:

## 1. Import and Use the Generated Code

In your Go handlers or services:

```go
import (
    "context"
    "your-project-path/internal/supabase/queries"
    "github.com/jackc/pgx/v5"
)

func GetTrackedItems(ctx context.Context, db *pgx.Conn, userID string) ([]queries.TrackedItem, error) {
    q := queries.New(db)

    // Use the generated query
    items, err := q.ListTrackedItems(ctx, userID)
    if err != nil {
        return nil, err
    }

    return items, nil
}
```

## 2. Database Connection

Make sure you have a database connection pool:

```go
import (
    "github.com/jackc/pgx/v5/pgxpool"
    "context"
)

func main() {
    ctx := context.Background()
    conn, err := pgxpool.New(ctx, "your_connection_string")
    if err != nil {
        // handle error
    }
    defer conn.Close()

    // Pass conn to your handlers
}
```

## 3. Using in HTTP Handlers

In your API handlers:

```go
func (h *Handler) GetTrackedItems(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r) // Your auth logic

    items, err := queries.ListTrackedItems(r.Context(), h.db, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return JSON response
    json.NewEncoder(w).Encode(items)
}
```

## 4. Transaction Support

For operations that need transactions:

```go
func UpdateItemWithHistory(ctx context.Context, db *pgx.Conn, itemID, userID string, newPrice float64) error {
    tx, err := db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    q := queries.New(db)

    // Update item price
    item, err := q.UpdateTrackedItemPrice(ctx, queries.UpdateTrackedItemPriceParams{
        ID: itemID,
        CurrentPrice: newPrice,
        UserID: userID,
    })
    if err != nil {
        return err
    }

    // The trigger will automatically create price history
    // But you could also manually create if needed:
    // _, err = q.CreatePriceHistory(ctx, queries.CreatePriceHistoryParams{
    //     ItemID: itemID,
    //     Price: newPrice,
    // })

    return tx.Commit(ctx)
}
```

## Key Points:

1. **Type Safety**: SQLC generates type-safe Go code with proper parameter structs
2. **Context Support**: All queries accept `context.Context` for cancellation
3. **Error Handling**: Returns proper errors that you can handle
4. **RLS Compliance**: The queries respect your Row Level Security policies
5. **Generated Types**: You get Go structs that match your database tables

The generated code will include:

- Struct types for each table (e.g., `TrackedItem`, `PriceHistory`)
- Parameter structs for each query (e.g., `CreateTrackedItemParams`)
- Query methods on the generated `Queries` interface

You'll find the generated code follows Go conventions and integrates well with your existing Go backend.
