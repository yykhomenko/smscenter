package main

import (
	"fmt"
	"log"

	"github.com/yykhomenko/smscenter/internal/server"
)

func main() {
	s := server.NewServer("localhost:3736")
	defer s.Close()
	log.Println("smscenter listen to:", s.Addr())
	fmt.Scanln()
}
