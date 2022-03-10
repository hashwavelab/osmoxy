package broadcast

import (
	"log"
	"time"
)

// Example of a simple broadcaster sending numbers to two workers.
//
// Five messages are sent.  The first worker prints all five.  The second worker prints the first and then unsubscribes.
func main() {
	b := NewBroadcaster()

	workerOne(b)
	workerTwo(b)

	for i := 0; i < 100; i++ {
		log.Printf("Sending %v", i)
		b.Submit(i)
	}
	defer b.Close()
	time.Sleep(1000 * time.Second)
}

func workerOne(b Broadcaster) {
	ch := make(chan interface{})
	b.Register(ch, 100)
	//defer b.Unregister(ch)
	wait := false
	// Dump out each message sent to the broadcaster.
	go func() {
		for v := range ch {
			log.Printf("workerOne read %v", v)
			if !wait {
				time.Sleep(1 * time.Second)
				wait = true
			}
		}
	}()
}

func workerTwo(b Broadcaster) {
	ch := make(chan interface{})
	b.Register(ch, 100)
	// defer b.Unregister(ch)
	// defer log.Printf("workerTwo is done\n")

	go func() {
		for v := range ch {
			log.Printf("workerTwo read %v", v)
			time.Sleep(100 * time.Millisecond)
		}
	}()
}
