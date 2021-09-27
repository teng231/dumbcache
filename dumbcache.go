package dumbcache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type DumbCache struct {
	client   *redis.Client
	timeout  time.Duration
	duration time.Duration
}

func (d *DumbCache) Connect(addr, pw string, db int, timeout, duration time.Duration) error {
	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        pw,
		MaxRetries:      10,
		MinRetryBackoff: 15 * time.Millisecond,
		MaxRetryBackoff: 1000 * time.Millisecond,
		DialTimeout:     10 * time.Second,
		DB:              db, // use default DB
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Print(err)
	}
	d.client = client
	d.timeout = timeout
	d.duration = duration
	return nil
}

func (d *DumbCache) Set(prefix string, input, payload interface{}) error {
	hash, err := d.MakeHash(input)
	if err != nil {
		return err
	}
	pBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	return d.client.Set(ctx, prefix+hash, string(pBytes), d.duration).Err()
}

func (d *DumbCache) MakeHash(in interface{}) (string, error) {
	payload, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	hash := fmt.Sprintf("%x", md5.Sum(payload))
	return hash, nil
}

// ParseData out is a pointer
func (d *DumbCache) ParseData(input, out interface{}) error {
	hash, err := d.MakeHash(input)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	data, err := d.client.Get(ctx, hash).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return err
	}
	return nil
}

// List list cache data
// out is a pointer
func (d *DumbCache) List(input, out interface{}, handler func() (interface{}, error)) error {
	hash, err := d.MakeHash(input)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	isParsing := false
	data, err := d.client.Get(ctx, CLIST+hash).Result()
	if err != nil && err.Error() == redis.Nil.Error() {
		payload, err := handler()
		if err != nil {
			return err
		}
		if err := d.Set(CLIST, input, payload); err != nil {
			return err
		}
		bin, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = string(bin)
		isParsing = true
	}
	if err != nil && !isParsing {
		return err
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return err
	}
	return nil
}

// Count cache count
func (d *DumbCache) Count(input interface{}, out *int64, handler func() (int64, error)) error {
	hash, err := d.MakeHash(input)
	if err != nil {
		return err
	}
	isParsing := false
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	data, err := d.client.Get(ctx, CCOUNT+hash).Result()
	if err != nil && err.Error() == redis.Nil.Error() {
		payload, err := handler()
		if err != nil {
			return err
		}
		if err := d.Set(CCOUNT, input, payload); err != nil {
			return err
		}
		bin, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = string(bin)
		isParsing = true
	}
	if err != nil && !isParsing {
		return err
	}
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return err
	}
	return nil
}

// Expire expire 2 key count and list
func (d *DumbCache) Expire(input interface{}) error {
	hash, err := d.MakeHash(input)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.client.Del(ctx, CCOUNT+hash, CLIST+hash).Err()
}
