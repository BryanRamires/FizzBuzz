package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/BryanRamires/FizzBuzz/internal/config"
	"github.com/BryanRamires/FizzBuzz/internal/stats"
	goredis "github.com/redis/go-redis/v9"
)

var _ stats.Repository = (*Repo)(nil)

type Repo struct {
	cfg     config.Config
	rdb     *goredis.Client
	rankKey string
	opTO    time.Duration
}

func NewRepo(cfg config.Config, rdb *goredis.Client) *Repo {
	return &Repo{
		cfg:     cfg,
		rdb:     rdb,
		rankKey: "fizzbuzz:stats:rank",
		opTO:    cfg.RedisOpTimeout,
	}
}

func memberForKey(k stats.Key) (string, error) {
	b, err := json.Marshal(k)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *Repo) Inc(ctx context.Context, k stats.Key) error {
	ctx, cancel := context.WithTimeout(ctx, r.opTO)
	defer cancel()

	member, err := memberForKey(k)
	if err != nil {
		return err
	}

	return r.rdb.ZIncrBy(ctx, r.rankKey, 1, member).Err()
}

func (r *Repo) Top(ctx context.Context) (stats.Top, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.opTO)
	defer cancel()

	items, err := r.rdb.ZRevRangeWithScores(ctx, r.rankKey, 0, 0).Result()
	if err != nil || len(items) == 0 {
		return stats.Top{}, false, err
	}

	memberStr, ok := items[0].Member.(string)
	if !ok {
		return stats.Top{}, false, fmt.Errorf("invalid redis member type")
	}

	var params stats.Key
	if err := json.Unmarshal([]byte(memberStr), &params); err != nil {
		return stats.Top{}, false, err
	}

	hits := uint64(math.Round(items[0].Score))
	return stats.Top{Parameters: params, Hits: hits}, true, nil
}
