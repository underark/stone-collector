package main

import (
	"net/http"
	"sync"
)

func main() {
	var wait sync.WaitGroup

	wait.Go(func() {
		http.Get("http://localhost:8080/trade")
	})

	wait.Go(func() {
		http.Get("http://localhost:8080/tick")
	})

	wait.Wait()
}
