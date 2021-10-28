package dumbcache

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

// Localmodule is addition module for add to and reduce request traffic to redis

type Message struct {
	Created int64
	Payload string
}

type ILocalModule interface {
	Add(key interface{}, value interface{}) error
	Get(key interface{}) (string, error)
}

type LocalModule struct {
	s        *lru.Cache
	duration time.Duration
}

func CreateLocalModule(size int, duration time.Duration) *LocalModule {
	s, err := lru.New(size)
	if err != nil {
		panic(err)
	}
	return &LocalModule{s: s, duration: duration}
}

func (l *LocalModule) Set(key interface{}, value interface{}) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	l.s.Add(key, Message{
		Created: time.Now().Unix(),
		Payload: string(payload),
	})
	return nil
}

func (l *LocalModule) Get(key interface{}) (string, error) {
	value, has := l.s.Get(key)
	if !has {
		return "", errors.New("not found")
	}
	msg := value.(Message)
	if time.Since(time.Unix(msg.Created, 0)) > l.duration {

		return "", errors.New("time out")
	}
	log.Print("get by localmodule")
	return msg.Payload, nil
}

func (l *LocalModule) Del(key interface{}) {
	l.s.Remove(key)
}
