package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	Redis  *redis.Client
	Prefix string
}

func (r *Redis) Push(key string, value any) error {
	err := r.Redis.RPush(r.getKey(key), value).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Retrieve(key string) []string {
	data, err := r.Redis.LRange(r.getKey(key), 0, -1).Result()
	if err != nil {
		fmt.Println("Error retrieving data:", err)
	}

	return data
}

func (r *Redis) Remove(key string, idx int) {
	data, err := r.Redis.LIndex(r.getKey(key), int64(idx)).Result()
	if err != nil {
		fmt.Println("Error get index data:", err)
	}

	err = r.Redis.LRem(r.getKey(key), 1, data).Err()

	if err != nil {
		fmt.Println("Error remove data:", err)
	}
}

func (r *Redis) Pop(key string) []string {
	result, err := r.Redis.BLPop(0*time.Second, r.getKey(key)).Result()

	if err != nil {
		fmt.Println(err.Error())
	}

	return result
}

// Get Retrieve an item from the cache by key.
func (r *Redis) Get(key string, defaults any) any {
	val, err := r.Redis.Get(r.getKey(key)).Result()
	if err != nil {
		switch s := defaults.(type) {
		case func() any:
			return s()
		default:
			return defaults
		}
	}

	return val
}

// Has Determine if an item exists in the cache.
func (r *Redis) Has(key string) bool {
	value, err := r.Redis.Exists(r.getKey(key)).Result()

	return !(err != nil || value == 0)
}

// Put Store an item in the cache for a given number of seconds.
func (r *Redis) Set(key string, value any, expiration time.Duration) error {
	switch value.(type) {
	case map[string]any, map[string]string, struct{}, any:
		// Serialize the data to JSON
		jsonData, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return r.Redis.Set(r.getKey(key), jsonData, expiration).Err()
	default:
		return r.Redis.Set(r.getKey(key), value, expiration).Err()
	}
}

// Pull Retrieve an item from the cache and delete it.
func (r *Redis) Pull(key string, defaults any) any {
	val, err := r.Redis.Get(r.getKey(key)).Result()
	r.Redis.Del(r.getKey(key))

	if err != nil {
		return defaults
	}

	return val
}

// Add Store an item in the cache if the key does not exist.
func (r *Redis) Add(key string, value any, seconds time.Duration) bool {
	val, err := r.Redis.SetNX(r.getKey(key), value, seconds).Result()
	if err != nil {
		return false
	}

	return val
}

// Remember Get an item from the cache, or execute the given Closure and store the result.
func (r *Redis) Remember(key string, ttl time.Duration, callback func() any) (any, error) {
	val := r.Get(key, nil)

	if val != nil {
		return val, nil
	}

	val = callback()

	if err := r.Set(key, val, ttl); err != nil {
		return nil, err
	}

	return val, nil
}

// RememberForever Get an item from the cache, or execute the given Closure and store the result forever.
func (r *Redis) RememberForever(key string, callback func() any) (any, error) {
	val := r.Get(key, nil)

	if val != nil {
		return val, nil
	}

	val = callback()

	if err := r.Set(key, val, 0); err != nil {
		return nil, err
	}

	return val, nil
}

// Forever Store an item in the cache indefinitely.
func (r *Redis) Forever(key string, value any) bool {
	if err := r.Set(key, value, 0); err != nil {
		return false
	}

	return true
}

// Forget Remove an item from the cache.
func (r *Redis) Forget(key string) bool {
	_, err := r.Redis.Del(r.getKey(key)).Result()

	return err == nil
}

// Flush Remove all items from the cache.
func (r *Redis) Flush() bool {
	res, err := r.Redis.FlushAll().Result()

	if err != nil || res != "OK" {
		return false
	}

	return true
}

func (r *Redis) getKey(key string) string {
	return r.Prefix + "_" + key
}
