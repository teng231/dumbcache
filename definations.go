package dumbcache

import (
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type IDumbCache interface {
	// connect redis
	Connect(config *Config) error

	// Set raw to get raw
	Set(prefix string, input, payload interface{}) error
	SetWithProto(prefix string, input interface{}, payload proto.Message) error
	// Make hash key
	MakeHash(in interface{}) (string, error)
	ParseData(input, out interface{}) error
	ParseDataWithProto(input interface{}, payload proto.Message) error

	Expire(input interface{}) error
	List(input, out interface{}, handler func() (interface{}, error)) error
	ListWithProto(input, out interface{}, handler func() (proto.Message, error)) error

	Count(input interface{}, out *int64, handler func() (int64, error)) error
	CalcInt(input interface{}, out *int64, handler func() (int64, error)) error
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
	CLIST    = "c_list_"
	CCOUNT   = "c_count_"
	CCACLINT = "c_calc_int_"
)
