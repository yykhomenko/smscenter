package main

import (
	"fmt"
	"log"

	"github.com/cbi-sh/smscenter/internal/app/server"
)

func main() {
	s := server.NewServer("192.168.0.2:3736")
	defer s.Close()

	log.Println("smscenter listen to:", s.Addr())

	fmt.Scanln()
}
