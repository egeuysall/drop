package scraper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PriceInfo struct {
	Price   float64
	InStock bool
}

type Scraper struct {
	client *http.Client
}

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Scraper) ScrapePrice(url string) (*PriceInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	priceText := s.extractPrice(doc)
	if priceText == "" {
		return nil, fmt.Errorf("item out of stock")
	}

	price, err := s.parsePrice(priceText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price: %w", err)
	}

	// Price found = in stock
	return &PriceInfo{
		Price:   price,
		InStock: true,
	}, nil
}

func (s *Scraper) extractPrice(doc *goquery.Document) string {
	// Amazon's split price format
	whole := doc.Find(".a-price-whole").First().Text()
	fraction := doc.Find(".a-price-fraction").First().Text()

	if whole != "" {
		// Remove any existing periods from whole
		whole = strings.ReplaceAll(strings.TrimSpace(whole), ".", "")

		if fraction != "" {
			fraction = strings.TrimSpace(fraction)
			return whole + "." + fraction
		}

		return whole
	}

	return ""
}

func (s *Scraper) parsePrice(priceText string) (float64, error) {
	priceText = strings.TrimSpace(priceText)
	priceText = strings.ReplaceAll(priceText, "$", "")
	priceText = strings.ReplaceAll(priceText, "£", "")
	priceText = strings.ReplaceAll(priceText, "€", "")
	priceText = strings.ReplaceAll(priceText, ",", "")

	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return price, nil
}
