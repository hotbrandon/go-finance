package jobs

import (
	"context"
	"log"
	"sync"
	"time"

	"go-finance/internal/client"

	"github.com/robfig/cron/v3"
)

type JobMeta struct {
	Name string
	Spec string
}

type Scheduler struct {
	cron *cron.Cron
	fc   *client.FinmindClient
	loc  *time.Location

	mu   sync.RWMutex
	jobs map[cron.EntryID]JobMeta
}

type JobStatus struct {
	ID   int       `json:"id"`
	Name string    `json:"name,omitempty"`
	Spec string    `json:"spec,omitempty"`
	Next time.Time `json:"next,omitempty"`
	Prev time.Time `json:"prev,omitempty"`
}

func NewScheduler(fc *client.FinmindClient, loc *time.Location) *Scheduler {
	c := cron.New(
		cron.WithLocation(loc),
		cron.WithSeconds(), // enable seconds field
	)
	return &Scheduler{
		cron: c, fc: fc, loc: loc,
		jobs: make(map[cron.EntryID]JobMeta),
	}
}

// AddJob wraps cron.AddFunc and records metadata
func (s *Scheduler) AddJob(spec, name string, job func()) (cron.EntryID, error) {
	id, err := s.cron.AddFunc(spec, job)
	if err != nil {
		return 0, err
	}
	s.mu.Lock()
	s.jobs[id] = JobMeta{Name: name, Spec: spec}
	s.mu.Unlock()
	return id, nil
}

// ListJobs returns a snapshot of scheduled jobs with next/prev run times
func (s *Scheduler) ListJobs() []JobStatus {
	entries := s.cron.Entries()
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]JobStatus, 0, len(entries))
	for _, entry := range entries {
		meta, exists := s.jobs[entry.ID]
		js := JobStatus{
			ID:   int(entry.ID),
			Next: entry.Next,
			Prev: entry.Prev,
		}
		if exists {
			js.Name = meta.Name
			js.Spec = meta.Spec
		}
		out = append(out, js)
	}
	return out
}

func (s *Scheduler) Start(ctx context.Context) {
	// "0 0 21 * * *" => at 21:00:00 every day (seconds, minutes, hours, day, month, dow)
	_, err := s.AddJob("0 0 21 * * *", "finmind-2337", s.finmindDailyJob(ctx, "2337"))
	if err != nil {
		log.Printf("scheduler add job error: %v", err)
		return
	}

	_, err = s.AddJob("0 1 21 * * *", "finmind-2885", s.finmindDailyJob(ctx, "2885"))
	if err != nil {
		log.Printf("scheduler add job error: %v", err)
		return
	}

	s.cron.Start()

	// stop when ctx is done
	<-ctx.Done()
	log.Println("Stopping the scheduler...")
	ctxStop := s.cron.Stop() // returns a context that waits for running jobs
	<-ctxStop.Done()
}
