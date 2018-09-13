package main
/*
rate limit 控制client对server的访问频率
golang.org/x/time/rate 提供了rate limit服务
*/


import (
	"sync"
	"log"
	"os"
	"context"
	"golang.org/x/time/rate"
)
func Open() *APIConnection {
	return &APIConnection{}
}
type APIConnection struct {
	rateLimiter * rate.Limiter
}
func (a *APIConnection) ReadFile(ctx context.Context) error {
	// Pretend we do work here
	if err := a.


	return nil
}
func (a *APIConnection) ResolveAddress(ctx context.Context) error {
	// Pretend we do work here
	return nil
}


func main() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	wg.Wait()
}
