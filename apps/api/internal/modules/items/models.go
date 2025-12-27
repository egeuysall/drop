package items

import "time"

type TrackedItem struct {
    ID            string
    UserID        string
    URL           string
    Name          string
    CurrentPrice  float64
    TargetPrice   *float64
    InStock       *bool
    CreatedAt     time.Time
    LastCheckedAt time.Time
}

type PriceHistory struct {
    ID        string
    ItemID    string
    Price     float64
    ScrapedAt time.Time
}

type PriceStats struct {
    ItemID       string
    MinPrice     float64
    MaxPrice     float64
    AvgPrice     float64
    HistoryCount int64
}
