package main

import (
	"context"
	"fmt"
	"go-redis/client"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestServerWithClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()

	time.Sleep(time.Second)

	nClients := 5
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := client.New("localhost:5001")
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

	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Errorf("expected 0 peers, got %d", len(server.peers))
	}
}

func TestOfficialRedisClient(t *testing.T) {
	listenAddress := ":5001"
	server := NewServer(Config{
		ListenAddress: listenAddress,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Millisecond * 400)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", ":5001"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	testCases := map[string]string{
		"key_1": "val_1",
		"key_2": "val_2",
		"key_3": "val_3",
	}
	for key, val := range testCases {
		if err := rdb.Set(context.Background(), key, val, 0).Err(); err != nil {
			t.Fatal(err)
		}
		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}
		if newVal != val {
			t.Fatalf("expected %s but got %s", val, newVal)
		}
	}
}
