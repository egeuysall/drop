package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/egeuysall/drop/internal/modules/items"
)

// ItemJob represents a single item to be processed by workers
type ItemJob struct {
	ID        string
	UserID    string
	URL       string
	Name      string
}

type PriceRefresherScheduler struct {
    itemsService items.Service
    interval     time.Duration
    workerCount  int
    stopChan     chan bool
}

func NewPriceRefresherScheduler(itemsService items.Service, interval time.Duration, workerCount int) *PriceRefresherScheduler {
    return &PriceRefresherScheduler{
        itemsService: itemsService,
        interval:     interval,
        workerCount:  workerCount,
        stopChan:     make(chan bool),
    }
}


func (s *PriceRefresherScheduler) Start() {
    s.refreshAllPrices();

    ticker := time.NewTicker(s.interval)

    go func() {
        for {
            select {
                case <- ticker.C:
                    s.refreshAllPrices();
                case <- s.stopChan:
                    ticker.Stop();
                    return;
            }
        }
    }()
}

func (s *PriceRefresherScheduler) Stop() {
    s.stopChan <- true;
}

func (s *PriceRefresherScheduler) refreshAllPrices() {
    ctx := context.Background()
    items, err := s.itemsService.GetItemsDueForCheck(ctx)

    if err != nil {
        log.Printf("Error while refreshing prices: %s", err.Error())
        return
    }

    if len(items) == 0 {
        log.Printf("No items due for price refresh")
        return
    }

    log.Printf("Starting concurrent refresh of %d items with %d workers", len(items), s.workerCount)

    // Create channels for work distribution
    jobs := make(chan ItemJob, len(items))
    results := make(chan string, len(items))

    // Start worker pool
    for w := 1; w <= s.workerCount; w++ {
        go s.priceRefreshWorker(w, jobs, results)
    }

    // Convert items to jobs and send to workers
    for _, item := range items {
        jobs <- ItemJob{
            ID:     item.ID,
            UserID: item.UserID,
            URL:    item.URL,
            Name:   item.Name,
        }
    }
    close(jobs)

    successCount := 0
    failCount := 0

    for range items {
        result := <-results
        if strings.HasPrefix(result, "SUCCESS:") {
            successCount++
            log.Println(result)
        } else {
            failCount++
            log.Println(result)
        }
    }

    log.Printf("Price refresh complete: %d succeeded, %d failed out of %d total", successCount, failCount, len(items))
}

// priceRefreshWorker processes individual refresh jobs
// Each worker runs independently with no shared state
func (s *PriceRefresherScheduler) priceRefreshWorker(workerID int, jobs <-chan ItemJob, results chan<- string) {
    for job := range jobs {
        // Debug log showing which worker is processing which item
        log.Printf("Worker %d processing item: %s", workerID, job.Name)

        err := s.itemsService.RefreshPrice(context.Background(), job.ID, job.UserID, job.URL)

        if err != nil {
            results <- fmt.Sprintf("FAILED: %s (%s): %v", job.ID, job.Name, err)
        } else {
            results <- fmt.Sprintf("SUCCESS: %s", job.Name)
        }
    }
}
