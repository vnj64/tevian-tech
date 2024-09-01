package main

import (
	"sync"
	"tevian/app"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.NewHttpServer().Start()
	}()

	wg.Wait()
}
