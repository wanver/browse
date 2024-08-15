package main

// Package main ...

import (
	"log"

	"github.com/wanver/browse"
)

func main() {
	err := browse.App()
	if err != nil {
		log.Fatal(err)
	}

}
