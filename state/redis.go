package state

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/amirrezaask/plumber"
	"github.com/go-redis/redis/v8"
)

type redisState struct {
	ctx  context.Context
	conn *redis.Client
	ttl  time.Duration
}

func NewRedis(ctx context.Context, host, port, user, password string, database int) (plumber.State, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		DB:       database,
		Username: user,
		Password: password,
	})
	res := c.Ping(ctx)
	if err := res.Err(); err != nil {
		return nil, err
	}
	return &redisState{
		conn: c,
		ctx:  ctx,
	}, nil
}

func (r *redisState) Flush() error {
	return nil
}

func (r *redisState) All() (map[string]interface{}, error) {
	res := r.conn.Keys(r.ctx, "*")
	if err := res.Err(); err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	for _, k := range res.Val() {
		val := r.conn.Get(r.ctx, k)
		if err := val.Err(); err != nil {
			return nil, err
		}
		m[k] = val.Val()
	}
	return m, nil
}
func (r *redisState) GetInt(key string) (int, error) {
	v, err := r.Get(key)
	if err != nil {
		return -1, err
	}
	if v == nil {
		return 0, nil
	}
	i, err := strconv.Atoi(v.(string))
	if err != nil {
		return -1, err
	}
	return i, nil
}
func (r *redisState) Get(key string) (interface{}, error) {
	res := r.conn.Get(r.ctx, key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return res.Val(), nil
}

func (r *redisState) GetBytes(key string) ([]byte, error) {
	res := r.conn.Get(r.ctx, key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return res.Bytes()
}
func (r *redisState) Set(key string, value interface{}) error {
	statusCmd := r.conn.Set(r.ctx, key, value, r.ttl)
	if err := statusCmd.Err(); err != nil {
		return err
	}
	return nil
}
