package redis

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"service/config"
	"strconv"
	"time"
)

type Redis struct {
	rdb *redis.Client
	rs  *redsync.Redsync
}

func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	return r.rdb.HGetAll(ctx, key).Scan(&dest)
}

func (r *Redis) Set(ctx context.Context, key string, data interface{}, second int) error {
	statusCmd := r.rdb.Set(ctx, key, data, time.Second*time.Duration(second))
	return statusCmd.Err()
}

//var _ cache.ICache = &Redis{}

func (r *Redis) Lock(ctx context.Context, key string, ttl int64, proses func(ctx context.Context) error) error {
	cfgExpiry := redsync.WithExpiry(time.Duration(ttl) * time.Second)

	mutex := r.rs.NewMutex(key, cfgExpiry)
	if err := mutex.LockContext(ctx); err != nil {
		return err
	}

	defer mutex.Unlock()

	err := proses(ctx)

	return err
}

func NewRedis(_ context.Context, config *config.Config) *Redis {

	dbName, err := strconv.Atoi(config.Redis.Name)
	if err != nil {
		dbName = 0
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password, // no password set
		DB:       dbName,                // use default DB
	})

	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)

	return &Redis{
		rdb: rdb,
		rs:  rs,
	}
}
