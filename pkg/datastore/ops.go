package datastore

import "time"

// Store 定义对 Redis 的常用操作接口
type Store interface {
	// 字符串
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Incr(key string) (int64, error)
	Decr(key string) (int64, error)
	Del(key string) error
	Exists(key string) (bool, error)
	Expire(key string, expiration time.Duration) error

	// 哈希
	HSet(key string, values ...interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HDel(key string, fields ...string) error

	// 列表
	LPush(key string, values ...interface{}) error
	RPop(key string) (string, error)

	// 集合
	SAdd(key string, members ...interface{}) error
	SMembers(key string) ([]string, error)
	SRem(key string, members ...interface{}) error

	// 分布式锁
	TryLock(key, value string, expiration time.Duration) (bool, error)
	Unlock(key string) error

	// Lua 脚本
	Eval(script string, keys []string, args ...interface{}) (interface{}, error)
}

// ============ 实现方法 ============

// --- 字符串操作 ---
func (r *redisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *redisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *redisClient) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *redisClient) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

func (r *redisClient) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *redisClient) Exists(key string) (bool, error) {
	res, err := r.client.Exists(r.ctx, key).Result()
	return res > 0, err
}

func (r *redisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// --- 哈希操作 ---
func (r *redisClient) HSet(key string, values ...interface{}) error {
	return r.client.HSet(r.ctx, key, values...).Err()
}

func (r *redisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

func (r *redisClient) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

func (r *redisClient) HDel(key string, fields ...string) error {
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// --- 列表操作 ---
func (r *redisClient) LPush(key string, values ...interface{}) error {
	return r.client.LPush(r.ctx, key, values...).Err()
}

func (r *redisClient) RPop(key string) (string, error) {
	return r.client.RPop(r.ctx, key).Result()
}

// --- 集合操作 ---
func (r *redisClient) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

func (r *redisClient) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

func (r *redisClient) SRem(key string, members ...interface{}) error {
	return r.client.SRem(r.ctx, key, members...).Err()
}

// --- 分布式锁 ---
func (r *redisClient) TryLock(key, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(r.ctx, key, value, expiration).Result()
}

func (r *redisClient) Unlock(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// --- Lua 脚本执行 ---
func (r *redisClient) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(r.ctx, script, keys, args...).Result()
}
