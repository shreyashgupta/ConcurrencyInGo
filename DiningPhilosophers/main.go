package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type Philosopher struct {
	id        int
	prodColor color.Attribute
}

var forks [5]chan bool

func PutDownFork(id int) {
	<-forks[id]
}

func PickUpFork(id int) {
	forks[id] <- true
}

func CreatePhilosopher(id int, color color.Attribute) *Philosopher {
	return &Philosopher{
		id:        id,
		prodColor: color,
	}
}

func (philosopher *Philosopher) Eat() {
	colorPrint := color.New(philosopher.prodColor).PrintfFunc()
	colorPrint("Philosopher %d is eating\n", philosopher.id)
	eatTime := 1 + rand.Intn(1)
	time.Sleep(time.Second * time.Duration(eatTime))
	colorPrint("Philosopher %d is done eating\n", philosopher.id)
}

func (philosopher *Philosopher) PutDownForks() {
	left := philosopher.id
	right := (philosopher.id + 1) % 5
	PutDownFork(left)
	PutDownFork(right)
	colorPrint := color.New(philosopher.prodColor).PrintfFunc()
	colorPrint("Philosopher %d has put down forks %d and %d\n", philosopher.id, left, right)
}

func (philosopher *Philosopher) Think() {
	colorPrint := color.New(philosopher.prodColor).PrintfFunc()
	colorPrint("Philosopher %d is thinking\n", philosopher.id)
	thinkTime := 1 + rand.Intn(2)
	time.Sleep(time.Second * time.Duration(thinkTime))
}

func (philosopher *Philosopher) Hungry() {
	for {
		colorPrint := color.New(philosopher.prodColor).PrintfFunc()
		colorPrint("Philosopher %d is hungry\n", philosopher.id)
		left := philosopher.id
		right := (philosopher.id + 1) % 5

		// To mitigate circular deadlock
		if left > right {
			left, right = right, left
		}

		PickUpFork(left)
		colorPrint("Philosopher %d picked up fork %d\n", philosopher.id, left)
		PickUpFork(right)
		colorPrint("Philosopher %d picked up fork %d\n", philosopher.id, right)
		philosopher.Eat()
		philosopher.PutDownForks()
		philosopher.Think()
	}
}

func main() {

	fmt.Println("Stating the feast")

	philosopher0 := CreatePhilosopher(0, color.FgBlue)

	philosopher1 := CreatePhilosopher(1, color.FgGreen)

	philosopher2 := CreatePhilosopher(2, color.FgMagenta)

	philosopher3 := CreatePhilosopher(3, color.FgRed)

	philosopher4 := CreatePhilosopher(4, color.FgYellow)

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

	x := make(chan bool)
	x <- true
}
