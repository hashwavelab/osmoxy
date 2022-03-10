package broadcast

import (
	"sync"
)

type broadcaster struct {
	sync.RWMutex
	on   bool
	subs map[chan<- interface{}]chan interface{}
}

type Broadcaster interface {
	Register(chan<- interface{}, int) bool
	Unregister(chan<- interface{})
	Close() error
	Submit(interface{})
	SubmitWait(interface{})
}

func NewBroadcaster() Broadcaster {
	b := &broadcaster{
		on:   true,
		subs: make(map[chan<- interface{}]chan interface{}),
	}
	return b
}

func (b *broadcaster) Register(newch chan<- interface{}, buflen int) bool {
	b.Lock()
	defer b.Unlock()
	if !b.on {
		return false
	}
	relayCh := make(chan interface{}, buflen)
	go func() {
		for m := range relayCh {
			newch <- m
		}
	}()
	b.subs[newch] = relayCh
	return true
}

func (b *broadcaster) Unregister(newch chan<- interface{}) {
	b.Lock()
	defer b.Unlock()
	delete(b.subs, newch)
}

func (b *broadcaster) Close() error {
	b.Lock()
	defer b.Unlock()
	b.on = false
	return nil
}

// Submit an item to be broadcast to all listeners.
// Message will be missed by some of listeners if the buffer is full.
func (b *broadcaster) Submit(m interface{}) {
	b.RLock()
	defer b.RUnlock()
	for _, relayCh := range b.subs {
		select {
		case relayCh <- m:
		default:
		}
	}
}

// Submit an item to be broadcast to all listeners.
// Wait for the message to be taken by all buffers
func (b *broadcaster) SubmitWait(m interface{}) {
	relayChs := make([]chan interface{}, 0)
	b.RLock()
	for _, relayCh := range b.subs {
		relayChs = append(relayChs, relayCh)
	}
	b.RUnlock()
	for _, relayCh := range relayChs {
		relayCh <- m
	}
}
