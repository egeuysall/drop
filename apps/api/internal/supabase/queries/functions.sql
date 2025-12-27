-- name: GetPriceHistoryForItemWithFunction :many
-- Get price history for an item using the helper function
SELECT * FROM get_price_history($1, $2);

-- name: GetPriceHistoryByDays :many
-- Get price history for a specific tracked item within X days
SELECT * FROM price_history 
WHERE item_id = $1 
AND scraped_at >= NOW() - (make_interval(days => $2))
ORDER BY scraped_at DESC;

-- name: CheckPriceDrop :one
-- Check if current price is below target price
SELECT
    ti.id,
    ti.name,
    ti.current_price,
    ti.target_price,
    ti.current_price < ti.target_price as is_price_drop,
    ti.current_price <= ti.target_price as is_at_or_below_target
FROM tracked_items ti
WHERE ti.id = $1 AND ti.user_id = $2;
