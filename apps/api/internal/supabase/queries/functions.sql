-- name: GetPriceHistoryForItemWithFunction :many
-- Get price history for an item using the helper function
SELECT * FROM get_price_history($1, $2);

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