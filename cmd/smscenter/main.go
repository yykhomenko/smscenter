package main

import (
	"time"

	"github.com/cbi-sh/smscenter/internal/app/server"
)

func main() {
	s := server.NewServer("192.168.0.2:3736")
	defer s.Close()

	time.Sleep(1<<63 - 1)
}
