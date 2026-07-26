package core

import (
	"github/litG-zen/litDB/conf"
	"time"
)

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64, oType uint8, oEnc uint8) *Obj {
	var expiresAt int64 = -1
	if durationMs > 1 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	return &Obj{
		Value:        value,
		ExpiresAt:    expiresAt,
		TypeEncoding: oType | oEnc,
	}
}

func Put(k string, obj *Obj) {
	if len(store) >= conf.KEYSLIMIT {
		evict()
	}
	store[k] = obj
}

func Get(k string) *Obj {
	v := store[k]
	if v != nil { //Passive deletion
		if v.ExpiresAt != -1 && v.ExpiresAt <= time.Now().UnixMilli() {
			delete(store, k)
			return nil
		}
	}
	return v

}

func Del(k string) bool {
	if _, ok := store[k]; ok {
		delete(store, k)
		return true
	}
	return false
}
