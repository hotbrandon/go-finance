package jobs

import (
	"context"
	"log"
	"time"

	"go-finance/internal/client"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
	fc   *client.FinmindClient
	loc  *time.Location
}

func NewScheduler(fc *client.FinmindClient, loc *time.Location) *Scheduler {
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithSeconds(), // enable seconds field
	)
	return &Scheduler{cron: c, fc: fc, loc: loc}
}

func (s *Scheduler) Start(ctx context.Context) {
	// "0 0 21 * * *" => at 21:00:00 every day (seconds, minutes, hours, day, month, dow)
	_, err := s.cron.AddFunc("0 0 21 * * *", func() {
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
		res, err := s.fc.GetTaiwanStockPrice(jobCtx, "TaiwanStockPrice", "2330", today)
		if err != nil {
			log.Printf("nightly job: fetch error: %v", err)
			return
		}

		if len(res.Data) == 0 {
			log.Printf("nightly job: no data for TaiwanStockPrice on %s", today)
			return
		}
		log.Printf("nightly job: fetched TaiwanStockPrice for %s", today)
	})
	if err != nil {
		log.Printf("scheduler add func error: %v", err)
		return
	}

	s.cron.Start()

	// stop when ctx is done
	<-ctx.Done()
	log.Println("Stopping the scheduler...")
	ctxStop := s.cron.Stop() // returns a context that waits for running jobs
	<-ctxStop.Done()
}
