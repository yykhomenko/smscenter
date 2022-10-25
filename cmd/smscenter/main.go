package main

import (
	"fmt"
	"log"

	"github.com/yykhomenko/smscenter/pkg/server"
)

func main() {
	s := server.NewServer("localhost:2775")
	defer s.Close()
	log.Println("smscenter listen to:", s.Addr())
	fmt.Scanln()
}
