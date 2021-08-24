package main

import (
	"time"
    "math/rand"
)

type (
    fork chan struct{}
    philosopher struct{
        left fork
        right fork
    }
    philosophers chan *philosopher
)

func newPhilosopher(diners, thinkers philosophers, left, right fork) *philosopher {
	var p = &philosopher{ left, right }
	go func() {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))
        for phil := range thinkers {
            time.Sleep(time.Microsecond * time.Duration(r.Int63()))
            diners <- phil
		}
	}()
	return p
}

func (p *philosopher) Eat(thinkers philosophers) {
	left := <-p.left
	right := <-p.right

	time.Sleep(time.Microsecond * 250)

	p.left <- left
	p.right <- right

    thinkers <- p
}

func main() {
	var (
		numPhilosophers = 5
		diners          = make(philosophers)                     // diners come as they become hungry and stop thinking
		forks           = make([]fork, numPhilosophers)          // buffered channel containing all empty seats and all unused fork pairs
		thinkers        = make(philosophers, numPhilosophers)    // buffered channel containing all philosophers in their starting state
	)

	for i := 0; i < 5; i++ {
		forks[i] = make(fork, 1)
		forks[i] <- struct{}{}
	}

	for i := 0; i < numPhilosophers; i++ {
		thinkers <- newPhilosopher(diners, thinkers, forks[i], forks[(i-1+numPhilosophers)%numPhilosophers]) // initialize our philosophers and start each goroutine
	}

    for p := range diners {
		go p.Eat(thinkers)
	}
}
