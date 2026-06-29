package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/timurdian/notification-service/internal/queue"
)

type SchedulerPoller interface {
	Start(ctx context.Context)
}

type schedulerPoller struct {
	queueMgr queue.QueueManager
}

func NewSchedulerPoller(queueMgr queue.QueueManager) SchedulerPoller {
	return &schedulerPoller{queueMgr: queueMgr}
}

func (p *schedulerPoller) Start(ctx context.Context) {
	log.Println("[SchedulerPoller] Starting 1-second interval scheduled/retry poller loop...")
	go p.pollLoop(ctx)
}

func (p *schedulerPoller) pollLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[SchedulerPoller] Stopped poller loop due to context cancellation")
			return
		case <-ticker.C:
			// Panggil atomic Lua script untuk memindahkan scheduled task ke instant queue
			count, err := p.queueMgr.MoveScheduledToReady(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					continue
				}
				log.Printf("[SchedulerPoller] Error polling scheduled tasks: %v", err)
				continue
			}

			if count > 0 {
				log.Printf("[SchedulerPoller] Promoted %d scheduled/retry tasks from Sorted Set to instant Queue", count)
			}
		}
	}
}
