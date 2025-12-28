package items

import "time"

type CreateItemRequest struct {
    URL          string   `json:"url"`
    Name         string   `json:"name"`
    CurrentPrice float64  `json:"current_price"`
    TargetPrice  *float64 `json:"target_price,omitempty"`
    InStock      *bool    `json:"in_stock,omitempty"`
}


type UpdateItemRequest struct {
    URL          *string   `json:"url,omitempty"`
    Name         *string   `json:"name,omitempty"`
    TargetPrice  *float64  `json:"target_price,omitempty"`
}

type UpdatePriceRequest struct {
    CurrentPrice float64  `json:"current_price"`
    InStock      bool     `json:"in_stock"`
}

type ItemResponse struct {
    ID            string     `json:"id"`
    UserID        string     `json:"user_id"`
    URL           string     `json:"url"`
    Name          string     `json:"name"`
    CurrentPrice  float64    `json:"current_price"`
    TargetPrice   *float64   `json:"target_price,omitempty"`
    InStock       bool       `json:"in_stock"`
    CreatedAt     time.Time  `json:"created_at"`
    LastCheckedAt time.Time  `json:"last_checked_at"`
}

type PriceHistoryResponse struct {
    Price     float64    `json:"price"`
    ScrapedAt time.Time  `json:"scraped_at"`
}

type PriceStatsResponse struct {
    ItemID       string   `json:"item_id"`
    MinPrice     float64  `json:"min_price"`
    MaxPrice     float64  `json:"max_price"`
    AvgPrice     float64  `json:"avg_price"`
    HistoryCount int64    `json:"history_count"`
}

type ListItemsResponse struct {
    Items []ItemResponse `json:"items"`
    Total int            `json:"total"`
}

type ErrorResponse struct {
    Error string `json:"error"`
    Code  int    `json:"code"`
}
