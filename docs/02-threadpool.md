# Add threadpool to save on thread creation time

## Why do we need to do this?

Again, go routines are very lightweight but they're not free. No mater how minial, spawning a new one does require some resources.

According to the current code, we are creating a new routine for each new request and disposing it after processing.

## Possible Solution

A solution is to have a pool of worker threads who are ready to process requests and reuse them for incoming requests without disposing or creating any new ones. 

## Implementation

Basically need to create all workers at start of program, and then just keep sending connections to them via a channel.

Reference I am going to use: [Go by Example: Worker pool](https://gobyexample.com/worker-pools)
