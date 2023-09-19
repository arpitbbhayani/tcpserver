# Add threadpool to save on thread creation time

## Why do we need to do this?



## Possible Solution

A solution would be to give up a bit of the performance and let some requests get accepted a bit later. We can do that by limiting the number of threads; in our case go routines. 

We can set a max cap to a calibrated number. So that  at a certain point in time, we only have N amount of goroutines running. And at that moment if a new request comes in, they have to wait.


## Implementation

### Buffered Channels

When a buffered channel is full, it does not let more items to be sent until a slot gets freed up.
More on buffered channels [here](https://www.geeksforgeeks.org/buffered-channel-in-golang/)

We can utilize that feature to our advantage, and implement a [counting semaphore](https://www.geeksforgeeks.org/semaphores-in-process-synchronization/). If we keep pushing to a semaphore, and the semaphore is full, the next push will be blocked until one of the acquirers releases. 

And in the context of go channels, we will acquire by sending a value to the channel and release by reading a value from it.

Let's see a small example of that, you can run it on [replit](https://replit.com/@ShuvojitSarkar/BufferedChannel-as-a-Semaphore)

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	MAX_GOROUTINES = 3
	sem            = make(chan int, MAX_GOROUTINES)
	wg             sync.WaitGroup
)

func process(i int) {
	fmt.Println("Processing: ", i)
	time.Sleep(3 * time.Second)
	fmt.Println("Processed: ", i)

	wg.Done()

	// Release the semaphore by receving a value from the channel
	<-sem
}

func main() {
	for i := 1; i <= 100; i++ {
        // Aquire a semaphore by sending a value to the channel
		sem <- 1

		wg.Add(1)
		go process(i)
	}

	// Wait for all workers to finish
	wg.Wait()
}
```

## Tracking Concurrent Running Goroutines

In order to debug and see if the solution actually works, we can have a thread-safe int to track whenever a goroutine starts or ends. 

We can accomplish that by using the [atomic](https://pkg.go.dev/sync/atomic) package. And we can spin up another go routine that will start a ticker to print the currently running goroutines.


```go
...

var runningRoutines int

func Add() {
	atomic.AddUint32(&runningRoutines, 1)
}

func Done() {
	atomic.AddUint32(&runningRoutines, ^uint32(0))
}

func runRoutineTracker() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("Currently running routines: ", runningRoutines)
		}
	}
}
...
```
