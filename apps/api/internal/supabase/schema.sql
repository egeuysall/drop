-- Migration: 001_create_tracked_items.sql
-- Create tracked_items table
CREATE TABLE IF NOT EXISTS tracked_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    name TEXT NOT NULL,
    current_price DECIMAL(10, 2) NOT NULL,
    target_price DECIMAL(10, 2),
    in_stock BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_checked_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_tracked_items_user_id ON tracked_items(user_id);
CREATE INDEX idx_tracked_items_last_checked ON tracked_items(last_checked_at);

-- Enable Row Level Security
ALTER TABLE tracked_items ENABLE ROW LEVEL SECURITY;

-- RLS Policies: Users can only access their own tracked items
CREATE POLICY "Users can view their own tracked items"
    ON tracked_items FOR SELECT
    USING (auth.uid() = user_id);

CREATE POLICY "Users can insert their own tracked items"
    ON tracked_items FOR INSERT
    WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own tracked items"
    ON tracked_items FOR UPDATE
    USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own tracked items"
    ON tracked_items FOR DELETE
    USING (auth.uid() = user_id);

-- Migration: 002_create_price_history.sql
-- Create price_history table
CREATE TABLE IF NOT EXISTS price_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES tracked_items(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    scraped_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_price_history_item_id ON price_history(item_id);
CREATE INDEX idx_price_history_scraped_at ON price_history(scraped_at DESC);

-- Enable Row Level Security
ALTER TABLE price_history ENABLE ROW LEVEL SECURITY;

-- RLS Policy: Users can view price history for their tracked items
CREATE POLICY "Users can view price history for their items"
    ON price_history FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM tracked_items
            WHERE tracked_items.id = price_history.item_id
            AND tracked_items.user_id = auth.uid()
        )
    );

CREATE POLICY "Service role can insert price history"
    ON price_history FOR INSERT
    WITH CHECK (true);

-- Migration: 003_create_helper_functions.sql
-- Function to automatically add price to history when tracked_item is updated
CREATE OR REPLACE FUNCTION add_price_to_history()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.current_price != OLD.current_price THEN
        INSERT INTO price_history (item_id, price, scraped_at)
        VALUES (NEW.id, NEW.current_price, NOW());
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function
CREATE TRIGGER trigger_add_price_to_history
    AFTER UPDATE OF current_price ON tracked_items
    FOR EACH ROW
    EXECUTE FUNCTION add_price_to_history();

-- Function to get price history for an item
CREATE OR REPLACE FUNCTION get_price_history(tracked_item_id UUID, days INTEGER DEFAULT 30)
RETURNS TABLE (
    price DECIMAL(10, 2),
    scraped_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT ph.price, ph.scraped_at
    FROM price_history ph
    WHERE ph.item_id = tracked_item_id
    AND ph.scraped_at >= NOW() - (days || ' days')::INTERVAL
    ORDER BY ph.scraped_at DESC;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
