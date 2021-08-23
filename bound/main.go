package main

import (
	"time"
    "math/rand"
)

type (
    forkPair [2]chan struct{}
    forks chan forkPair
    philosopher struct{}
    philosophers chan *philosopher
)

func newForks(n int) forks {
    var (
	    out = make(forks, 5)
	    pair = make([]chan struct{}, 5, 5)
    )

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

func (p *philosopher) Eat(pair forkPair, forks forks, thinkers philosophers) {
	left := <-pair[0]
	right := <-pair[1]

	time.Sleep(time.Microsecond * 250)

	pair[0] <- left
	pair[1] <- right

	forks <- pair
	thinkers <- p
}

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

    for phil := range diners {
		go phil.Eat(<-forks, forks, thinkers)
	}
}
