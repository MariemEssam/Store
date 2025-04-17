package main

import (
	"fmt"
	"sync"
	"time"
)

//Struct representing a purchase requests includes my buyers  name, product name, quantity, and a reply channel to get a response
// PurchaseRequest defines the structure of a customer request

type PurchaseRequest struct {
	BuyerName string
	Product   string
	Quantity  int
	ReplyChan chan string
}

// Store struct represents the store
// Holds available products and their quantities
// Uses a mutex to ensure safe concurrent access "No data races"

type Store struct {
	Products map[string]int
}

// ProcessPurchase handles a customer's purchase request
func (s *Store) ProcessPurchase(req PurchaseRequest) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	availableQty, exists := s.Products[req.Product]

	if !exists {
		req.ReplyChan <- fmt.Sprintf("%s Product %s not available", req.BuyerName, req.Product)
		return
	}

	if availableQty >= req.Quantity {
		s.Products[req.Product] -= req.Quantity
		req.ReplyChan <- fmt.Sprintf("%s bought %d of %s. Remaining: %d", req.BuyerName, req.Quantity, req.Product, s.Products[req.Product])
	} else if availableQty > 0 {
		req.ReplyChan <- fmt.Sprintf("%s found Not enough stock for %s. Available: %d", req.BuyerName, req.Product, availableQty)
	} else {
		req.ReplyChan <- fmt.Sprintf("%s %s is out of stock!", req.BuyerName, req.Product)
	}
}

// Buyer simulates a customer sending a purchase request
func Buyer(req PurchaseRequest, storeChan chan PurchaseRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	storeChan <- req
}

// StoreWorker continuously handles incoming purchase requests
func StoreWorker(store *Store, storeChan chan PurchaseRequest) {
	for req := range storeChan {
		store.ProcessPurchase(req)
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {

	// Initialize my store with some products and their quantity
	store := Store{
		Products: map[string]int{
			"Apple":  6,
			"Banana": 6,
			"Orange": 3,
		},
		Mutex: sync.Mutex{},
	}

	// Channel to send purchase requests to my store
	storeChan := make(chan PurchaseRequest)

	// Use goroutine to start my store worker
	go StoreWorker(&store, storeChan)

	// Identification of my buyers and their purchase requests
	buyers := []PurchaseRequest{
		{"Ali", "Apple", 2, make(chan string)},
		{"Mariem", "Apple", 4, make(chan string)},
		{"Farida", "Banana", 3, make(chan string)},
		{"Salim", "Orange", 5, make(chan string)},
		{"Yahia", "Kiwi", 1, make(chan string)},
		{"Noura", "Banana", 3, make(chan string)},
		{"Hadi", "Banana", 1, make(chan string)},
	}

	var wg sync.WaitGroup

	// Launch each buyer in a separate goroutine
	for _, buyer := range buyers {
		wg.Add(1)
		go Buyer(buyer, storeChan, &wg)
	}

	wg.Wait()        // Wait for all buyers to finish
	close(storeChan) // Close the channel after use
}
