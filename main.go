package main

import (
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

var (
	MAX_THREADS     = 5
	sem             = make(chan int, MAX_THREADS)
	runningRoutines uint32
)

func randomDelayBetween(min int, max int) {
	rd := rand.Intn(max-min) + min
	log.Printf("I may take %d seconds to process...\n", rd)
	time.Sleep(time.Duration(rd) * time.Second)
}

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

func process(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	randomDelayBetween(1, 10)
	log.Println("processing the request")

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, World!\r\n"))
	conn.Close()
	Done()

	// Release
	<-sem
}

func main() {
	listener, err := net.Listen("tcp", ":1729")
	if err != nil {
		log.Fatal(err)
	}

	go runRoutineTracker()
	for {
		log.Println("ready to accept a new connection")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Aquire
		sem <- 1
		go process(conn)
		Add()
	}
}
