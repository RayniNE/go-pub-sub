package pubSub

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

type PubSub[T any] struct {
	subscribers []chan T
	mu          sync.RWMutex
	closed      bool
}

func NewPubSub[T any]() *PubSub[T] {
	return &PubSub[T]{
		mu: sync.RWMutex{},
	}
}

func (ps *PubSub[T]) Subcribe() <-chan T {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return nil
	}

	r := make(chan T)

	ps.subscribers = append(ps.subscribers, r)

	return r
}

func (ps *PubSub[T]) Publish(value T) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return
	}

	for _, sub := range ps.subscribers {
		sub <- value
	}
}

func (ps *PubSub[T]) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return
	}

	for _, sub := range ps.subscribers {
		close(sub)
	}

	ps.closed = true
}

func AddSubscriber(pubSub *PubSub[string], id int) {
	sub := pubSub.Subcribe()
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			value, ok := <-sub
			if !ok {
				fmt.Printf("Subscriber #%d is exiting\n", id)
				return
			}

			fmt.Printf("Subscriber #%d, value: %v\n", id, value)
		}
	}()
}

func Run() {
	pubSub := NewPubSub[string]()

	for i := 1; i <= 5; i++ {
		AddSubscriber(pubSub, i)
	}

	pubSub.Publish("Game of Thrones")
	pubSub.Publish("Breaking Bad")
	pubSub.Publish("Prison Break")
	pubSub.Publish("Supernatural")

	pubSub.Close()

	wg.Wait()
}
