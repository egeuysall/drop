-- name: CreateTrackedItem :one
-- Create a new tracked item
INSERT INTO tracked_items (
    user_id, 
    url, 
    name, 
    current_price, 
    target_price, 
    in_stock
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTrackedItem :one
-- Get a tracked item by ID
SELECT * FROM tracked_items 
WHERE id = $1 AND user_id = $2;

-- name: ListTrackedItems :many
-- List all tracked items for a user
SELECT * FROM tracked_items 
WHERE user_id = $1 
ORDER BY created_at DESC;

-- name: UpdateTrackedItem :one
-- Update a tracked item
UPDATE tracked_items SET
    url = $2,
    name = $3,
    current_price = $4,
    target_price = $5,
    in_stock = $6,
    last_checked_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteTrackedItem :exec
-- Delete a tracked item
DELETE FROM tracked_items 
WHERE id = $1 AND user_id = $2;

-- name: GetTrackedItemsDueForCheck :many
-- Get tracked items that need price checking (older than 6 hours)
SELECT * FROM tracked_items 
WHERE last_checked_at < NOW() - INTERVAL '6 hours'
ORDER BY last_checked_at ASC;

-- name: UpdateTrackedItemPrice :one
-- Update only the price of a tracked item
UPDATE tracked_items SET
    current_price = $2,
    last_checked_at = NOW()
WHERE id = $1 AND user_id = $3
RETURNING *;