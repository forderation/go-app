package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/forderation/go-app/http_app"
)

type InMemoryPayerStore struct {
	sync  sync.Mutex
	Store map[string]int
}

func NewInMemoryPayerStore() *InMemoryPayerStore {
	return &InMemoryPayerStore{Store: map[string]int{}}
}

func (i *InMemoryPayerStore) GetPlayerScore(name string) int {
	return i.Store[name]
}

func (i *InMemoryPayerStore) RecordWin(name string) {
	// i.sync.Lock()
	i.Store[name]++
	// defer i.sync.Unlock()
}

func main() {
	server := &http_app.PlayerServer{
		Store: &InMemoryPayerStore{},
	}
	log.Fatal(http.ListenAndServe(":5000", server))
}
