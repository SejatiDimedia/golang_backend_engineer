package worker

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/timurdian/notification-service/internal/provider"
	"github.com/timurdian/notification-service/internal/queue"
	"github.com/timurdian/notification-service/internal/service"
)

type WorkerPool interface {
	Start(ctx context.Context)
}

type workerPool struct {
	queueMgr     queue.QueueManager
	notifService service.NotificationService
	provider     provider.NotificationProvider
	concurrency  int
}

func NewWorkerPool(
	queueMgr queue.QueueManager,
	notifService service.NotificationService,
	provider provider.NotificationProvider,
	concurrency int,
) WorkerPool {
	return &workerPool{
		queueMgr:     queueMgr,
		notifService: notifService,
		provider:     provider,
		concurrency:  concurrency,
	}
}

func (p *workerPool) Start(ctx context.Context) {
	log.Printf("[WorkerPool] Starting %d background workers...", p.concurrency)
	for i := 1; i <= p.concurrency; i++ {
		go p.worker(ctx, i)
	}
}

func (p *workerPool) worker(ctx context.Context, workerID int) {
	log.Printf("[Worker %d] Started and waiting for tasks...", workerID)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[Worker %d] Stopped due to context cancellation", workerID)
			return
		default:
		}

		// Dequeue blocking wait
		task, err := p.queueMgr.Dequeue(ctx)
		if err != nil {
			// Jika context cancel, ignore error brpop
			if errors.Is(err, context.Canceled) {
				continue
			}
			log.Printf("[Worker %d] Dequeue error: %v", workerID, err)
			time.Sleep(1 * time.Second) // backoff singkat agar tidak loop super-cepat jika Redis drop
			continue
		}

		p.processTask(ctx, workerID, task)
	}
}

func (p *workerPool) processTask(ctx context.Context, workerID int, task *queue.Task) {
	log.Printf("[Worker %d] Processing notification ID %d (Attempt %d)", workerID, task.NotificationID, task.Attempt+1)

	// 1. Cek status notifikasi terupdate dari DB
	notif, err := p.notifService.GetStatus(ctx, task.NotificationID)
	if err != nil {
		log.Printf("[Worker %d] Failed to fetch status for ID %d: %v", workerID, task.NotificationID, err)
		return
	}

	// Mencegah pengiriman ganda jika status sudah SENT atau FAILED (DLQ)
	if notif.Status == "SENT" || notif.Status == "FAILED" {
		log.Printf("[Worker %d] Skipping task ID %d: status is already %s", workerID, task.NotificationID, notif.Status)
		return
	}

	// Update status ke PROCESSING di DB
	_ = p.notifService.UpdateStatus(ctx, notif.ID, "PROCESSING", task.Attempt, "")

	// 2. Kirim Notifikasi via Provider
	err = p.provider.Send(ctx, task.Type, task.Target, task.Content)
	if err == nil {
		// SUKSES
		log.Printf("[Worker %d] SUCCESS: Notification ID %d successfully sent", workerID, task.NotificationID)
		_ = p.notifService.UpdateStatus(ctx, notif.ID, "SENT", task.Attempt+1, "")
		return
	}

	// GAGAL - Lakukan Audit Log Gagal
	log.Printf("[Worker %d] FAILED: Notification ID %d failed: %v", workerID, task.NotificationID, err)
	task.Attempt++

	// Cek limit maximum retries
	if task.Attempt >= notif.MaxRetries {
		log.Printf("[Worker %d] DLQ REACHED: Notification ID %d reached maximum retries (%d). Marking as FAILED", workerID, task.NotificationID, notif.MaxRetries)
		_ = p.notifService.UpdateStatus(ctx, notif.ID, "FAILED", task.Attempt, err.Error())
		return
	}

	// 3. Hitung Exponential Backoff Delay
	// Delay = (2^attempt) * 2 detik -> Attempt 1: 4s, Attempt 2: 8s, Attempt 3: 16s, Attempt 4: 32s
	delaySeconds := int(math.Pow(2, float64(task.Attempt))) * 2
	sendAt := time.Now().Add(time.Duration(delaySeconds) * time.Second)

	log.Printf("[Worker %d] RETRY QUEUED: Notification ID %d will retry in %d seconds (at %s)", workerID, task.NotificationID, delaySeconds, sendAt.Format("15:04:05"))

	// Update status di DB relasional kembali ke PENDING dengan riwayat log
	_ = p.notifService.UpdateStatus(ctx, notif.ID, "PENDING", task.Attempt, err.Error())

	// Enqueue back into Redis Sorted Set (ZSET)
	err = p.queueMgr.EnqueueScheduled(ctx, task, sendAt)
	if err != nil {
		log.Printf("[Worker %d] Failed to requeue scheduled task for ID %d: %v", workerID, task.NotificationID, err)
	}
}
