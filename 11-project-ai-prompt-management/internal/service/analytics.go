package service

import (
	"context"
	"log"

	"github.com/timurdian/prompt-management/internal/entity"
	"github.com/timurdian/prompt-management/internal/repository"
)

type AnalyticsService interface {
	LogAsync(logEntry *entity.AnalyticsLog)
	StartWorker(ctx context.Context)
}

type analyticsService struct {
	repo repository.PromptRepository
	ch   chan *entity.AnalyticsLog
}

func NewAnalyticsService(repo repository.PromptRepository) AnalyticsService {
	// Buffered channel size 1000
	return &analyticsService{
		repo: repo,
		ch:   make(chan *entity.AnalyticsLog, 1000),
	}
}

func (s *analyticsService) LogAsync(logEntry *entity.AnalyticsLog) {
	select {
	case s.ch <- logEntry:
		// Berhasil masuk antrean
	default:
		// Antrean penuh, buang atau log error (mencegah bottleneck)
		log.Println("[WARNING] Analytics buffer channel is full, dropping log entry")
	}
}

func (s *analyticsService) StartWorker(ctx context.Context) {
	go func() {
		log.Println("Analytics worker daemon started successfully")
		for {
			select {
			case logEntry := <-s.ch:
				err := s.repo.LogAnalytics(context.Background(), logEntry)
				if err != nil {
					log.Printf("[ERROR] Failed to save analytics log: %v", err)
				}
			case <-ctx.Done():
				log.Println("Analytics worker daemon shutting down...")
				// Habiskan antrean sisa sebelum exit
				close(s.ch)
				for logEntry := range s.ch {
					_ = s.repo.LogAnalytics(context.Background(), logEntry)
				}
				return
			}
		}
	}()
}
