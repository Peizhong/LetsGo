package proxy

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/peizhong/letsgo/pkg/log"
)

// 储存room和server的信息
type Repository interface {
	GetString(key string) (string, bool)
	SetString(key string, value string, expiration time.Duration) error
	DeleteString(key string)
	Keys(pattern string) []string

	GetSortedSetMemberCount(key string) int64
	GetSortedSetMemberScore(key, member string) float64
	SetSortedSetMemberScore(key, member string, score float64) error
	IncrSortedSetMemberScore(key, member string, incr float64) (float64, error)
	RemoveSortedSetMember(key, member string)
	RemoveSortedSetMemberScoreRange(key string, min, max float64)
	GetSortedSetMembersWithScore(key string) map[string]float64
}

type RedisRepository struct {
	rc *redis.Client
}

var onceRepo sync.Once
var singelRedis Repository

func NewRedisRepository(conf Config) Repository {
	if singelRedis != nil {
		return singelRedis
	}
	onceRepo.Do(func() {
		// 单节点
		client := redis.NewClient(&redis.Options{
			Addr:     conf.GetString(redisAddrKey),
			Password: conf.GetString(redisPasswordKey),
			DB:       0,
		})
		if err := client.Ping().Err(); err != nil {
			panic(err)
		}
		singelRedis = &RedisRepository{
			rc: client,
		}
	})
	return singelRedis
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

func (r *RedisRepository) SetString(key string, value string, expiration time.Duration) error {
	err := r.rc.Set(key, value, expiration).Err()
	return err
}

func (r *RedisRepository) DeleteString(key string) {
	r.rc.Del(key)
}

func (r *RedisRepository) GetSortedSetMemberCount(key string) int64 {
	val, _ := r.rc.ZCard(key).Result()
	return val
}

func (r *RedisRepository) GetSortedSetMemberScore(key, member string) float64 {
	val, err := r.rc.ZScore(key, member).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		log.Info("redis err", err.Error())
		return -1
	}
	return val
}

// SetSortedSet
func (r *RedisRepository) SetSortedSetMemberScore(key, member string, score float64) error {
	err := r.rc.ZAdd(key, redis.Z{
		Member: member,
		Score:  score,
	}).Err()
	return err
}

// IncrSortedSetMemberScore redis: ZINCRBY
func (r *RedisRepository) IncrSortedSetMemberScore(key, member string, incr float64) (float64, error) {
	val, err := r.rc.ZIncrBy(key, incr, member).Result()
	if err != nil {
		log.Info("redis err", err.Error())
		return -1, err
	}
	return val, nil
}

// RemoveSortedSetMember redis: ZREM
func (r *RedisRepository) RemoveSortedSetMember(key, member string) {
	r.rc.ZRem(key, member)
}

// RemoveSortedSetScoreRange redis: ZREMRANGEBYSCORE
func (r *RedisRepository) RemoveSortedSetMemberScoreRange(key string, min, max float64) {
	r.rc.ZRemRangeByScore(key, fmt.Sprintf("%v", min), fmt.Sprintf("%v", max))
}

func (r *RedisRepository) GetSortedSetMembersWithScore(key string) map[string]float64 {
	res := make(map[string]float64)
	vals, err := r.rc.ZRangeWithScores(key, 0, -1).Result()
	if err == redis.Nil {
		return res
	} else if err != nil {
		log.Info("redis err", err.Error())
		return res
	}
	for _, v := range vals {
		res[v.Member.(string)] = v.Score
	}
	return res
}
