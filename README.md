## Lockless Map

### A simple threadsafe lockless map written in golang


### Usage

```
	lm := NewLocklessMap()

	key := "cat" // key to be "contended" for

	// putter 1
	go func() {
		for i:=0; i < 1000; i++ {
			value := i
			lm.Put(key, value)
			time.Sleep(someInterval) 
			// simulating some periodic put
		}
	}()

	// putter 2
	go func() {
		for i:=0; i < 1000; i++ {
			value := i
			lm.Put(key, value)
			time.Sleep(someOtherInterval) 
			// simulating some periodic put
		}
	}()

	time.Sleep(verySmallInterval) // simulating activity in the application

	// taker
	t, err := lm.Take(key)
	if err != nil {
		// err is only when t is nil	
	}
	// do stuff with t
```

### Please let me know in the comments if you encounter any issues or would like to contribute improvements.
