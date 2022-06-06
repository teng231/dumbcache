package dumbcache

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type PartnerRequest struct {
	Id    int64   `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name  string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Type  int32   `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`
	State int32   `protobuf:"varint,5,opt,name=state,proto3" json:"state,omitempty"`
	Limit int32   `protobuf:"varint,7,opt,name=limit,proto3" json:"limit,omitempty"`
	From  int64   `protobuf:"varint,8,opt,name=from,proto3" json:"from,omitempty"`
	To    int64   `protobuf:"varint,9,opt,name=to,proto3" json:"to,omitempty"`
	Ids   []int64 `protobuf:"varint,11,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

type ARequest struct {
	Id    int64  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Limit int32  `protobuf:"varint,7,opt,name=limit,proto3" json:"limit,omitempty"`
}

type A struct {
	Id   string `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	X    int32  `protobuf:"varint,7,opt,name=limit,proto3" json:"x,omitempty"`
}

type Partner struct {
	// `xorm:"pk autoincr notnull"`
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" xorm:"pk autoincr notnull"`
	// `xorm:"text"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty" xorm:"text"`
	// `xorm:"text"`
	Address string `protobuf:"bytes,4,opt,name=address,proto3" json:"address,omitempty" xorm:"text"`
	Phone   string `protobuf:"bytes,5,opt,name=phone,proto3" json:"phone,omitempty"`
	// `xorm:"text"`
	Slogan string `protobuf:"bytes,6,opt,name=slogan,proto3" json:"slogan,omitempty" xorm:"text"`
}

type Partners struct {
	Partners []*Partner `protobuf:"varint,1,opt,name=partner,proto3" json:"partner,omitempty"`
}

func TestCacheGetSet(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:     "localhost:6379",
		Timeout:  5 * time.Second,
		Duration: 5 * time.Minute,
	})
	if err != nil {
		log.Print(1, err)
	}
	insert := []*Partner{
		{Id: 1, Name: "te1"},
		{Id: 2, Name: "te2"},
	}
	data := []*Partner{}
	req := &PartnerRequest{Id: 10, Limit: 1}
	err = d.Set("", req, insert)
	if err != nil {
		log.Print(2, err)
	}
	if err := d.ParseData(req, &data); err != nil {
		log.Print(err, req)
	}
	bin, _ := json.Marshal(data)
	log.Print(string(bin))
}

func TestCacheList(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:     "localhost:6379",
		Timeout:  5 * time.Second,
		Duration: 5 * time.Minute,
	})
	if err != nil {
		log.Print(1, err)
	}
	data := []*Partner{}
	err = d.List(&PartnerRequest{Id: 10, Limit: 2}, &data, func() (interface{}, error) {
		return []*Partner{
			{Id: 1, Name: "te1"},
			{Id: 4, Name: "te4"},
		}, nil
	})
	if err != nil {
		log.Print(2, err)
	}
	log.Print(data)
}

func TestCacheListWithProto(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:     "localhost:6379",
		Timeout:  5 * time.Second,
		Duration: 5 * time.Minute,
	})
	if err != nil {
		log.Print(1, err)
	}
	data := &Partners{}
	err = d.ListWithProto(&PartnerRequest{Id: 10, Limit: 2}, data, func() (proto.Message, error) {
		return &Partners{Partners: []*Partner{
			{Id: 1, Name: "te1"},
			{Id: 4, Name: "te4"},
		}}, nil
	})
	if err != nil {
		log.Print(2, err)
	}
	log.Print(data)
}

func TestCacheCount(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:     "localhost:6379",
		Timeout:  5 * time.Second,
		Duration: 5 * time.Minute,
	})
	if err != nil {
		log.Print(1, err)
	}
	var c int64 = 0
	err = d.Count(&PartnerRequest{Id: 10, Limit: 3}, &c, func() (int64, error) {
		return 20, nil
	})
	if err != nil {
		log.Print(2, err)
	}
	log.Print(c)
}

func TestCacheListWithRedisDie(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:     "localhost:6379",
		Timeout:  5 * time.Second,
		Duration: 5 * time.Minute,
	})
	if err != nil {
		log.Print(1, err)
	}
	d.client.Close()
	data := []*Partner{}
	err = d.List(&PartnerRequest{Id: 10, Limit: 2}, &data, func() (interface{}, error) {
		return []*Partner{
			{Id: 1, Name: "te1"},
			{Id: 4, Name: "te4"},
		}, nil
	})
	if err != nil {
		log.Print(2, err)
	}
	log.Print(data)
}

func TestCacheListWithLocalCache(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:          "localhost:6379",
		Timeout:       5 * time.Second,
		Duration:      5 * time.Minute,
		LocalDuration: 1 * time.Second,
		MaxSizeLocal:  200,
	})
	if err != nil {
		log.Print(1, err)
	}
	data := []*Partner{}
	err = d.List(&PartnerRequest{Id: 10, Limit: 2}, &data, func() (interface{}, error) {
		return []*Partner{
			{Id: 1, Name: "te1"},
			{Id: 4, Name: "te4"},
		}, nil
	})
	if err != nil {
		log.Print(2, err)
	}
	log.Print(data)
	time.Sleep(100 * time.Millisecond)

	/// next step
	d.List(&PartnerRequest{Id: 10, Limit: 2}, &data, func() (interface{}, error) {
		return []*Partner{
			{Id: 1, Name: "te1"},
			{Id: 4, Name: "te4"},
		}, nil
	})
	log.Print(data)
	time.Sleep(2 * time.Second)

	/// next step
	d.List(&PartnerRequest{Id: 10, Limit: 2}, &data, func() (interface{}, error) {
		return []*Partner{
			{Id: 1, Name: "te1"},
			{Id: 4, Name: "te4"},
		}, nil
	})

	log.Print(data)
}

func TestCacheItemWithFreeKey(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect(&Config{
		Addr:          "localhost:6379",
		Timeout:       5 * time.Second,
		Duration:      5 * time.Minute,
		LocalDuration: 1 * time.Second,
		MaxSizeLocal:  200,
	})
	if err != nil {
		log.Print(1, err)
	}

	cases := map[interface{}]interface{}{
		&PartnerRequest{Id: 10, Limit: 2}: []*Partner{
			{Id: 1, Name: "te1"},
		},
		&ARequest{Id: 10, Limit: 2}: []*A{
			{Id: "ss", Name: "2"},
		},
	}
	for req, item := range cases {
		var data interface{}
		err = d.List(req, &data, func() (interface{}, error) {
			return item, nil
		})
		if err != nil {
			log.Print(2, err)
		}
		log.Print(data)
	}
}

func TestRun(t *testing.T) {
	x := "11.222"
	y := "333"
	log.Print(strings.Split(x, "."))
	log.Print(strings.Split(y, "."))
	xx := strings.Split(x, ".")
	yy := strings.Split(y, ".")
	log.Print(xx[len(xx)-1])
	log.Print(yy[len(yy)-1])

}
