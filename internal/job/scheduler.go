package job

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	fetcher  *PriceFetcher
	interval time.Duration
	stopCh   chan struct{}
}

func NewScheduler(fetcher *PriceFetcher, interval time.Duration) *Scheduler {
	return &Scheduler{
		fetcher:  fetcher,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("Scheduler started: fetching prices every %v", s.interval)

	// 起動直後に1回実行
	s.fetcher.FetchAndSaveAll(ctx)

	for {
		select {
		case <-ticker.C:
			s.fetcher.FetchAndSaveAll(ctx)
		case <-s.stopCh:
			log.Println("Scheduler stopped")
			return
		case <-ctx.Done():
			log.Println("Scheduler stopped by context")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}
