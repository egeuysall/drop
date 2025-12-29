package items

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	generated "github.com/egeuysall/drop/internal/supabase/generated"
	"github.com/egeuysall/drop/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
    CreateItem(ctx context.Context, item TrackedItem) (*TrackedItem, error)
    GetItemByID(ctx context.Context, id, userID string) (*TrackedItem, error)
    ListItemsByUserID(ctx context.Context, userID string) ([]TrackedItem, error)
    UpdateItemByID(ctx context.Context, id, userID string, req UpdateItemRequest) (*TrackedItem, error)
    DeleteItem(ctx context.Context, id, userID string) error
    UpdateItemPrice(ctx context.Context, id, userID string, currentPrice float64, inStock bool) (*TrackedItem, error)
    GetPriceHistory(ctx context.Context, id string, days int) ([]PriceHistory, error)
    GetPriceStats(ctx context.Context, id string) (*PriceStats, error)
    CheckPriceDrop(ctx context.Context, id, userID string) (*PriceDropCheck, error)
    GetItemsDueForCheck(ctx context.Context) ([]TrackedItem, error)
}

type repository struct {
    queries *generated.Queries
}

func NewRepository(queries *generated.Queries) Repository {
    return &repository{
        queries: queries,
    }
}

func (r *repository) CreateItem(ctx context.Context, item TrackedItem) (*TrackedItem, error) {
    // Try to parse user ID as UUID, if it fails, generate a UUID from the string
    var userID pgtype.UUID
    if parsedUUID, err := utils.ParseUUID(item.UserID); err == nil {
        userID = parsedUUID
    } else {
        // Generate a deterministic UUID from the string user ID
        // This ensures we can always map the string ID to the same UUID
        hash := sha256.Sum256([]byte(item.UserID))
        userIDBytes := hash[:16] // Use first 16 bytes for UUID
        var uuid pgtype.UUID
        copy(uuid.Bytes[:], userIDBytes)
        uuid.Valid = true
        userID = uuid
    }

    params := generated.CreateTrackedItemParams{
        UserID: userID,
        Url: item.URL,
        Name: item.Name,
        CurrentPrice: utils.Float64ToNumeric(item.CurrentPrice),
        TargetPrice: utils.Float64PtrToNumeric(item.TargetPrice),
        InStock: utils.BoolPtrToPgBool(item.InStock),
    }

    result, err := r.queries.CreateTrackedItem(ctx, params)

    if err != nil {
        return nil, fmt.Errorf("failed to create item: %w", err)
    }

    return &TrackedItem{
        ID: utils.UUIDToString(result.ID),
        UserID: utils.UUIDToString(result.UserID),
        URL: result.Url,
        Name: result.Name,
        CurrentPrice: utils.NumericToFloat64(result.CurrentPrice),
        TargetPrice: utils.NumericToFloat64Ptr(result.TargetPrice),
        InStock: utils.PgBoolToBoolPtr(result.InStock),
        CreatedAt: result.CreatedAt.Time,
        LastCheckedAt: result.LastCheckedAt.Time,
    }, nil
}

func (r *repository) GetItemByID(ctx context.Context, id, userID string) (*TrackedItem, error) {
    // Try to parse user ID as UUID, if it fails, generate a UUID from the string
    var parsedUserID pgtype.UUID
    if uuid, err := utils.ParseUUID(userID); err == nil {
        parsedUserID = uuid
    } else {
        // Generate a deterministic UUID from the string user ID
        hash := sha256.Sum256([]byte(userID))
        userIDBytes := hash[:16] // Use first 16 bytes for UUID
        var uuid pgtype.UUID
        copy(uuid.Bytes[:], userIDBytes)
        uuid.Valid = true
        parsedUserID = uuid
    }

    parsedID, err := utils.ParseUUID(id)

    if err != nil {
        return nil, fmt.Errorf("invalid item ID: %w", err)
    }

    params := generated.GetTrackedItemParams{
        UserID: parsedUserID,
        ID: parsedID,
    }

    result, err := r.queries.GetTrackedItem(ctx, params)

    if err != nil {
        return nil, fmt.Errorf("failed to get item: %w", err)
    }

    return &TrackedItem{
        ID: utils.UUIDToString(result.ID),
        UserID: utils.UUIDToString(result.UserID),
        URL: result.Url,
        Name: result.Name,
        CurrentPrice: utils.NumericToFloat64(result.CurrentPrice),
        TargetPrice: utils.NumericToFloat64Ptr(result.TargetPrice),
        InStock: utils.PgBoolToBoolPtr(result.InStock),
        CreatedAt: result.CreatedAt.Time,
        LastCheckedAt: result.LastCheckedAt.Time,
    }, nil
}

func (r *repository) ListItemsByUserID(ctx context.Context, userID string) ([]TrackedItem, error) {
    parsedUserID, err := utils.ParseUUID(userID)

    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    results, err := r.queries.ListTrackedItems(ctx, parsedUserID)

    if err != nil {
        return nil, fmt.Errorf("failed to list items: %w", err)
    }

    items := make([]TrackedItem, len(results))

    for i, item := range results {
        items[i] = TrackedItem{
            ID: utils.UUIDToString(item.ID),
            UserID: utils.UUIDToString(item.UserID),
            URL: item.Url,
            Name: item.Name,
            CurrentPrice: utils.NumericToFloat64(item.CurrentPrice),
            TargetPrice: utils.NumericToFloat64Ptr(item.TargetPrice),
            InStock: utils.PgBoolToBoolPtr(item.InStock),
            CreatedAt: item.CreatedAt.Time,
            LastCheckedAt: item.LastCheckedAt.Time,
        }
    }

    return items, nil
}

func (r *repository) UpdateItemByID(ctx context.Context, id, userID string, req UpdateItemRequest) (*TrackedItem, error) {
    parsedUserID, err := utils.ParseUUID(userID)

    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    parsedID, err := utils.ParseUUID(id)

    if err != nil {
        return nil, fmt.Errorf("invalid item ID: %w", err)
    }

    current, err := r.GetItemByID(ctx, id, userID)

    if err != nil {
        return nil, fmt.Errorf("failed to get current item: %w", err)
    }

    url := current.URL

    if req.URL != nil {
        url = *req.URL
    }

    name := current.Name

    if req.Name != nil {
        name = *req.Name
    }

    targetPrice := current.TargetPrice
    if req.TargetPrice != nil {
        targetPrice = req.TargetPrice
    }

    params := generated.UpdateTrackedItemParams{
        ID:           parsedID,
        Url:          url,
        Name:         name,
        CurrentPrice: utils.Float64ToNumeric(current.CurrentPrice),
        TargetPrice:  utils.Float64PtrToNumeric(targetPrice),
        InStock:      utils.BoolPtrToPgBool(current.InStock),
        UserID:       parsedUserID,
    }

    result, err := r.queries.UpdateTrackedItem(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("failed to update item: %w", err)
    }

    return &TrackedItem{
        ID:            utils.UUIDToString(result.ID),
        UserID:        utils.UUIDToString(result.UserID),
        URL:           result.Url,
        Name:          result.Name,
        CurrentPrice:  utils.NumericToFloat64(result.CurrentPrice),
        TargetPrice:   utils.NumericToFloat64Ptr(result.TargetPrice),
        InStock:       utils.PgBoolToBoolPtr(result.InStock),
        CreatedAt:     result.CreatedAt.Time,
        LastCheckedAt: result.LastCheckedAt.Time,
    }, nil
}

func (r *repository) DeleteItem(ctx context.Context, id, userID string) error {
    parsedUserID, err := utils.ParseUUID(userID)

    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }

    parsedID, err := utils.ParseUUID(id)

    if err != nil {
        return fmt.Errorf("invalid item ID: %w", err)
    }

    params := generated.DeleteTrackedItemParams{
        ID: parsedID,
        UserID: parsedUserID,
    }

    if err = r.queries.DeleteTrackedItem(ctx, params); err != nil {
        return fmt.Errorf("failed to delete item: %w", err)
    }

    return nil
}

func (r *repository) UpdateItemPrice(ctx context.Context, id, userID string, currentPrice float64, inStock bool) (*TrackedItem, error) {
    log.Printf("UpdateItemPrice called: id=%s, userID=%s, currentPrice=%.2f, inStock=%t", id, userID, currentPrice, inStock)

    // Try to parse user ID as UUID, if it fails, generate a UUID from the string
    var parsedUserID pgtype.UUID
    if uuid, err := utils.ParseUUID(userID); err == nil {
        parsedUserID = uuid
        log.Printf("Parsed userID as UUID: %s -> %s", userID, uuid)
    } else {
        // Generate a deterministic UUID from the string user ID
        hash := sha256.Sum256([]byte(userID))
        userIDBytes := hash[:16] // Use first 16 bytes for UUID
        var uuid pgtype.UUID
        copy(uuid.Bytes[:], userIDBytes)
        uuid.Valid = true
        parsedUserID = uuid
        log.Printf("Generated UUID from string userID: %s -> %s", userID, uuid)
    }

    parsedID, err := utils.ParseUUID(id)

    if err != nil {
        return nil, fmt.Errorf("invalid item ID: %w", err)
    }

    params := generated.UpdateTrackedItemPriceParams{
        ID: parsedID,
        UserID: parsedUserID,
        CurrentPrice: utils.Float64ToNumeric(currentPrice),
        InStock: utils.BoolToPgBool(inStock),
    }

    result, err := r.queries.UpdateTrackedItemPrice(ctx, params)

    if err != nil {
        return nil, fmt.Errorf("failed to update item price: %w", err)
    }

    return &TrackedItem{
        ID: utils.UUIDToString(result.ID),
        UserID: utils.UUIDToString(result.UserID),
        URL: result.Url,
        Name: result.Name,
        CurrentPrice: utils.NumericToFloat64(result.CurrentPrice),
        TargetPrice: utils.NumericToFloat64Ptr(result.TargetPrice),
        InStock: utils.PgBoolToBoolPtr(result.InStock),
        CreatedAt: result.CreatedAt.Time,
        LastCheckedAt: result.LastCheckedAt.Time,
    }, nil
}

func (r *repository) GetPriceHistory(ctx context.Context, id string, days int) ([]PriceHistory, error) {
    parsedID, err := utils.ParseUUID(id)
    if err != nil {
        return nil, fmt.Errorf("invalid item ID: %w", err)
    }

    params := generated.GetPriceHistoryByDaysParams{
        ItemID: parsedID,
        Days: int32(days),
    }

    generatedHistory, err := r.queries.GetPriceHistoryByDays(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("failed to get price history: %w", err)
    }

    history := make([]PriceHistory, len(generatedHistory))
    for i, gh := range generatedHistory {
        history[i] = PriceHistory{
            ID: utils.UUIDToString(gh.ID),
            ItemID: utils.UUIDToString(gh.ItemID),
            Price: utils.NumericToFloat64(gh.Price),
            ScrapedAt: gh.ScrapedAt.Time,
        }
    }

    return history, nil
}

func (r *repository) GetPriceStats(ctx context.Context, id string) (*PriceStats, error) {
    parsedID, err := utils.ParseUUID(id)
    if err != nil {
        return nil, fmt.Errorf("invalid item ID: %w", err)
    }

    generatedStats, err := r.queries.GetPriceHistoryStats(ctx, parsedID)
    if err != nil {
        return nil, fmt.Errorf("failed to get price stats: %w", err)
    }

    // Handle min price - could be pgtype.Numeric or nil
    var minPrice float64
    switch v := generatedStats.MinPrice.(type) {
    case pgtype.Numeric:
        minPrice = utils.NumericToFloat64(v)
    case nil:
        minPrice = 0.0
    default:
        return nil, fmt.Errorf("unsupported min price type: %T", v)
    }

    // Handle max price - could be pgtype.Numeric or nil
    var maxPrice float64
    switch v := generatedStats.MaxPrice.(type) {
    case pgtype.Numeric:
        maxPrice = utils.NumericToFloat64(v)
    case nil:
        maxPrice = 0.0
    default:
        return nil, fmt.Errorf("unsupported max price type: %T", v)
    }

    // AvgPrice is already a float64
    avgPrice := generatedStats.AvgPrice

    stats := PriceStats{
        ItemID:       utils.UUIDToString(generatedStats.ItemID),
        MinPrice:     minPrice,
        MaxPrice:     maxPrice,
        AvgPrice:     avgPrice,
        HistoryCount: generatedStats.HistoryCount,
    }

    return &stats, nil
}

func (r *repository) CheckPriceDrop(ctx context.Context, id, userID string) (*PriceDropCheck, error) {
    parsedUserID, err := utils.ParseUUID(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }

    parsedID, err := utils.ParseUUID(id)
    if err != nil {
        return nil, fmt.Errorf("invalid item ID: %w", err)
    }

    params := generated.CheckPriceDropParams{
        ID: parsedID,
        UserID: parsedUserID,
    }

    generatedCheck, err := r.queries.CheckPriceDrop(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("failed to check price drop: %w", err)
    }

    check := PriceDropCheck{
        ID: utils.UUIDToString(generatedCheck.ID),
        Name: generatedCheck.Name,
        CurrentPrice: utils.NumericToFloat64(generatedCheck.CurrentPrice),
        TargetPrice: utils.NumericToFloat64Ptr(generatedCheck.TargetPrice),
        IsPriceDrop: generatedCheck.IsPriceDrop,
        IsAtOrBelowTarget: generatedCheck.IsAtOrBelowTarget,
    }

    return &check, nil
}

func (r *repository) GetItemsDueForCheck(ctx context.Context) ([]TrackedItem, error) {
    generatedItems, err := r.queries.GetTrackedItemsDueForCheck(ctx)

    if err != nil {
        return nil, fmt.Errorf("failed to get items due for check: %w", err)
    }

    items := make([]TrackedItem, len(generatedItems))
    for i, item := range generatedItems {
        items[i] = TrackedItem{
            ID:            utils.UUIDToString(item.ID),
            UserID:        utils.UUIDToString(item.UserID),
            URL:           item.Url,
            Name:          item.Name,
            CurrentPrice:  utils.NumericToFloat64(item.CurrentPrice),
            TargetPrice:   utils.NumericToFloat64Ptr(item.TargetPrice),
            InStock:       utils.PgBoolToBoolPtr(item.InStock),
            CreatedAt:     item.CreatedAt.Time,
            LastCheckedAt: item.LastCheckedAt.Time,
        }
    }
    return items, nil
}
