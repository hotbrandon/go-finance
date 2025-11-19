package jobs

import (
	"context"
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

		// watcher: print when the job is cancelled or deadline exceeded
		go func() {
			<-jobCtx.Done()
			switch jobCtx.Err() {
			case context.Canceled:
				log.Printf("nightly job: cancelled")
			case context.DeadlineExceeded:
				log.Printf("nightly job: deadline exceeded")
			}
		}()

		// example: fetch for stock 2337 from today (adjust as needed)
		// yesterday := time.Now().In(s.loc).AddDate(0, 0, -1).Format("2006-01-02")
		// _, err := s.fc.GetTaiwanStockPrice(jobCtx, "TaiwanStockPrice", "2330", yesterday)
		today := time.Now().In(s.loc).Format("2006-01-02")
		res, err := s.fc.GetTaiwanStockPrice(jobCtx, "TaiwanStockPrice", stock_id, today)
		if err != nil {
			log.Printf("nightly job: fetch error: %v", err)
			return
		}

		if len(res.Data) == 0 {
			log.Printf("nightly job: no data for TaiwanStockPrice on %s", today)
			return
		}
		log.Printf("nightly job: fetched TaiwanStockPrice for %s", today)
	}
}
