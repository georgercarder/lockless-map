package lockless_map

import (
	"fmt"
)

type locklessMap struct {
	CH     (chan *kvPair)
	reqCH  (chan interface{}) // key
	takeCH (chan interface{}) // value
}

func NewLocklessMap() (lt *locklessMap) {
	lt = new(locklessMap)
	lt.CH = make(chan *kvPair, 1)
	lt.reqCH = make(chan interface{}, 1)
	lt.takeCH = make(chan interface{}, 1)
	go func() {
		latestMap := make(map[interface{}]interface{})
		kv := new(kvPair)
		var key interface{} 
		for {
			select {
			case kv = <-lt.CH:
				latestMap[kv.K] = kv.V
				continue
			case key = <-lt.reqCH:
				lt.takeCH <- latestMap[key]
				//c = nil
			}
		}
	}()
	return
}

type kvPair struct {
	K string
	V interface{}
}

func (lt *locklessMap) Take(key interface{}) (s interface{}, err error) {
	lt.reqCH <- key 
	s = <-lt.takeCH
	if s == nil {
		err = fmt.Errorf("*latest.take: Channel is empty.")
	}
	return
}

func (lt *locklessMap) Put(key string, s interface{}) {
//	fmt.Println("debug put", key, s)
	lt.CH <- &kvPair{K:key, V:s}
	return
}
