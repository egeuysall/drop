package items

import (
	"context"
	"fmt"
	"strings"

	"github.com/egeuysall/drop/internal/modules/scraper"
	"github.com/egeuysall/drop/internal/utils"
)

// Service defines the interface for items business logic
type Service interface {
	CreateItem(ctx context.Context, userID string, req CreateItemRequest) (*ItemResponse, error)
	GetItemByID(ctx context.Context, id, userID string) (*ItemResponse, error)
	ListItemsByUserID(ctx context.Context, userID string) ([]ItemResponse, error)
	UpdateItemByID(ctx context.Context, id, userID string, req UpdateItemRequest) (*ItemResponse, error)
	DeleteItem(ctx context.Context, id, userID string) error
	UpdateItemPrice(ctx context.Context, id, userID string, currentPrice float64, inStock bool) (*ItemResponse, error)
	GetPriceHistory(ctx context.Context, id string, days int) ([]PriceHistoryResponse, error)
	GetPriceStats(ctx context.Context, id string) (*PriceStatsResponse, error)
	CheckPriceDrop(ctx context.Context, id, userID string) (*PriceDropCheck, error)
	GetItemsDueForCheck(ctx context.Context) ([]ItemResponse, error)
	RefreshPrice(ctx context.Context, itemID, userID, url string) error
}

// service implements the Service interface
type service struct {
	repo Repository
	scraper *scraper.Scraper
}

func NewService(repo Repository, scr *scraper.Scraper) Service {
	return &service{
		repo: repo,
		scraper: scr,
	}
}

// CreateItem handles the full item creation workflow
func (s *service) CreateItem(ctx context.Context, userID string, req CreateItemRequest) (*ItemResponse, error) {
    if err := utils.ValidateURL(req.URL); err != nil {
        return nil, fmt.Errorf("invalid URL: %w", err)
    }

    if err := s.checkForDuplicates(ctx, userID, req.URL); err != nil {
        return nil, fmt.Errorf("duplicate item: %w", err)
    }

    priceInfo, err := s.scraper.ScrapePrice(req.URL)

    if err != nil {
        if strings.Contains(err.Error(), "out of stock") {
            currentPrice := 0.0

            if req.TargetPrice != nil {
                currentPrice = *req.TargetPrice
            }

            inStock := false
            req.CurrentPrice = currentPrice
            req.InStock = &inStock
        } else {
            return nil, fmt.Errorf("failed to scrape price: %w", err)
        }
    } else {
        req.CurrentPrice = priceInfo.Price
        req.InStock = &priceInfo.InStock
    }

    item := TrackedItem{
        UserID:       userID,
        URL:          req.URL,
        Name:         req.Name,
        CurrentPrice: req.CurrentPrice,
        TargetPrice:  req.TargetPrice,
        InStock:      req.InStock,
    }

    createdItem, err := s.repo.CreateItem(ctx, item)
    if err != nil {
        return nil, fmt.Errorf("failed to create item: %w", err)
    }

    return &ItemResponse{
        ID:           createdItem.ID,
        UserID:       createdItem.UserID,
        URL:          createdItem.URL,
        Name:         createdItem.Name,
        CurrentPrice: createdItem.CurrentPrice,
        TargetPrice:  createdItem.TargetPrice,
        InStock:      *createdItem.InStock,
        CreatedAt:    createdItem.CreatedAt,
    }, nil
}

// GetItemByID retrieves a single item by ID with ownership validation
func (s *service) GetItemByID(ctx context.Context, id, userID string) (*ItemResponse, error) {
    item, err := s.repo.GetItemByID(ctx, id, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get item: %w", err)
    }

    return &ItemResponse{
        ID:            item.ID,
        UserID:        item.UserID,
        URL:           item.URL,
        Name:          item.Name,
        CurrentPrice:  item.CurrentPrice,
        TargetPrice:   item.TargetPrice,
        InStock:       *item.InStock,
        CreatedAt:     item.CreatedAt,
        LastCheckedAt: item.LastCheckedAt,
    }, nil
}

// ListItemsByUserID retrieves all items for a specific user
func (s *service) ListItemsByUserID(ctx context.Context, userID string) ([]ItemResponse, error) {
    items, err := s.repo.ListItemsByUserID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to list items: %w", err)
    }

    responses := make([]ItemResponse, len(items))
    for i, item := range items {
        responses[i] = ItemResponse{
            ID:            item.ID,
            UserID:        item.UserID,
            URL:           item.URL,
            Name:          item.Name,
            CurrentPrice:  item.CurrentPrice,
            TargetPrice:   item.TargetPrice,
            InStock:       *item.InStock,
            CreatedAt:     item.CreatedAt,
            LastCheckedAt: item.LastCheckedAt,
        }
    }

    return responses, nil
}

// UpdateItemByID updates an existing item with validation
func (s *service) UpdateItemByID(ctx context.Context, id, userID string, req UpdateItemRequest) (*ItemResponse, error) {
    // Validate URL if provided
    if req.URL != nil {
        if err := utils.ValidateURL(*req.URL); err != nil {
            return nil, fmt.Errorf("invalid URL: %w", err)
        }
        // Check for duplicates if URL is being changed
        if err := s.checkForDuplicateURLChange(ctx, userID, id, *req.URL); err != nil {
            return nil, fmt.Errorf("duplicate item: %w", err)
        }
    }

    updatedItem, err := s.repo.UpdateItemByID(ctx, id, userID, req)
    if err != nil {
        return nil, fmt.Errorf("failed to update item: %w", err)
    }

    return &ItemResponse{
        ID:            updatedItem.ID,
        UserID:        updatedItem.UserID,
        URL:           updatedItem.URL,
        Name:          updatedItem.Name,
        CurrentPrice:  updatedItem.CurrentPrice,
        TargetPrice:   updatedItem.TargetPrice,
        InStock:       *updatedItem.InStock,
        CreatedAt:     updatedItem.CreatedAt,
        LastCheckedAt: updatedItem.LastCheckedAt,
    }, nil
}

// DeleteItem removes an item with ownership validation
func (s *service) DeleteItem(ctx context.Context, id, userID string) error {
    return s.repo.DeleteItem(ctx, id, userID)
}

// UpdateItemPrice updates the price and stock status with validation
func (s *service) UpdateItemPrice(ctx context.Context, id, userID string, currentPrice float64, inStock bool) (*ItemResponse, error) {
    // Validate price
    if currentPrice < 0 {
        return nil, fmt.Errorf("price cannot be negative")
    }

    updatedItem, err := s.repo.UpdateItemPrice(ctx, id, userID, currentPrice, inStock)
    if err != nil {
        return nil, fmt.Errorf("failed to update price: %w", err)
    }

    return &ItemResponse{
        ID:            updatedItem.ID,
        UserID:        updatedItem.UserID,
        URL:           updatedItem.URL,
        Name:          updatedItem.Name,
        CurrentPrice:  updatedItem.CurrentPrice,
        TargetPrice:   updatedItem.TargetPrice,
        InStock:       *updatedItem.InStock,
        CreatedAt:     updatedItem.CreatedAt,
        LastCheckedAt: updatedItem.LastCheckedAt,
    }, nil
}

// GetPriceHistory retrieves price history with parameter validation
func (s *service) GetPriceHistory(ctx context.Context, id string, days int) ([]PriceHistoryResponse, error) {
    // Validate days parameter
    if days <= 0 || days > 365 {
        return nil, fmt.Errorf("days must be between 1 and 365")
    }

    history, err := s.repo.GetPriceHistory(ctx, id, days)
    if err != nil {
        return nil, fmt.Errorf("failed to get price history: %w", err)
    }

    responses := make([]PriceHistoryResponse, len(history))
    for i, entry := range history {
        responses[i] = PriceHistoryResponse{
            Price:     entry.Price,
            ScrapedAt: entry.ScrapedAt,
        }
    }

    return responses, nil
}

// GetPriceStats retrieves price statistics for an item
func (s *service) GetPriceStats(ctx context.Context, id string) (*PriceStatsResponse, error) {
    stats, err := s.repo.GetPriceStats(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get price stats: %w", err)
    }

    return &PriceStatsResponse{
        ItemID:       stats.ItemID,
        MinPrice:     stats.MinPrice,
        MaxPrice:     stats.MaxPrice,
        AvgPrice:     stats.AvgPrice,
        HistoryCount: stats.HistoryCount,
    }, nil
}

// CheckPriceDrop checks if an item's price has dropped below target
func (s *service) CheckPriceDrop(ctx context.Context, id, userID string) (*PriceDropCheck, error) {
    return s.repo.CheckPriceDrop(ctx, id, userID)
}

// GetItemsDueForCheck retrieves items that need price checking
func (s *service) GetItemsDueForCheck(ctx context.Context) ([]ItemResponse, error) {
    items, err := s.repo.GetItemsDueForCheck(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get items due for check: %w", err)
    }

    responses := make([]ItemResponse, len(items))
    for i, item := range items {
        responses[i] = ItemResponse{
            ID:            item.ID,
            URL:           item.URL,
            Name:          item.Name,
            CurrentPrice:  item.CurrentPrice,
            TargetPrice:   item.TargetPrice,
            InStock:       *item.InStock,
            CreatedAt:     item.CreatedAt,
            LastCheckedAt: item.LastCheckedAt,
        }
    }

    return responses, nil
}

// checkForDuplicates checks if user already tracks this URL
func (s *service) checkForDuplicates(ctx context.Context, userID, url string) error {
    items, err := s.repo.ListItemsByUserID(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to check for duplicates: %w", err)
    }

    for _, item := range items {
        if item.URL == url {
            return fmt.Errorf("item with this URL already exists")
        }
    }

    return nil
}

// checkForDuplicateURLChange checks for duplicates when updating URL
func (s *service) checkForDuplicateURLChange(ctx context.Context, userID, currentItemID, newURL string) error {
    items, err := s.repo.ListItemsByUserID(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to check for duplicates: %w", err)
    }

    for _, item := range items {
        if item.ID != currentItemID && item.URL == newURL {
            return fmt.Errorf("item with this URL already exists")
        }
    }

    return nil
}

func (s *service) RefreshPrice(ctx context.Context, itemID, userID, url string) error {
    priceInfo, err := s.scraper.ScrapePrice(url)

    if err != nil {
        if strings.Contains(err.Error(), "out of stock") {
            _, err := s.repo.UpdateItemPrice(ctx, itemID, userID, 0, false)
            return err
        }

        return fmt.Errorf("failed to scrape price: %w", err)
    }

    _, err = s.repo.UpdateItemPrice(ctx, itemID, userID, priceInfo.Price, priceInfo.InStock)

    if err != nil {
        return fmt.Errorf("failed to update price: %w", err)
    }

    return nil
}
