package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrLockAcquireFailed = errors.New("failed to acquire distributed lock")

type LockManager interface {
	AcquireLock(ctx context.Context, key string, ttl time.Duration, retries int, backoff time.Duration) (string, error)
	ReleaseLock(ctx context.Context, key string, token string) error
}

type redisLockManager struct {
	client *redis.Client
}

func NewRedisLockManager(client *redis.Client) LockManager {
	return &redisLockManager{client: client}
}

func (m *redisLockManager) AcquireLock(ctx context.Context, key string, ttl time.Duration, retries int, backoff time.Duration) (string, error) {
	// Generate random token unik untuk lock session ini
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	for i := 0; i <= retries; i++ {
		// SET key token NX PX ttl
		ok, err := m.client.SetNX(ctx, key, token, ttl).Result()
		if err == nil && ok {
			return token, nil
		}
		
		if i < retries {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}
	}
	return "", ErrLockAcquireFailed
}

func (m *redisLockManager) ReleaseLock(ctx context.Context, key string, token string) error {
	// Script Lua atomik untuk mencocokkan token unik dan menghapus kunci
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	res, err := m.client.Eval(ctx, luaScript, []string{key}, token).Result()
	if err != nil {
		return err
	}
	
	if val, ok := res.(int64); ok && val == 1 {
		return nil
	}
	return errors.New("lock release failed: token mismatch or lock expired")
}
