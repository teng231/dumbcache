package dumbcache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type IDumbCache interface {
	// connect redis
	Connect(config *Config) error

	// Set raw to get raw
	Set(prefix string, input, payload interface{}) error
	// Make hash key
	MakeHash(in interface{}) (string, error)
	ParseData(input, out interface{}) error

	Expire(input interface{}) error
	List(input, out interface{}, handler func() (interface{}, error)) error
	Count(input interface{}, out *int64, handler func() (int64, error)) error
}

type DumbCache struct {
	client   *redis.Client
	timeout  time.Duration
	duration time.Duration
	module   *LocalModule
}

type Config struct {
	Addr          string
	Password      string
	Db            int
	Timeout       time.Duration
	Duration      time.Duration
	MaxSizeLocal  int
	LocalDuration time.Duration
}

const (
	CLIST  = "c_list_"
	CCOUNT = "c_count_"
)
