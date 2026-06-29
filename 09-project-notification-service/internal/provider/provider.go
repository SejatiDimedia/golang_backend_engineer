package provider

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type NotificationProvider interface {
	Send(ctx context.Context, notifType, target, content string) error
}

type mockNotificationProvider struct {
	failureRate float64
	randGen     *rand.Rand
}

func NewMockNotificationProvider(failureRate float64) NotificationProvider {
	return &mockNotificationProvider{
		failureRate: failureRate,
		randGen:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *mockNotificationProvider) Send(ctx context.Context, notifType, target, content string) error {
	// Simulasi sedikit latency jaringan (50 - 150 ms)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Duration(50+p.randGen.Intn(100)) * time.Millisecond):
	}

	// Uji coba kegagalan jaringan acak
	if p.randGen.Float64() < p.failureRate {
		return fmt.Errorf("provider connection failure: server returned 503 service unavailable")
	}

	log.Printf("[Provider] SUCCESS: Sent %s notification to '%s'. Content preview: %.30s...", notifType, target, content)
	return nil
}
