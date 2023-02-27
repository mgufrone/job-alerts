package cache

import (
	"context"
	"encoding/json"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	KeySeparator = ":"
)

type redisCache struct {
	baseCache
	logger   *logrus.Entry
	redisCli *redis.Client
}

func (c redisCache) Has(ctx context.Context, key string) bool {
	k := c.makeKey(key)
	c.logger.Debugln("check if cache key exists", k)
	v, err := c.redisCli.Exists(ctx, k).Result()
	return err == nil && v != 0
}

func (c redisCache) Get(ctx context.Context, key string, out interface{}) (err error) {
	c.logger.Debugln("retrieve cache", c.makeKey(key))
	outString, err := c.redisCli.Get(ctx, c.makeKey(key)).Result()
	if err != nil {
		return
	}
	return json.Unmarshal([]byte(outString), out)
}

func (c redisCache) Set(ctx context.Context, key string, value interface{}, duration TTL) (err error) {
	if duration.resolve() == 0 {
		return
	}
	k := c.makeKey(key)
	c.logger.Debugln("setting cache", k)
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.logger.Debugf("storing cache %s with size %d", k, len(string(val)))
	return c.redisCli.Set(ctx, k, string(val), duration.resolve()).Err()
}
func (c redisCache) Delete(ctx context.Context, key string) (err error) {
	c.logger.Debugln("deleting cache", c.makeKey(key))
	return c.redisCli.Del(ctx, c.makeKey(key)).Err()
}

func (c redisCache) CacheOrCreate(ctx context.Context, key string, out interface{}, duration TTL, generator ICacheGenerator) (err error) {
	if c.Has(ctx, key) {
		return c.Get(ctx, key, out)
	}
	resultSet, err := generator(ctx)
	if err != nil {
		return
	}
	err = c.Set(ctx, key, resultSet, duration)
	if err != nil {
		return
	}
	return c.Get(ctx, key, out)
}

func (c redisCache) Clear(ctx context.Context) error {
	return c.ClearByPrefix(ctx, "")
}

func (c redisCache) ClearByPrefix(ctx context.Context, prefix string) error {
	kp := c.makeKey(fmt.Sprintf("%s:%s", prefix, "*"))
	if prefix == "" {
		kp = c.makeKey("*")
	}
	c.logger.Debugln("clearing cache", kp)
	iter := c.redisCli.Scan(ctx, 0, kp, 0).Iterator()
	for iter.Next(ctx) {
		c.logger.Debugln("deleting cache from prefix", iter.Val())
		err := c.redisCli.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func Redis(client *redis.Client, logger *logrus.Entry, prefix string) (ICache, error) {
	cache := &redisCache{redisCli: client}
	cache.logger = logger
	cache.prefix = []string{prefix}
	cache.separator = KeySeparator
	return cache, nil
}
