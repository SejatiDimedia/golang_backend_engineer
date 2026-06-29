package service

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/timurdian/notification-service/internal/entity"
	"github.com/timurdian/notification-service/internal/queue"
)

// Mock NotificationRepository
type mockNotificationRepository struct {
	notifs map[uint]*entity.Notification
	logs   []*entity.NotificationLog
	nextID uint
}

func (m *mockNotificationRepository) Create(ctx context.Context, notif *entity.Notification) error {
	m.nextID++
	notif.ID = m.nextID
	m.notifs[notif.ID] = notif
	return nil
}

func (m *mockNotificationRepository) Update(ctx context.Context, notif *entity.Notification) error {
	m.notifs[notif.ID] = notif
	return nil
}

func (m *mockNotificationRepository) GetByID(ctx context.Context, id uint) (*entity.Notification, error) {
	n, ok := m.notifs[id]
	if !ok {
		return nil, nil
	}
	return n, nil
}

func (m *mockNotificationRepository) CreateLog(ctx context.Context, log *entity.NotificationLog) error {
	m.logs = append(m.logs, log)
	return nil
}

// Mock QueueManager
type mockQueueManager struct {
	instant   []*queue.Task
	scheduled map[int64][]*queue.Task
}

func (m *mockQueueManager) Enqueue(ctx context.Context, task *queue.Task) error {
	m.instant = append(m.instant, task)
	return nil
}

func (m *mockQueueManager) Dequeue(ctx context.Context) (*queue.Task, error) {
	if len(m.instant) == 0 {
		return nil, nil
	}
	t := m.instant[0]
	m.instant = m.instant[1:]
	return t, nil
}

func (m *mockQueueManager) EnqueueScheduled(ctx context.Context, task *queue.Task, sendAt time.Time) error {
	sec := sendAt.Unix()
	m.scheduled[sec] = append(m.scheduled[sec], task)
	return nil
}

func (m *mockQueueManager) MoveScheduledToReady(ctx context.Context) (int64, error) {
	return 0, nil
}

func TestNotificationService_Create_Instant(t *testing.T) {
	repo := &mockNotificationRepository{notifs: make(map[uint]*entity.Notification)}
	qm := &mockQueueManager{scheduled: make(map[int64][]*queue.Task)}
	svc := NewNotificationService(repo, qm)

	ctx := context.Background()
	notif, err := svc.Create(ctx, "email", "target@email.com", "my text", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if notif.Status != "PENDING" {
		t.Errorf("expected PENDING status, got %s", notif.Status)
	}
	if len(qm.instant) != 1 {
		t.Errorf("expected 1 instant task queued, got %d", len(qm.instant))
	}
	if qm.instant[0].NotificationID != notif.ID {
		t.Errorf("expected queued ID %d, got %d", notif.ID, qm.instant[0].NotificationID)
	}
}

func TestNotificationService_Create_Scheduled(t *testing.T) {
	repo := &mockNotificationRepository{notifs: make(map[uint]*entity.Notification)}
	qm := &mockQueueManager{scheduled: make(map[int64][]*queue.Task)}
	svc := NewNotificationService(repo, qm)

	ctx := context.Background()
	sendAt := time.Now().Add(1 * time.Hour)
	notif, err := svc.Create(ctx, "webhook", "https://api.myweb.com/callback", "hello json", &sendAt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(qm.instant) != 0 {
		t.Errorf("expected 0 instant tasks, got %d", len(qm.instant))
	}
	if len(qm.scheduled) != 1 {
		t.Errorf("expected 1 scheduled task bucket, got %d", len(qm.scheduled))
	}
	
	// Verifikasi timestamp key
	tasks := qm.scheduled[sendAt.Unix()]
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task in scheduled timestamp bucket, got %d", len(tasks))
	}
	if tasks[0].NotificationID != notif.ID {
		t.Errorf("expected ID %d, got %d", notif.ID, tasks[0].NotificationID)
	}
}

func TestExponentialBackoffCalculation(t *testing.T) {
	// Verifikasi rumus backoff matematika 2^attempt * 2
	tests := []struct {
		attempt       int
		expectedDelay int
	}{
		{attempt: 1, expectedDelay: 4},  // 2^1 * 2 = 4
		{attempt: 2, expectedDelay: 8},  // 2^2 * 2 = 8
		{attempt: 3, expectedDelay: 16}, // 2^3 * 2 = 16
		{attempt: 4, expectedDelay: 32}, // 2^4 * 2 = 32
	}

	for _, tc := range tests {
		delaySeconds := int(math.Pow(2, float64(tc.attempt))) * 2
		if delaySeconds != tc.expectedDelay {
			t.Errorf("for attempt %d: expected delay %d, got %d", tc.attempt, tc.expectedDelay, delaySeconds)
		}
	}
}
