package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewClients(t *testing.T) {
	nClients := 5
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()

			key := fmt.Sprintf("client_foo_%d", i)
			val := fmt.Sprintf("client_bar_%d", i)
			if err := c.Set(context.Background(), key, val); err != nil {
				log.Fatal(err)
			}
			value, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %d got this val back => %s\n", it, value)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestNewClient(t *testing.T) {
	c, err := New("localhost:5001")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
	for i := 0; i < 10; i++ {
		fmt.Println("SET =>", fmt.Sprintf("bar%d", i))
		if err := c.Set(context.Background(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}

		val, err := c.Get(context.Background(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", val)
	}
}
