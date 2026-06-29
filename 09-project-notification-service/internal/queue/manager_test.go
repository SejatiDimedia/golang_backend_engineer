package queue

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedisQueueManager_EnqueueAndDequeue(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	qm := NewQueueManager(rdb)
	task := &Task{
		NotificationID: 42,
		Type:           "email",
		Target:         "test@email.com",
		Content:        "Hello world",
		Attempt:        0,
	}

	ctx := context.Background()

	// 1. Enqueue
	err = qm.Enqueue(ctx, task)
	if err != nil {
		t.Fatalf("failed to enqueue: %v", err)
	}

	// 2. Dequeue
	dequeued, err := qm.Dequeue(ctx)
	if err != nil {
		t.Fatalf("failed to dequeue: %v", err)
	}

	if dequeued.NotificationID != task.NotificationID {
		t.Errorf("expected ID %d, got %d", task.NotificationID, dequeued.NotificationID)
	}
	if dequeued.Target != task.Target {
		t.Errorf("expected Target %s, got %s", task.Target, dequeued.Target)
	}
}

func TestRedisQueueManager_MoveScheduledToReady(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	qm := NewQueueManager(rdb)
	task := &Task{
		NotificationID: 99,
		Type:           "webhook",
		Target:         "https://webhook.site",
		Content:        "payload",
		Attempt:        1,
	}

	ctx := context.Background()

	// 1. Enqueue Scheduled 5 detik di masa depan
	sendAt := time.Now().Add(5 * time.Second)
	err = qm.EnqueueScheduled(ctx, task, sendAt)
	if err != nil {
		t.Fatalf("failed to enqueue scheduled: %v", err)
	}

	// 2. Coba trigger MoveScheduledToReady sekarang -> harusnya belum jatuh tempo (0 dipromosikan)
	count, err := qm.MoveScheduledToReady(ctx)
	if err != nil {
		t.Fatalf("failed to move scheduled tasks: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 promoted tasks, got %d", count)
	}

	// 3. Fast-forward waktu di miniredis melewati sendAt
	mr.FastForward(6 * time.Second)

	// 4. Trigger MoveScheduledToReady kembali -> harusnya 1 dipromosikan ke list
	count, err = qm.MoveScheduledToReady(ctx)
	if err != nil {
		t.Fatalf("failed to move scheduled tasks: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 promoted task, got %d", count)
	}

	// 5. Dequeue tugas yang dipromosikan untuk verifikasi
	dequeued, err := qm.Dequeue(ctx)
	if err != nil {
		t.Fatalf("failed to dequeue: %v", err)
	}
	if dequeued.NotificationID != task.NotificationID {
		t.Errorf("expected ID %d, got %d", task.NotificationID, dequeued.NotificationID)
	}
}
