package main

import (
	"time"
    "math/rand"
)

type (
    forks chan struct{}
    philosopher struct{}
    philosophers chan *philosopher
)

func newPhilosopher(diners, thinkers philosophers) *philosopher {
	var p = &philosopher{}
	go func() {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))
        for phil := range thinkers {
            time.Sleep(time.Microsecond * time.Duration(r.Int63()))
            diners <- phil
		}
	}()
	return p
}

func (p *philosopher) Eat(forks forks, thinkers philosophers) {
	left := <-forks
	right := <-forks

	time.Sleep(time.Microsecond * 250)

	forks <- right
	forks <- left

	thinkers <- p
}

func main() {
	var (
		numPhilosophers = 5
		diners          = make(philosophers)                  // diners come as they become hungry and stop thinking
		thinkers        = make(philosophers, numPhilosophers) // buffered channel containing all philosophers in their starting state
		forks           = make(forks, numPhilosophers)        // buffered channel containing all empty seats and all unused chopsticks
	)

	for i := 0; i < numPhilosophers; i++ {
		thinkers <- newPhilosopher(diners, thinkers) // initialize our philosophers and start each goroutine
		forks <- struct{}{}                          // initialize our forks
	}

    for phil := range diners {
		go phil.Eat(forks, thinkers)
	}
}
