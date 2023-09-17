# Limiting number of threads

## Why do we need to do this?

Becaue of the `go` keyword, the tasks are instantly offloaded to a seperate goroutine, thus the tcp server currenlty has negligible overhead to accept a request and can handle a lots of connections coming in. 

Go routines are very lightweight, can a lot of them can be spawed, however even they come with a cost, add a lot of them up and you can easily run of our resources. 

In such a scenario, our application will panic and die. Possibly also severing ties to all the connections which were accepted but not served yet. 

## Possible Solution

A solution would be to give up a bit of the performance and let some requests get accepted a bit later. We can do that by limiting the number of threads; in our case go routines. 

We can set a max cap to a calibrated number. So that  at a certain point in time we only have N amount of goroutines running. And at that moment if a new request comes in, they have to wait.


## Implementation

### Buffered Channels

When a buffered channel is full, it does not let more items to be sent until a slot gets freed up.
More on buffered channels [here](https://www.geeksforgeeks.org/buffered-channel-in-golang/)

We can utilise that feature to our advantage, and implement a [counting semaphore](https://www.geeksforgeeks.org/semaphores-in-process-synchronization/). If we keep pushing to a semaphore, and the semaphore is full, the next push will be blocked until one of the aquirees releases. 

And in the context of go channels, we will aquire by sending a value to the channel and release by reading a value from it.

Let's see a small example of that, you can run it on [replit](https://replit.com/@ShuvojitSarkar/BufferedChannel-as-a-Semaphore)

```
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


