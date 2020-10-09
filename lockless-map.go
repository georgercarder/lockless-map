package lockless_map

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
			}
		}
	}()
	return
}

type kvPair struct {
	K interface{}
	V interface{}
}

func (lt *locklessMap) Take(keys ...interface{}) (s interface{}) {
	_lt := lt
	for _, k := range keys {
		s = _lt.take(k)
		switch s.(type) {
		case *locklessMap:
			_lt = s.(*locklessMap)
		default:
			return
		}
	}
	return
}

func (lt *locklessMap) take(key interface{}) (s interface{}) {
	lt.reqCH <- key
	s = <-lt.takeCH
	return
}

func (lt *locklessMap) Put(keysNVal ...interface{}) {
	_lt := lt
	for i := 0; i < len(keysNVal)-2; i++ {
		t := _lt.Take(keysNVal[i])
		if t == nil {
			t = NewLocklessMap()
		}
		_lt.put(keysNVal[i], t)
		_lt = t.(*locklessMap)
	}
	_lt.put(keysNVal[len(keysNVal)-2], keysNVal[len(keysNVal)-1])
	return
}

func (lt *locklessMap) put(key interface{}, value interface{}) {
	lt.CH <- &kvPair{K: key, V: value}
	return
}
