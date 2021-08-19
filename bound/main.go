package main

import (
	"time"
)

type forkPair [2]chan struct{}

func newForks(n int) chan forkPair {
	out := make(chan forkPair, 5)

	var pair = make([]chan struct{}, 5, 5)
	for i := 0; i < 5; i++ {
		pair[i] = make(chan struct{}, 1)
		pair[i] <- struct{}{}
	}

	for i := 0; i < 5; i++ {
		var f forkPair
		f[0] = pair[i]         //left fork
		f[1] = pair[(i-1+n)%n] //right fork
		out <- f
	}

	return out
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

func (p *philosopher) Eat(pair forkPair, forks chan forkPair, thinkers philosophers) {
	left := <-pair[0]
	right := <-pair[1]

	time.Sleep(time.Microsecond * 250)

	pair[0] <- left
	pair[1] <- right

	forks <- pair
	thinkers <- p
}

type philosophers chan *philosopher

func main() {
	var (
		numPhilosophers = 5
		forks           = newForks(numPhilosophers)           // buffered channel containing all empty seats and all unused fork pairs
		diners          = make(philosophers)                  // diners come as they become hungry and stop thinking
		thinkers        = make(philosophers, numPhilosophers) // buffered channel containing all philosophers in their starting state
	)

	for i := 0; i < numPhilosophers; i++ {
		thinkers <- newPhilosopher(diners, thinkers) // initialize our philosophers and start each goroutine
	}

	for {
		philosopher := <-diners
		pair := <-forks
		go philosopher.Eat(pair, forks, thinkers)
	}
}
