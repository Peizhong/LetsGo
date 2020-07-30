package proxy

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/peizhong/letsgo/pkg/log"
)

const RepoErrorVal float64 = -999

// 储存room和server的信息
type Repository interface {
	GetString(key string) (string, bool)
	SetString(key string, value string)
	DeleteString(key string)
	Keys(pattern string) []string

	GetSortedSetMemberCount(key string) int64
	GetSortedSetMemberScore(key, member string) float64
	SetSortedSetMemberScore(key, member string, score float64)
	IncrSortedSetMemberScore(key, member string, incr float64) float64
	RemoveSortedSetMember(key, member string)
	RemoveSortedSetMemberScoreRange(key string, min, max float64)
	GetSortedSetMembersWithScore(key string) map[string]float64
}

type RedisRepository struct {
	rc *redis.Client
}

func NewRedisRepository(conf Config) Repository {
	// 单节点
	client := redis.NewClient(&redis.Options{
		Addr:     conf.GetString(redisAddrKey),
		Password: conf.GetString(redisPasswordKey),
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

func (r *RedisRepository) GetSortedSetMemberScore(key, member string) float64 {
	val, err := r.rc.ZScore(key, member).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		log.Info("redis err", err.Error())
		return RepoErrorVal
	}
	return val
}

// SetSortedSet
func (r *RedisRepository) SetSortedSetMemberScore(key, member string, score float64) {
	r.rc.ZAdd(key, redis.Z{
		Member: member,
		Score:  score,
	})
}

// IncrSortedSetMemberScore redis: ZINCRBY
func (r *RedisRepository) IncrSortedSetMemberScore(key, member string, incr float64) float64 {
	val, err := r.rc.ZIncrBy(key, incr, member).Result()
	if err != nil {
		log.Info("redis err", err.Error())
		return RepoErrorVal
	}
	return val
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
