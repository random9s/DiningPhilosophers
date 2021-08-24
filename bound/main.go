package main

import (
	"math/rand"
	"time"
)

type (
	fork        chan struct{}
	philosopher struct {
		left  fork
		right fork
	}
	philosophers chan *philosopher
)

func (p *philosopher) Eat(r *rand.Rand, diners philosophers) {
	left := <-p.left
	right := <-p.right

	time.Sleep(time.Microsecond * 250) //spend some fixed time eating

	p.left <- left
	p.right <- right

    time.Sleep(time.Microsecond * time.Duration(r.Int63())) //spend some time thinking
	diners <- p
}

func main() {
	var (
		numPhilosophers = 5
		diners          = make(philosophers) // diners come as they become hungry and stop thinking
		forks           = []fork{make(fork, 1), make(fork, 1), make(fork, 1), make(fork, 1), make(fork, 1)} // buffered channel containing all empty seats and all unused fork pairs
        r = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

    for i := 0; i < numPhilosophers; i++ {
		forks[i] <- struct{}{}
    }

	for i := 0; i < numPhilosophers; i++ {
        time.Sleep(time.Microsecond * time.Duration(r.Int63())) //spend some time thinking
		diners <- &philosopher{forks[i], forks[(i-1+numPhilosophers)%numPhilosophers]} // initialize our philosophers and start each goroutine
	}

	for p := range diners {
		go p.Eat(r, diners)
	}
}
