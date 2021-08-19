package main

import (
	"time"
)

type forks chan struct{}

func newForks() forks {
	var forks = make(forks, 5)
	for i := 0; i < 5; i++ {
		forks <- struct{}{}
	}
	return forks
}

type philosopher struct{}

func newPhilosopher(diners, thinkers philosophers) *philosopher {
	var p = &philosopher{}

	go func() {
		for {
			diners <- <-thinkers
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

type philosophers chan *philosopher

func main() {
	var (
		numPhilosophers = 5
		forks           = newForks()                          // buffered channel containing all empty seats and all unused chopsticks
		diners          = make(philosophers)                  // diners come as they become hungry and stop thinking
		thinkers        = make(philosophers, numPhilosophers) // buffered channel containing all philosophers in their starting state
	)

	for i := 0; i < numPhilosophers; i++ {
		thinkers <- newPhilosopher(diners, thinkers) // initialize our philosophers and start each goroutine
	}

	for {
		philosopher := <-diners
		go philosopher.Eat(forks, thinkers)
	}
}
