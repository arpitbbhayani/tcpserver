package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	WORKERS_COUNT = 100
	wg            sync.WaitGroup
)

func randomDelayBetween(min int, max int) {
	rd := rand.Intn(max-min) + min
	log.Printf("I may take %d seconds to process...\n", rd)
	time.Sleep(time.Duration(rd) * time.Second)
}

func process(ctx context.Context, workerId int, connChan chan net.Conn) {
	defer wg.Done()
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
			log.Printf("[%d] Worker dying, bye", workerId)
			return
		}
	}
}

func main() {
	wg.Add(WORKERS_COUNT)
	listener, err := net.Listen("tcp", ":1729")
	if err != nil {
		log.Fatal(err)
	}

	connChan := make(chan net.Conn)
	defer close(connChan)

	ctx, cancelWorkers := context.WithCancel(context.TODO())
	defer cancelWorkers()

	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigC
		cancelWorkers()
		wg.Wait()

		log.Println("Exiting..")
		os.Exit(0)
	}()

	go func(ctx context.Context, c chan net.Conn) {
		// Start workers
		for i := 0; i < WORKERS_COUNT; i++ {
			go process(ctx, i, c)
		}
	}(ctx, connChan)

	for {
		log.Println("ready to accept a new connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		connChan <- conn
	}

}
