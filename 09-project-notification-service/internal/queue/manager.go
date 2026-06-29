package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Task struct {
	NotificationID uint   `json:"notification_id"`
	Type           string `json:"type"` // email, webhook, push
	Target         string `json:"target"`
	Content        string `json:"content"`
	Attempt        int    `json:"attempt"`
}

type QueueManager interface {
	Enqueue(ctx context.Context, task *Task) error
	Dequeue(ctx context.Context) (*Task, error)
	EnqueueScheduled(ctx context.Context, task *Task, sendAt time.Time) error
	MoveScheduledToReady(ctx context.Context) (int64, error)
}

type redisQueueManager struct {
	rdb     *redis.Client
	listKey string // "queue:notifications"
	zsetKey string // "queue:scheduled"
}

func NewQueueManager(rdb *redis.Client) QueueManager {
	return &redisQueueManager{
		rdb:     rdb,
		listKey: "queue:notifications",
		zsetKey: "queue:scheduled",
	}
}

func (q *redisQueueManager) Enqueue(ctx context.Context, task *Task) error {
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return q.rdb.LPush(ctx, q.listKey, taskBytes).Err()
}

func (q *redisQueueManager) Dequeue(ctx context.Context) (*Task, error) {
	// BRPop is a blocking pop from the right (tail) of the list
	// Timeout 0 means wait indefinitely
	res, err := q.rdb.BRPop(ctx, 0, q.listKey).Result()
	if err != nil {
		return nil, err
	}

	// res[0] is the key name ("queue:notifications")
	// res[1] is the value (task bytes)
	if len(res) < 2 {
		return nil, fmt.Errorf("unexpected brpop response length")
	}

	var task Task
	err = json.Unmarshal([]byte(res[1]), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (q *redisQueueManager) EnqueueScheduled(ctx context.Context, task *Task, sendAt time.Time) error {
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return q.rdb.ZAdd(ctx, q.zsetKey, redis.Z{
		Score:  float64(sendAt.Unix()),
		Member: taskBytes,
	}).Err()
}

func (q *redisQueueManager) MoveScheduledToReady(ctx context.Context) (int64, error) {
	// Atomic Lua Script to fetch scheduled tasks, remove them from ZSET, and push to List
	luaScript := `
		local tasks = redis.call("zrangebyscore", KEYS[1], "-inf", ARGV[1])
		if #tasks > 0 then
			redis.call("zremrangebyscore", KEYS[1], "-inf", ARGV[1])
			for i = 1, #tasks do
				redis.call("lpush", KEYS[2], tasks[i])
			end
		end
		return #tasks
	`

	now := time.Now().Unix()
	res, err := q.rdb.Eval(ctx, luaScript, []string{q.zsetKey, q.listKey}, now).Result()
	if err != nil {
		return 0, err
	}

	count, ok := res.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected lua response format")
	}

	return count, nil
}
