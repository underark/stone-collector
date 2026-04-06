package main

import (
	"net/http"
	"sync"
)

func main() {
	var wait sync.WaitGroup

	for range 100 {
		wait.Go(func() {
			http.Get("http://localhost:8080/tick")
		})
	}
	wait.Wait()
}
