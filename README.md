# Store

This is a simple concurrent shopping store simulation built in Go.

# About the Project

The program simulates multiple buyers trying to purchase items from a store concurrently using Go's goroutines and channels. It ensures safe access to shared resources using mutexes.

# Features

- Buyers send purchase requests concurrently.
- Store handles requests safely using a mutex.
- Proper responses based on product availability.
- Use of goroutines, channels, and wait groups.
  
# Concepts Used 

- Golang (Go)
- Concurrency (goroutines, channels, mutex)
- Structs
  
# How to Run

1. Install Go: https://go.dev/dl/
2. Clone this repo:
   ```bash
   git clone https://github.com/MariemEssam/Store.git
   cd Store-go
   go run main.go


