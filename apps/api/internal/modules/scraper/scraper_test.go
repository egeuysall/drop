package scraper

import (
	"testing"
)

func TestScrapePrice(t *testing.T) {
    scraper := NewScraper()

    url := "https://a.co/d/avMQ7hA"

    info, err := scraper.ScrapePrice(url)
    if err != nil {
        t.Logf("Error scraping: %v", err)
        return
    }

    t.Logf("Price: $%.2f", info.Price)
    t.Logf("In Stock: %v", info.InStock)
}
