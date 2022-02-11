# v1.1.4
Fix Duplicate key with md5 hash
see at test `TestCacheItemWithFreeKey`

# v1.1.3
add new cache type calcInt

# v1.1.2
remove log

# v1.1.1
### Breaking change
```go
// change Connect func from

func (d *DumbCache) Connect(addr, pw string, db int, timeout, duration time.Duration) error

// to new connect func

func (d *DumbCache) Connect(config *Config) error


```

Add struct `Config`

```go
type Config struct {
	Addr          string
	Password      string
	Db            int
	Timeout       time.Duration
	Duration      time.Duration
	MaxSizeLocal  int
	LocalDuration time.Duration
}
```
Modify Dumbcache struct add option localModule

# v1.0.3
Add dumbcache release