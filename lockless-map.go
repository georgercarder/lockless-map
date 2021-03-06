package lockless_map

import (
	"time"
)

type LocklessMap interface {
	Put(keysNVal ...interface{})
	Take(keys ...interface{}) interface{}
	Dump() *DumpPacket
}

type LocklessMap_ struct {
	cH           (chan *kvPair)
	reqCH        (chan interface{}) // key
	takeCH       (chan interface{}) // value
	dumpReqCH    (chan bool)
	dumpPacketCH (chan *DumpPacket)
	tCH          (chan bool)
}

type DumpPacket struct {
	Keys   []interface{}
	Values []interface{}
}

func NewLocklessMap() (lt *LocklessMap_) {
	lt = new(LocklessMap_)
	lt.cH = make(chan *kvPair, 1)
	lt.reqCH = make(chan interface{}, 1)
	lt.takeCH = make(chan interface{}, 1)
	lt.dumpReqCH = make(chan bool, 1)
	lt.dumpPacketCH = make(chan *DumpPacket, 1)
	lt.tCH = make(chan bool, 1)
	go func() {
		latestMap := make(map[interface{}]interface{})
		kv := new(kvPair)
		var key interface{}
		for {
			select {
			case kv = <-lt.cH:
				latestMap[kv.K] = kv.V
				continue
			case key = <-lt.reqCH:
				lt.takeCH <- latestMap[key]
				continue
			case <-lt.dumpReqCH:
				dp := new(DumpPacket)
				for k, v := range latestMap {
					dp.Keys = append(dp.Keys, k)
					dp.Values = append(dp.Values, v)
				}
				lt.dumpPacketCH <- dp
			}
		}
	}()
	return
}

type kvPair struct {
	K interface{}
	V interface{}
}

func (lt *LocklessMap_) Take(keys ...interface{}) (s interface{}) {
	_lt := lt
	for _, k := range keys {
		s = _lt.take(k)
		switch s.(type) {
		case *LocklessMap_:
			_lt = s.(*LocklessMap_)
		default:
			return
		}
	}
	return
}

func (lt *LocklessMap_) take(key interface{}) (s interface{}) {
	lt.reqCH <- key
	s = <-lt.takeCH
	return
}

func (lt *LocklessMap_) Dump() (dp *DumpPacket) {
	lt.dumpReqCH <- true
	dp = <-lt.dumpPacketCH
	return
}

func (lt *LocklessMap_) Put(keysNVal ...interface{}) {
	_lt := lt
	for i := 0; i < len(keysNVal)-2; i++ {
		t := _lt.Take(keysNVal[i])
		if t == nil {
			t = NewLocklessMap()
		}
		_lt.put(keysNVal[i], t)
		_lt = t.(*LocklessMap_)
	}
	_lt.put(keysNVal[len(keysNVal)-2], keysNVal[len(keysNVal)-1])
	return
}

func (lt *LocklessMap_) put(key interface{}, value interface{}) {
	for len(lt.cH) > 0 {
		time.Sleep(time.Nanosecond)
	}
	lt.cH <- &kvPair{K: key, V: value}
	return
}
