package proxy

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/peizhong/letsgo/pkg/log"
)

// 储存room和server的信息
type Repository interface {
	GetString(key string) (string, bool)
	SetString(key string, value string)
	DeleteString(key string)
	Keys(pattern string) []string

	GetSortedSetMemberCount(key string) int64
	GetSortedSetMemberScore(key, member string) (float64, bool)
	SetSortedSetMemberScore(key, member string, score float64)
	IncrSortedSetMemberScore(key, member string, incr float64)
	RemoveSortedSetMemberScoreRange(key string, min, max float64)
}

var (
	RedisHost     = "localhost"
	RedisPort     = 6379
	RedisPassword = ""
)

type RedisRepository struct {
	rc *redis.Client
}

func NewRedisRepository() Repository {
	// 单节点
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", RedisHost, RedisPort),
		Password: RedisPassword,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		panic(err)
	}
	return &RedisRepository{
		rc: client,
	}
}

func (r *RedisRepository) Keys(pattern string) []string {
	val, _ := r.rc.Keys(pattern).Result()
	return val
}

func (r *RedisRepository) GetString(key string) (string, bool) {
	val, err := r.rc.Get(key).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Info("redis err", err.Error())
		return "", false
	}
	return val, true
}

func (r *RedisRepository) SetString(key string, value string) {
	r.rc.Set(key, value, 0)
}

func (r *RedisRepository) DeleteString(key string) {
	r.rc.Del(key)
}

func (r *RedisRepository) GetSortedSetMemberCount(key string) int64 {
	val, _ := r.rc.ZCard(key).Result()
	return val
}

func (r *RedisRepository) GetSortedSetMemberScore(key, member string) (float64, bool) {
	val, err := r.rc.ZScore(key, member).Result()
	if err == redis.Nil {
		return 0, false
	} else if err != nil {
		log.Info("redis err", err.Error())
		return 0, false
	}
	return val, true
}

// SetSortedSet
func (r *RedisRepository) SetSortedSetMemberScore(key, member string, score float64) {
	r.rc.ZAdd(key, redis.Z{
		Member: member,
		Score:  score,
	})
}

// IncrSortedSetMemberScore redis: ZINCRBY
func (r *RedisRepository) IncrSortedSetMemberScore(key, member string, incr float64) {
	r.rc.ZIncrBy(key, incr, member)
}

// RemoveSortedSetScoreRange redis: ZREMRANGEBYSCORE
func (r *RedisRepository) RemoveSortedSetMemberScoreRange(key string, min, max float64) {
	r.rc.ZRemRangeByScore(key, fmt.Sprintf("%d", min), fmt.Sprintf("%d", max))
}
