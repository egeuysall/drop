-- name: CreatePriceHistory :one
-- Create a new price history entry
INSERT INTO price_history (
    item_id,
    price
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetPriceHistoryForItem :many
-- Get price history for a specific tracked item
SELECT * FROM price_history
WHERE item_id = $1
ORDER BY scraped_at DESC
LIMIT $2;

-- name: GetRecentPriceHistory :many
-- Get recent price history for multiple items
SELECT ph.* FROM price_history ph
JOIN (
    SELECT ph_inner.item_id, MAX(ph_inner.scraped_at) as max_scraped_at
    FROM price_history ph_inner
    WHERE ph_inner.item_id = ANY($1)
    GROUP BY ph_inner.item_id
) latest ON ph.item_id = latest.item_id AND ph.scraped_at = latest.max_scraped_at
ORDER BY ph.item_id;

-- name: GetPriceHistoryStats :one
-- Get price statistics for a tracked item
SELECT
    item_id,
    MIN(price) as min_price,
    MAX(price) as max_price,
    AVG(price) as avg_price,
    COUNT(*) as history_count
FROM price_history
WHERE item_id = $1
GROUP BY item_id;
