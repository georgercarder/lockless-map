## Lockless Map

### A simple threadsafe lockless hashmap written in golang


### Usage

```
	lm := NewLocklessMap()

	key1 := "cat" // key to be "contended" for
	key2 := 123
	key3 := uint64(456)
	key4 := byte(0xF4)

	// putter 1
	go func() {
		for i:=0; i < 1000; i++ {
			value := i
			lm.Put(key1, key2, key3, key4, value)
			time.Sleep(someInterval) 
			// simulating some periodic put
		}
	}()

	// putter 2
	go func() {
		for i:=0; i < 1000; i++ {
			value := i
			lm.Put(key1, key2, key3, key4, value)
			time.Sleep(someOtherInterval) 
			// simulating some periodic put
		}
	}()

	time.Sleep(verySmallInterval) // simulating activity in the application

	// taker
	t, err := lm.Take(key1, key2, key3, key4)
	if err != nil {
		// err is only when t is nil	
	}
	// do stuff with t
```

### Please let me know in git issues if you encounter any issues or would like to contribute improvements.
