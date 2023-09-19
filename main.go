package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"
)

var (
	WORKERS_COUNT = 100
	connChan      = make(chan net.Conn)
)

func randomDelayBetween(min int, max int) {
	rd := rand.Intn(max-min) + min
	log.Printf("I may take %d seconds to process...\n", rd)
	time.Sleep(time.Duration(rd) * time.Second)
}

func process(ctx context.Context, workerId int) {
	log.Printf("Worker %d on duty!\n", workerId)
	for {
		select {
		case conn := <-connChan:
			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			randomDelayBetween(1, 2)
			log.Printf("[%d] processing the request", workerId)
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, World!\r\n"))
			conn.Close()

			break
		case <-ctx.Done():
			log.Println("Worker dying, bye")
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":1729")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelWorkers := context.WithCancel(context.TODO())
	defer cancelWorkers()

	go func(ctx context.Context) {
		// Start workers
		for i := 0; i < WORKERS_COUNT; i++ {
			go process(ctx, i)
		}
	}(ctx)

	for {
		log.Println("ready to accept a new connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		connChan <- conn
	}
}
