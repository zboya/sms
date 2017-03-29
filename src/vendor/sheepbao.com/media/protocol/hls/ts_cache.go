package hls

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
)

type TSCache struct {
	// lock   sync.RWMutex
	entrys map[string]*TSCacheItem
}

func NewTSCache() *TSCache {
	return &TSCache{
		entrys: make(map[string]*TSCacheItem),
	}
}

func (self *TSCache) Set(key string, e *TSCacheItem) {
	// self.lock.Lock()
	v, ok := self.entrys[key]
	if !ok {
		self.entrys[key] = e
	}
	if v.ID() != e.ID() {
		self.entrys[key] = e
	}
	// self.lock.Unlock()
}

func (self *TSCache) Get(key string) *TSCacheItem {
	// self.lock.Lock()
	v := self.entrys[key]
	// self.lock.Unlock()
	return v
}

const (
	maxTSCacheNum = 3
)

var (
	ErrNoKey = errors.New("No key for cache")
)

type TSCacheItem struct {
	id  string
	num int
	ll  *list.List
	// lock sync.RWMutex
	lm map[string]TSItem
}

func NewTSCacheItem(id string) *TSCacheItem {
	return &TSCacheItem{
		id:  id,
		ll:  list.New(),
		num: maxTSCacheNum,
		lm:  make(map[string]TSItem),
	}
}

func (self *TSCacheItem) ID() string {
	return self.id
}

func (self *TSCacheItem) GenM3U8PlayList() ([]byte, error) {
	var seq int
	var getSeq bool
	var maxDuration int
	m3u8body := bytes.NewBuffer(nil)
	// self.lock.Lock()
	for e := self.ll.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		v, ok := self.lm[key]
		if ok {
			if v.Duration > maxDuration {
				maxDuration = v.Duration
			}
			if !getSeq {
				getSeq = true
				seq = v.SeqNum
			}
			fmt.Fprintf(m3u8body, "#EXTINF:%.3f,\n%s\n", float64(v.Duration)/float64(1000), v.Name)
		}
	}
	// self.lock.Unlock()
	w := bytes.NewBuffer(nil)
	fmt.Fprintf(w,
		"#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-ALLOW-CACHE:NO\n#EXT-X-TARGETDURATION:%d\n#EXT-X-MEDIA-SEQUENCE:%d\n\n",
		maxDuration/1000+1, seq)
	w.Write(m3u8body.Bytes())
	return w.Bytes(), nil
}

func (self *TSCacheItem) SetItem(key string, item TSItem) {
	// self.lock.Lock()
	if self.ll.Len() == self.num {
		e := self.ll.Front()
		self.ll.Remove(e)
		k := e.Value.(string)
		delete(self.lm, k)
	}
	self.lm[key] = item
	self.ll.PushBack(key)
	// self.lock.Unlock()
}

func (self *TSCacheItem) GetItem(key string) (TSItem, error) {
	// self.lock.RLock()
	item, ok := self.lm[key]
	if !ok {
		// self.lock.RUnlock()
		return item, ErrNoKey
	}
	// self.lock.RUnlock()
	return item, nil
}

type TSItem struct {
	Name     string
	SeqNum   int
	Duration int
	Data     []byte
}

func NewTSItem(name string, duration, seqNum int, b []byte) TSItem {
	var item TSItem
	item.Name = name
	item.SeqNum = seqNum
	item.Duration = duration
	item.Data = make([]byte, len(b))
	copy(item.Data, b)
	return item
}
