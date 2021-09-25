package dumbcache

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

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

func TestCacheGetSet(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect("localhost:6379", "", 0, 5*time.Second, 5*time.Minute)
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
		log.Print(err)
	}
	bin, _ := json.Marshal(data)
	log.Print(string(bin))
}

func TestCacheList(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect("localhost:6379", "", 0, 5*time.Second, 5*time.Minute)
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

func TestCacheCount(t *testing.T) {
	d := &DumbCache{}
	err := d.Connect("localhost:6379", "", 0, 5*time.Second, 5*time.Minute)
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
