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
				//c = nil
			}
		}
	}()
	return
}

type kvPair struct {
	K interface{}
	V interface{}
}

func (lt *locklessMap) Take(key interface{}) (s interface{}) {
	lt.reqCH <- key
	s = <-lt.takeCH
	return
}

func (lt *locklessMap) Put(key interface{}, s interface{}) {
	lt.CH <- &kvPair{K: key, V: s}
	return
}
