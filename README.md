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
		}
	}()

	// putter 2
	go func() {
		for i:=0; i < 1000; i++ {
			lm.Put("cat", i)
			time.Sleep(someOtherInterval)
		}
	}()

	time.Sleep(verySmallInterval) // so that Take wont give nil

	// taker
	t, err := lm.Take("cat")
	if err != nil {
	
	}
	fmt.Println("debug take", t, er)
```

### Please let me know in the comments if you encounter any issues or would like to contribute improvements.
