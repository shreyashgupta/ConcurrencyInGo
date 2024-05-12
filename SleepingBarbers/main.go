package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/icrowley/fake"
)

type Customer struct {
	name       string
	hairLength int
}

type Barber struct {
	name  string
	skill int
}

type Salon struct {
	numBarbers     int
	chairs         chan *Customer
	closedChan     chan bool
	barberDoneChan chan bool
	shopOpen       bool
	unservedCount  int
	openDuration   time.Duration
}

func createSalon(capacity, numBarbers, openDurationSec int) *Salon {
	return &Salon{
		numBarbers:     numBarbers,
		chairs:         make(chan *Customer, capacity),
		closedChan:     make(chan bool),
		barberDoneChan: make(chan bool),
		shopOpen:       false,
		unservedCount:  0,
		openDuration:   time.Second * time.Duration(openDurationSec),
	}
}

func createCustomer(name string, hairLength int) *Customer {
	return &Customer{
		name,
		hairLength,
	}
}

func createBarber(name string, skill int) *Barber {
	return &Barber{
		name,
		skill,
	}
}

func (salon *Salon) open() {
	// Open salon and make barbers available

	salon.shopOpen = true
	color.Green("Salon is open for business !")

	for i := 0; i < salon.numBarbers; i++ {
		barber := createBarber(fake.FirstName(), 1+rand.Intn(2))
		go barber.serve(salon)
	}
}

func (salon *Salon) close() {

	color.Yellow("Salon is closing for the day")

	// Stop sending new customers
	salon.closedChan <- true
	close(salon.chairs)

	// Set the flag to false to make sure barbers leave for the day once no customer is waiting
	salon.shopOpen = false

	color.Magenta("Waiting for babers to finish up, still %d customers left\n", len(salon.chairs))
	for i := 0; i < salon.numBarbers; i++ {
		// Wait for all the barbers to leave
		<-salon.barberDoneChan
	}
	color.Yellow("Salon is closed for the day")
	color.Yellow("%d Customer were unserved", salon.unservedCount)

	// Close the channels
	close(salon.closedChan)
	close(salon.barberDoneChan)
}

func (barber *Barber) leaveForTheDay(salon *Salon) {
	color.Magenta("%s is done for the day", barber.name)
	// Signal that barber is done for the day
	salon.barberDoneChan <- true
}

func (barber *Barber) serve(salon *Salon) {
	fmt.Printf("%s is available to attend customer!\n", barber.name)

	for {
		// Get customer from chair and the salon open status
		customer, salonOpen := <-salon.chairs

		if salonOpen {
			// In case salon is open or customers are waiting, serve them
			fmt.Printf("%s is cutting hair of %s \n", barber.name, customer.name)
			time.Sleep(time.Second * time.Duration(customer.hairLength-barber.skill))
			fmt.Printf(" - Customer %s is all groomed\n", customer.name)
		} else {
			// In case salon is closed and no more customers are left to serve,
			// leave for the day
			barber.leaveForTheDay(salon)
			return
		}
	}
}

func createAndSendCustomers(salon *Salon) {
	// This function creates and sends customers to salon until the salon is open
	for {
		customer := createCustomer(fake.FirstName(), 2+rand.Intn(5))

		select {
		case <-salon.closedChan:
			// In case salon is closed, exit the function
			color.Cyan("Stopped sending new customers")
			return
		case salon.chairs <- customer:
		default:
			// If no space is available in salon, customer leaves
			salon.unservedCount++
			color.Red("No space available, customer %s is leaving\n", customer.name)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func signalEOD(salon *Salon) {
	// Close the salon after a specified duration
	<-time.After(salon.openDuration)
	salon.close()
}

func main() {

	salon := createSalon(10, 4, 20)
	salon.open()

	go createAndSendCustomers(salon)

	signalEOD(salon)
}
