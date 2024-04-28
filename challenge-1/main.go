package main

import (
	"fmt"
	"sync"
)

var msg string

func updateMessage(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	msg = s
}

func printMessage() {
	fmt.Println(msg)
}

func main() {

	// challenge: modify this code so that the calls to updateMessage() on lines
	// 27, 30, and 33 run as goroutines, and implement wait groups so that
	// the program runs properly, and prints out three different messages.
	// Then, write a test for all three functions in this program: updateMessage(),
	// printMessage(), and main().

	var wg sync.WaitGroup
	wg.Add(1)

	msg = "Hello, world!"

	go updateMessage("Hello, universe!", &wg)
	wg.Wait()

	printMessage()

	wg.Add(1)

	go updateMessage("Hello, cosmos!", &wg)
	wg.Wait()

	printMessage()

	wg.Add(1)

	go updateMessage("Hello, world!", &wg)
	wg.Wait()

	printMessage()
}
