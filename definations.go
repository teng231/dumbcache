package dumbcache

import (
	"time"
)

type IDumbCache interface {
	// connect redis
	Connect(addr, pw string, db int, timeout, duration time.Duration) error
	// Set raw to get raw
	Set(prefix string, input, payload interface{}) error
	// Make hash key
	MakeHash(in interface{}) (string, error)
	ParseData(input, out interface{}) error

	Expire(input interface{}) error
	List(input, out interface{}, handler func() (interface{}, error)) error
	Count(input interface{}, out *int64, handler func() (int64, error)) error
}

const (
	CLIST  = "c_list_"
	CCOUNT = "c_count_"
)
