## Lockless Map

### A threadsafe lockless map written in golang


### Usage

```
	lm := NewLocklessMap()
	// putter 1
	go func() {
		for i:=0; i < 1000; i++ {
			lm.Put("cat", i)
			time.Sleep(someInterval) 
			// simulating some period put
		}
	}()

	// putter 2
	go func() {
		for i:=0; i < 1000; i++ {
			lm.Put("cat", i)
			time.Sleep(someOtherInterval) 
			// simulating some periodic put
		}
	}()

	time.Sleep(verySmallInterval) // simulating activity in the application

	// taker
	t, err := lm.Take("cat")
	if err != nil {
		// err is only when t is nil	
	}
	fmt.Println("debug take", t, er)
```

### Please let me know in the comments if you encounter any issues or would like to contribute improvements.
