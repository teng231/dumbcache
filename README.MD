# Dumbcache
[![GitHub Actions](https://github.com/princjef/gomarkdoc/workflows/Test/badge.svg)](https://github.com/teng231/dumbcache/actions?query=workflow)
[![Go Reference](https://pkg.go.dev/badge/github.com/princjef/gomarkdoc.svg)](https://pkg.go.dev/github.com/teng231/dumbcache)

Cache list, count with filter param golang, using struct, hashkey

### Structure

we hash your request object to md5 hashing and add a prefix count or list.
data storage to redis using get/set with expire time.(default by 2 mins)


### Installing

```bash
go get github.com/teng231/dumbcache
```

### Example
* LIST
``` go
func ListData() {
	d := &DumbCache{}
	err := d.Connect("localhost:6379", "", 0, 5*time.Second, 5*time.Minute)
	if err != nil {
		log.Print(1, err)
	}
	data := []*Partner{}
	err = d.List(&PartnerRequest{Id: 10, Limit: 2}, &data, func() (interface{}, error) {
        list, err := db.ListPartners(req)
		return list, err
	})
	if err != nil {
		log.Print(err)
	}
	log.Print(data)
}
```

* Count
``` go
func CountData(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect("localhost:6379", "", 0, 5*time.Second, 5*time.Minute)
	if err != nil {
		log.Print(1, err)
	}
	var c int64 = 0
	err = d.Count(&PartnerRequest{Id: 10, Limit: 3}, &c, func() (int64, error) {
		count, err := db.CountPartners(req)
		return count, err
	})
	if err != nil {
		log.Print(err)
	}
	log.Print(c)
}
```

* Expire
``` go
func Expire() {
	d := &DumbCache{}
	err := d.Connect("localhost:6379", "", 0, 5*time.Second, 5*time.Minute)
	if err != nil {
		log.Print(1, err)
	}
	var c int64 = 0
	err = d.Expire(&PartnerRequest{Id: 10, Limit: 3})
	if err != nil {
		log.Print(err)
	}
	log.Print(c)
}
```