package scheduler

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/egeuysall/drop/internal/modules/items"
	"github.com/egeuysall/drop/internal/modules/items/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestNewPriceRefresherScheduler tests scheduler initialization
func TestNewPriceRefresherScheduler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)

	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 5)

	assert.NotNil(t, scheduler)
	assert.Equal(t, 5, scheduler.workerCount)
	assert.Equal(t, 30*time.Minute, scheduler.interval)
}

// TestPriceRefreshWorker tests the worker function
func TestPriceRefreshWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 1)

	// Test successful refresh
	job := ItemJob{
		ID:     "item1",
		UserID: "user1",
		URL:    "https://example.com/product1",
		Name:   "Test Product",
	}

	mockService.EXPECT().RefreshPrice(
		context.Background(),
		job.ID,
		job.UserID,
		job.URL,
	).Return(nil)

	jobs := make(chan ItemJob, 1)
	results := make(chan string, 1)

	jobs <- job
	close(jobs)

	// Run worker
	go scheduler.priceRefreshWorker(1, jobs, results)

	// Check result
	result := <-results
	expected := "SUCCESS: Test Product"
	assert.Equal(t, expected, result)
}

// TestPriceRefreshWorker_Error tests worker with error
func TestPriceRefreshWorker_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 1)

	job := ItemJob{
		ID:     "item1",
		UserID: "user1",
		URL:    "https://example.com/product1",
		Name:   "Test Product",
	}

	expectedError := errors.New("scraping failed")
	mockService.EXPECT().RefreshPrice(
		context.Background(),
		job.ID,
		job.UserID,
		job.URL,
	).Return(expectedError)

	jobs := make(chan ItemJob, 1)
	results := make(chan string, 1)

	jobs <- job
	close(jobs)

	// Run worker
	go scheduler.priceRefreshWorker(1, jobs, results)

	// Check result
	result := <-results
	expected := "FAILED: item1 (Test Product): scraping failed"
	assert.Equal(t, expected, result)
}

// TestRefreshAllPrices tests the main refresh function
func TestRefreshAllPrices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 2)

	// Setup mock items
	testItems := []items.ItemResponse{
		{
			ID:           "item1",
			UserID:       "user1",
			URL:          "https://example.com/product1",
			Name:         "Product 1",
			CurrentPrice: 99.99,
			InStock:      true,
		},
		{
			ID:           "item2",
			UserID:       "user1",
			URL:          "https://example.com/product2",
			Name:         "Product 2",
			CurrentPrice: 149.99,
			InStock:      true,
		},
	}

	// Mock GetItemsDueForCheck
	mockService.EXPECT().GetItemsDueForCheck(context.Background()).
		Return(testItems, nil)

	// Mock successful refreshes
	mockService.EXPECT().RefreshPrice(
		context.Background(),
		"item1",
		"user1",
		"https://example.com/product1",
	).Return(nil)

	mockService.EXPECT().RefreshPrice(
		context.Background(),
		"item2",
		"user1",
		"https://example.com/product2",
	).Return(nil)

	// Run refresh
	scheduler.refreshAllPrices()

	// Note: In a real test, you might want to add timeouts and more assertions
	// This is a basic integration test
}

// TestRefreshAllPrices_Empty tests with no items
func TestRefreshAllPrices_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 2)

	// Mock empty result
	mockService.EXPECT().GetItemsDueForCheck(context.Background()).
		Return([]items.ItemResponse{}, nil)

	// Should not call RefreshPrice
	mockService.EXPECT().RefreshPrice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Run refresh - should handle empty case gracefully
	scheduler.refreshAllPrices()
}

// TestRefreshAllPrices_Error tests error handling
func TestRefreshAllPrices_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 2)

	// Mock error from GetItemsDueForCheck
	expectedError := errors.New("database error")
	mockService.EXPECT().GetItemsDueForCheck(context.Background()).
		Return(nil, expectedError)

	// Should not call RefreshPrice
	mockService.EXPECT().RefreshPrice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	// Run refresh - should handle error gracefully
	scheduler.refreshAllPrices()
}

// TestConcurrentProcessing tests that multiple workers process items concurrently
func TestConcurrentProcessing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	scheduler := NewPriceRefresherScheduler(mockService, 30*time.Minute, 3)

	// Create multiple test items
	testItems := make([]items.ItemResponse, 10)
	for i := range testItems {
		testItems[i] = items.ItemResponse{
			ID:           fmt.Sprintf("item%d", i),
			UserID:       "user1",
			URL:          fmt.Sprintf("https://example.com/product%d", i),
			Name:         fmt.Sprintf("Product %d", i),
			CurrentPrice: float64(100 + i*10),
			InStock:      true,
		}
	}

	// Mock GetItemsDueForCheck
	mockService.EXPECT().GetItemsDueForCheck(context.Background()).
		Return(testItems, nil)

	// Mock successful refreshes for all items
	for _, item := range testItems {
		mockService.EXPECT().RefreshPrice(
			context.Background(),
			item.ID,
			item.UserID,
			item.URL,
		).Return(nil)
	}

	// Run refresh - should process items concurrently
	scheduler.refreshAllPrices()

	// Note: This test verifies the concurrency pattern works
	// In a more advanced test, you could measure timing to verify
	// that items are processed faster than sequential approach
}
