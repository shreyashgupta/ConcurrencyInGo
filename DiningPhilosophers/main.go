package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Philospher struct {
	id        int
	prodColor color.Attribute
}

type Fork struct {
	id        int
	available bool
	channel   chan bool
	lock      sync.Mutex
}

func (fork *Fork) pickup() {
	fork.lock.Lock()
	defer fork.lock.Unlock()
	fork.available = false
}

func (fork *Fork) isAvailable() bool {
	fork.lock.Lock()
	defer fork.lock.Unlock()
	return fork.available
}

func (fork *Fork) makeAvailable() {
	fork.lock.Lock()
	defer fork.lock.Unlock()
	fork.available = true
}

func (fork *Fork) informAvailability() {
	if len(fork.channel) == 0 {
		fork.channel <- true
	}
}

var forks [5]Fork

var forksLock sync.Mutex

func (philospher *Philospher) Eat() {
	colorPrint := color.New(philospher.prodColor).PrintfFunc()
	colorPrint("Philosopher %d is eating\n", philospher.id)
	eatTime := 2 + rand.Intn(5)
	time.Sleep(time.Second * time.Duration(eatTime))
	colorPrint("Philosopher %d is done eating\n", philospher.id)

	philospher.PutDownForks()
}

func (philospher *Philospher) PutDownForks() {
	colorPrint := color.New(philospher.prodColor).PrintfFunc()
	forks[philospher.id].makeAvailable()
	forks[(philospher.id+1)%5].makeAvailable()
	forks[philospher.id].informAvailability()
	forks[(philospher.id+1)%5].informAvailability()
	colorPrint("Philosopher %d has put down forks %d and %d\n", philospher.id, philospher.id, (philospher.id+1)%5)
}

func (philospher *Philospher) Hungry() {
	for {
		colorPrint := color.New(philospher.prodColor).PrintfFunc()
		colorPrint("Philosopher %d is hungry\n", philospher.id)
		left := philospher.id
		right := (philospher.id + 1) % 5
		select {
		case <-forks[left].channel:
			{
				forksLock.Lock()
				canEat := false
				if forks[left].isAvailable() && forks[right].isAvailable() {
					forks[left].pickup()
					forks[right].pickup()
					canEat = true
				} else {
					colorPrint("Philosopher %d is can't eat, going back to thinking\n", philospher.id)
				}
				forksLock.Unlock()
				if canEat {
					philospher.Eat()
					colorPrint("Philosopher %d is thinking\n", philospher.id)
				}
			}
		case <-forks[right].channel:
			{
				forksLock.Lock()
				canEat := false
				if forks[left].isAvailable() && forks[right].isAvailable() {
					forks[left].pickup()
					forks[right].pickup()
					canEat = true
				} else {
					colorPrint("Philosopher %d is can't eat, going back to thinking\n", philospher.id)
				}
				forksLock.Unlock()
				if canEat {
					philospher.Eat()
					colorPrint("Philosopher %d is thinking\n", philospher.id)
				}
			}

		}
		time.Sleep(5 * time.Second)
	}
}

func main() {

	fmt.Println("Stating the feast")

	Philospher0 := Philospher{
		id:        0,
		prodColor: color.FgBlue,
	}
	Philospher1 := Philospher{
		id:        1,
		prodColor: color.FgGreen,
	}
	Philospher2 := Philospher{
		id:        2,
		prodColor: color.FgHiMagenta,
	}
	Philospher3 := Philospher{
		id:        3,
		prodColor: color.FgRed,
	}
	Philospher4 := Philospher{
		id:        4,
		prodColor: color.FgHiYellow,
	}

	go Philospher0.Hungry()
	go Philospher1.Hungry()
	go Philospher2.Hungry()
	go Philospher3.Hungry()
	go Philospher4.Hungry()

	forks = [5]Fork{
		Fork{
			id:        0,
			available: true,
			channel:   make(chan bool, 1),
		},
		Fork{
			id:        1,
			available: true,
			channel:   make(chan bool, 1),
		},
		Fork{
			id:        2,
			available: true,
			channel:   make(chan bool, 1),
		},
		Fork{
			id:        3,
			available: true,
			channel:   make(chan bool, 1),
		},
		Fork{
			id:        4,
			available: true,
			channel:   make(chan bool, 1),
		},
	}

	go forks[0].informAvailability()
	go forks[1].informAvailability()
	go forks[2].informAvailability()
	go forks[3].informAvailability()
	go forks[4].informAvailability()
	x := make(chan bool)
	x <- true
}
