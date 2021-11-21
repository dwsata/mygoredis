package caches

import (
	"sync/atomic"
	"time"

	"github.com/opentechnologysel/mygoredis/helpers"
)

const NeverDie = 0 //ttl = 0
type value struct {
	data  []byte //storage real data
	ttl   int64  //time to live
	ctime int64  //create time
}

func newValue(data []byte, ttl int64) *value {
	return &value{
		data:  helpers.Copy(data),
		ttl:   ttl,
		ctime: time.Now().Unix(),
	}
}

func (v *value) alive() bool {
	return v.ttl == NeverDie || time.Now().Unix()-v.ctime < v.ttl
}

func (v *value) visit() []byte {
	atomic.SwapInt64(&v.ctime, time.Now().Unix())
	return v.data
}
