package main

import (
	"log"

	"github.com/xjdrew/gofin"
)

func main() {
	s, err := gofin.GetLastPrice("sh600000")
	if err != nil {
		log.Fatal("gofin.GetLastPrice:", err)
	}
	log.Printf("%+v", s)
}
