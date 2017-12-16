package main

import (
	"time"

	"github.com/sh3rp/shooze"
)

func main() {
	svc := shooze.NewWebservice()
	svc.Init(8080)
	for {
		time.Sleep(1 * time.Second)
	}
}
