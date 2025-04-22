package main

import (
	"fmt"
	"sync"
	"time"
)

type PurchaseRequest struct {
	BuyerName string
	Product   string
	Quantity  int
	// ReplyChan chan string (removed )
}

type Store struct {
	Products map[string]int
	Mutex    sync.Mutex
	Channel  chan PurchaseRequest
}

type ProductItem struct {
	Name     string
	Quantity int
}

func NewStore(responseChan chan string, pi ...ProductItem) *Store {
	store := &Store{
		Products: make(map[string]int),
		Channel:  make(chan PurchaseRequest, 20),
	}
	for _, item := range pi {
		store.Products[item.Name] = item.Quantity
	}
	go store.storeWorker(responseChan) //here returns responseChan
	return store
}

func (s *Store) ProcessPurchase(req PurchaseRequest, responseChan chan string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	availableQty, exists := s.Products[req.Product]

	if !exists {
		responseChan <- fmt.Sprintf("%s: Product %s not available", req.BuyerName, req.Product) //check if respinsechan exists
		return
	}
	// here I've checked with responseChan the only channel for all buyers
	if availableQty >= req.Quantity {
		s.Products[req.Product] -= req.Quantity
		responseChan <- fmt.Sprintf("%s bought %d of %s. Remaining: %d", req.BuyerName, req.Quantity, req.Product, s.Products[req.Product])
	} else if availableQty > 0 {
		responseChan <- fmt.Sprintf("%s: Not enough stock for %s. Available: %d", req.BuyerName, req.Product, availableQty)
	} else {
		responseChan <- fmt.Sprintf("%s: %s is out of stock!", req.BuyerName, req.Product)
	}
}

func Buyer(req PurchaseRequest, storeChan chan PurchaseRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	storeChan <- req
}

func (s *Store) storeWorker(responseChan chan string) {
	for req := range s.Channel {
		s.ProcessPurchase(req, responseChan)
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *Store) Dispose() {
	close(s.Channel)
}

func main() {
	responseChan := make(chan string)

	store := NewStore(responseChan,
		ProductItem{"Apple", 6},
		ProductItem{"Banana", 6},
		ProductItem{"Orange", 3},
	)

	buyers := []PurchaseRequest{
		{"Ali", "Apple", 2},
		{"Mariem", "Apple", 4},
		{"Farida", "Banana", 3},
		{"Salim", "Orange", 5},
		{"Yahia", "Kiwi", 1},
		{"Noura", "Banana", 3},
		{"Hadi", "Banana", 1},
	}

	var wg sync.WaitGroup
	// Here’s my little listener goroutine. this is how I faced any confusion
	// It just listens to the response channel and prints out any message it receives.
	go func() {
		for res := range responseChan {
			fmt.Println(res)
		}
	}()

	for _, buyer := range buyers {
		wg.Add(1)
		fmt.Printf("%s is buying now.\n", buyer.BuyerName)
		go Buyer(buyer, store.Channel, &wg)
		fmt.Printf("%s has finished buying.\n", buyer.BuyerName)
	}

	wg.Wait()
	store.Dispose()
	time.Sleep(time.Second)
	close(responseChan)
}
