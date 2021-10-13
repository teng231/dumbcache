package dumbcache

import (
	"log"
	"testing"
	"time"
)

func TestLocalModuleCache(t *testing.T) {
	module := CreateLocalModule(10, 5*time.Second)
	module.Set("test1", "bimmer")
	module.Set("test2", map[string]interface{}{"name": "luna"})
	// get data
	data, err := module.Get("test1")
	log.Print(err, " ", data)

	data2, err := module.Get("test2")
	// m := make(map[string]interface{})
	// json.Unmarshal(data2, &m)
	log.Print(data2, err)

	time.Sleep(5 * time.Second)
	data, err = module.Get("test1")
	log.Print(err, " ", data)
}
