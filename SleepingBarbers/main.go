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
	open           bool
	unservedCount  int
	openDuration   time.Duration
}

func CreateSalon(capacity, numBarbers, openDurationSec int) *Salon {
	return &Salon{
		numBarbers:     numBarbers,
		chairs:         make(chan *Customer, capacity),
		closedChan:     make(chan bool),
		barberDoneChan: make(chan bool),
		open:           false,
		unservedCount:  0,
		openDuration:   time.Second * time.Duration(openDurationSec),
	}
}

func (salon *Salon) Open() {
	salon.open = true
	color.Green("Salon is open for business !")

	for i := 0; i < salon.numBarbers; i++ {
		barber := CreateBarber(fake.FirstName(), 1+rand.Intn(2))
		go barber.Serve(salon)
	}
}

func (salon *Salon) Close() {
	color.Yellow("Salon is closing for the day")
	salon.closedChan <- true
	close(salon.chairs)
	salon.open = false
	color.Magenta("Waiting for babers to finish up, still %d customers left\n", len(salon.chairs))
	for i := 0; i < salon.numBarbers; i++ {
		<-salon.barberDoneChan
	}
	color.Yellow("Salon is closed for the day")
	color.Yellow("%d Customer were unserved", salon.unservedCount)
	close(salon.closedChan)
	close(salon.barberDoneChan)
}

func CreateCustomer(name string, hairLength int) *Customer {
	return &Customer{
		name,
		hairLength,
	}
}

func CreateBarber(name string, skill int) *Barber {
	return &Barber{
		name,
		skill,
	}
}

func (barber *Barber) LeaveForTheDay(salon *Salon) {
	color.Magenta("%s is done for the day", barber.name)
	salon.barberDoneChan <- true
}

func (barber *Barber) Serve(salon *Salon) {
	fmt.Printf("%s is available to attend customer!\n", barber.name)

	for {
		customer, salonOpen := <-salon.chairs

		if salonOpen {
			fmt.Printf("%s is cutting hairs of %s \n", barber.name, customer.name)
			time.Sleep(time.Second * time.Duration(customer.hairLength))
			fmt.Printf(" - Customer %s is all groomed\n", customer.name)
		} else {
			barber.LeaveForTheDay(salon)
			return
		}
	}
}

func CreateAndSendCustomers(salon *Salon) {

	for {
		customer := CreateCustomer(fake.FirstName(), 2+rand.Intn(5))

		select {
		case <-salon.closedChan:
			color.Cyan("Stopped sending new customers")
			return
		case salon.chairs <- customer:
		default:
			salon.unservedCount++
			color.Red("No space available, customer %s is leaving\n", customer.name)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func SignalEOD(salon *Salon) {
	<-time.After(salon.openDuration)
	salon.Close()
}

func main() {

	salon := CreateSalon(10, 4, 20)
	salon.Open()

	go CreateAndSendCustomers(salon)

	SignalEOD(salon)
}
