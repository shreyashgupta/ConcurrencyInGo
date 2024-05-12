package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type Philosopher struct {
	id        int
	prodColor *color.Color
}

var forks [5]chan bool

func putDownFork(id int) {
	// Empty the fork channel
	<-forks[id]
}

func pickUpFork(id int) {
	// Try to insert into the fork channel,
	// this will block until channel is empty
	// i.e no other philosopher has a hold on it
	forks[id] <- true
}

func createPhilosopher(id int, color *color.Color) *Philosopher {
	return &Philosopher{
		id:        id,
		prodColor: color,
	}
}

func (philosopher *Philosopher) eat() {
	philosopher.prodColor.Printf("Philosopher %d is eating\n", philosopher.id)
	eatTime := 1 + rand.Intn(1)
	time.Sleep(time.Second * time.Duration(eatTime))
	philosopher.prodColor.Printf("Philosopher %d is done eating\n", philosopher.id)
}

func (philosopher *Philosopher) think() {
	philosopher.prodColor.Printf("Philosopher %d is thinking\n", philosopher.id)
	thinkTime := 1 + rand.Intn(2)
	time.Sleep(time.Second * time.Duration(thinkTime))
}

func (philosopher *Philosopher) putDownForks(left, right int) {
	// This function empties the left and right fork channel
	putDownFork(left)
	putDownFork(right)
	philosopher.prodColor.Printf("Philosopher %d has put down forks %d and %d\n", philosopher.id, left, right)
}

func (philosopher *Philosopher) pickUpForks(left, right int) {

	// Philosopher tries to pick up the left fork
	pickUpFork(left)
	philosopher.prodColor.Printf("Philosopher %d picked up fork %d\n", philosopher.id, left)

	// Philosopher tries to pick up the right fork
	pickUpFork(right)
	philosopher.prodColor.Printf("Philosopher %d picked up fork %d\n", philosopher.id, right)
}

func (philosopher *Philosopher) Hungry() {
	for {
		philosopher.prodColor.Printf("Philosopher %d is hungry\n", philosopher.id)
		left := philosopher.id
		right := (philosopher.id + 1) % 5

		// To mitigate circular deadlock
		if left > right {
			left, right = right, left
		}

		// Philospher tries to pick up forks
		philosopher.pickUpForks(left, right)

		// Now that both forks are with philosopher, he starts to eat
		philosopher.eat()

		// Once philospher is done eating, release both the forks
		philosopher.putDownForks(left, right)

		// After releasing both forks, philosopher goes back to thinking
		philosopher.think()
	}
}

func main() {

	fmt.Println("Stating the feast")

	philosopher0 := createPhilosopher(0, color.New(color.FgBlue))

	philosopher1 := createPhilosopher(1, color.New(color.FgGreen))

	philosopher2 := createPhilosopher(2, color.New(color.FgMagenta))

	philosopher3 := createPhilosopher(3, color.New(color.FgRed))

	philosopher4 := createPhilosopher(4, color.New(color.FgYellow))

	color.New(color.BgBlue)

	forks = [5]chan bool{
		make(chan bool, 1),
		make(chan bool, 1),
		make(chan bool, 1),
		make(chan bool, 1),
		make(chan bool, 1),
	}

	go philosopher0.Hungry()
	go philosopher1.Hungry()
	go philosopher2.Hungry()
	go philosopher3.Hungry()
	go philosopher4.Hungry()

	// Block the main function to know if other go routines
	// got into a deadlock
	x := make(chan bool)
	x <- true
}
