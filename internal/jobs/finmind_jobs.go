package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"go-finance/internal/models"
	"log"
	"time"
)

// finmindDailyJob returns the job function to schedule. It captures the parent context
// so cancellation from main will propagate into the running job.
func (s *Scheduler) finmindDailyJob(ctx context.Context, stock_id string) func() {
	return func() {
		// run the job with a per-job timeout
		jobCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		// example: fetch for stock 2337 from today (adjust as needed)
		// yesterday := time.Now().In(s.loc).AddDate(0, 0, -1).Format("2006-01-02")
		// _, err := s.fc.GetTaiwanStockPrice(jobCtx, "TaiwanStockPrice", "2330", yesterday)
		today := time.Now().In(s.loc).Format("2006-01-02")
		res, err := s.fc.GetTaiwanStockPrice(jobCtx, "TaiwanStockPrice", stock_id, today)
		if err != nil {
			// distinguish cancellation/timeout from other errors
			if errors.Is(err, context.Canceled) || jobCtx.Err() == context.Canceled {
				log.Printf("nightly job: fetch cancelled for %s on %s", stock_id, today)
				return
			}
			if errors.Is(err, context.DeadlineExceeded) || jobCtx.Err() == context.DeadlineExceeded {
				log.Printf("nightly job: fetch deadline exceeded for %s on %s", stock_id, today)
				return
			}
			log.Printf("nightly job: fetch error: %v", err)
			return
		}

		if len(res.Data) == 0 {
			log.Printf("nightly job: no data for TaiwanStockPrice on %s", today)
			return
		}
		var prices []models.StockPrice
		if err := json.Unmarshal(res.Data, &prices); err != nil {
			log.Printf("nightly job: invalid data format for stock_id %s: %s", stock_id, err.Error())
			return
		}

		log.Printf("nightly job: fetched TaiwanStockPrice for %s: %v", prices[0].Date, prices[0].Close)
	}
}
