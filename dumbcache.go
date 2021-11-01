package dumbcache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func (d *DumbCache) Connect(config *Config) error {
	if config.Addr == "" {
		return errors.New("not found addr redis")
	}
	if config.Duration == 0 {
		return errors.New("not found duration")
	}
	client := redis.NewClient(&redis.Options{
		Addr:            config.Addr,
		Password:        config.Password,
		MaxRetries:      10,
		MinRetryBackoff: 15 * time.Millisecond,
		MaxRetryBackoff: 1000 * time.Millisecond,
		DialTimeout:     10 * time.Second,
		DB:              config.Db, // use default DB
	})
	if config.Timeout == 0 {
		config.Timeout = 3 * time.Second
	}
	if config.MaxSizeLocal == 0 {
		config.MaxSizeLocal = 300
	}
	if config.LocalDuration == 0 {
		config.LocalDuration = config.Duration / 2
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Print(err)
	}
	d.client = client
	d.timeout = config.Timeout
	d.duration = config.Duration
	d.module = CreateLocalModule(config.MaxSizeLocal, config.LocalDuration)
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
	if d.module != nil {
		err := d.module.Set(prefix+hash, payload)
		if err != nil {
			log.Print(err)
		}
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

	if d.module != nil {
		data, err := d.module.Get(hash)
		if err == nil {
			if err := json.Unmarshal([]byte(data), out); err != nil {
				return err
			}
			return nil
		}

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

	if d.module != nil {
		data, err := d.module.Get(CLIST + hash)
		if err == nil {
			if err := json.Unmarshal([]byte(data), out); err != nil {
				return err
			}
		}

	}

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
		payload, err := handler()
		if err != nil {
			return err
		}
		bin, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = string(bin)
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
	if d.module != nil {
		data, err := d.module.Get(CCOUNT + hash)
		if err == nil {
			if err := json.Unmarshal([]byte(data), out); err != nil {
				return err
			}
			return nil
		}
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
		log.Print(err)
		payload, err := handler()
		if err != nil {
			return err
		}
		bin, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = string(bin)
	}
	log.Print("get from redis")
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
	d.module.Del(CCOUNT + hash)
	d.module.Del(CLIST + hash)
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	return d.client.Del(ctx, CCOUNT+hash, CLIST+hash).Err()
}

// CalcInt cache count
func (d *DumbCache) CalcInt(input interface{}, out *int64, handler func() (int64, error)) error {
	hash, err := d.MakeHash(input)
	if err != nil {
		return err
	}
	if d.module != nil {
		data, err := d.module.Get(CCACLINT + hash)
		if err == nil {
			if err := json.Unmarshal([]byte(data), out); err != nil {
				return err
			}
			return nil
		}
	}
	isParsing := false
	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	data, err := d.client.Get(ctx, CCACLINT+hash).Result()
	if err != nil && err.Error() == redis.Nil.Error() {
		payload, err := handler()
		if err != nil {
			return err
		}
		if err := d.Set(CCACLINT, input, payload); err != nil {
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
		log.Print(err)
		payload, err := handler()
		if err != nil {
			return err
		}
		bin, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		data = string(bin)
	}
	log.Print("get from redis")
	if err := json.Unmarshal([]byte(data), out); err != nil {
		return err
	}
	return nil
}
