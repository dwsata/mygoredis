package caches

import (
	"sync/atomic"
	"time"

	"github.com/opentechnologyself/mygoredis/helpers"
)

const NeverDie = 0 //ttl = 0
type value struct {
	Data  []byte //storage real data
	Ttl   int64  //time to live
	Ctime int64  //create time
}

func newValue(data []byte, ttl int64) *value {
	return &value{
		Data:  helpers.Copy(data),
		Ttl:   ttl,
		Ctime: time.Now().Unix(),
	}
}

func (v *value) alive() bool {
	return v.Ttl == NeverDie || time.Now().Unix()-v.Ctime < v.Ttl
}

func (v *value) visit() []byte {
	atomic.SwapInt64(&v.Ctime, time.Now().Unix())
	return v.Data
}
