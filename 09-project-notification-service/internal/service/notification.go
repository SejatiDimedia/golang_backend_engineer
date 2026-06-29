package service

import (
	"context"
	"errors"
	"time"

	"github.com/timurdian/notification-service/internal/entity"
	"github.com/timurdian/notification-service/internal/queue"
	"github.com/timurdian/notification-service/internal/repository"
)

var (
	ErrInvalidNotificationType = errors.New("invalid notification type, allowed: email, webhook, push")
	ErrInvalidTarget           = errors.New("target parameter is required")
	ErrInvalidContent          = errors.New("content parameter is required")
	ErrNotificationNotFound    = errors.New("notification not found")
)

type NotificationService interface {
	Create(ctx context.Context, notifType, target, content string, sendAt *time.Time) (*entity.Notification, error)
	GetStatus(ctx context.Context, id uint) (*entity.Notification, error)
	UpdateStatus(ctx context.Context, id uint, status string, attempt int, errMsg string) error
}

type notificationService struct {
	notifRepo repository.NotificationRepository
	queueMgr  queue.QueueManager
}

func NewNotificationService(notifRepo repository.NotificationRepository, queueMgr queue.QueueManager) NotificationService {
	return &notificationService{
		notifRepo: notifRepo,
		queueMgr:  queueMgr,
	}
}

func (s *notificationService) Create(ctx context.Context, notifType, target, content string, sendAt *time.Time) (*entity.Notification, error) {
	// 1. Validasi Input
	if notifType != "email" && notifType != "webhook" && notifType != "push" {
		return nil, ErrInvalidNotificationType
	}
	if target == "" {
		return nil, ErrInvalidTarget
	}
	if content == "" {
		return nil, ErrInvalidContent
	}

	scheduledTime := time.Now()
	if sendAt != nil {
		scheduledTime = *sendAt
	}

	// 2. Simpan metadata awal ke PostgreSQL (status = PENDING)
	notif := &entity.Notification{
		Type:         notifType,
		Target:       target,
		Content:      content,
		Status:       "PENDING",
		MaxRetries:   5,
		AttemptCount: 0,
		SendAt:       scheduledTime,
	}

	if err := s.notifRepo.Create(ctx, notif); err != nil {
		return nil, err
	}

	// 3. Masukkan ke Antrean (Redis)
	task := &queue.Task{
		NotificationID: notif.ID,
		Type:           notif.Type,
		Target:         notif.Target,
		Content:        notif.Content,
		Attempt:        0,
	}

	// Jika scheduledTime adalah masa depan, masukkan ke Scheduled Queue ZSET
	if scheduledTime.After(time.Now()) {
		if err := s.queueMgr.EnqueueScheduled(ctx, task, scheduledTime); err != nil {
			return nil, err
		}
	} else {
		// Jika instan, masukkan ke instant List Queue
		if err := s.queueMgr.Enqueue(ctx, task); err != nil {
			return nil, err
		}
	}

	return notif, nil
}

func (s *notificationService) GetStatus(ctx context.Context, id uint) (*entity.Notification, error) {
	notif, err := s.notifRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if notif == nil {
		return nil, ErrNotificationNotFound
	}
	return notif, nil
}

func (s *notificationService) UpdateStatus(ctx context.Context, id uint, status string, attempt int, errMsg string) error {
	notif, err := s.notifRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notif == nil {
		return ErrNotificationNotFound
	}

	// Update status
	notif.Status = status
	notif.AttemptCount = attempt

	if err := s.notifRepo.Update(ctx, notif); err != nil {
		return err
	}

	// Catat audit log percobaan
	logEntry := &entity.NotificationLog{
		NotificationID: notif.ID,
		Attempt:        attempt,
		Status:         status,
		ErrorMessage:   errMsg,
	}

	return s.notifRepo.CreateLog(ctx, logEntry)
}
