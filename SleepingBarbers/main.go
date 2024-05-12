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
	name        string
	barberColor *color.Color
}

func CreateCustomer(name string, hairLength int) *Customer {
	return &Customer{
		name:       name,
		hairLength: hairLength,
	}
}

func CreateBarber(name string, color *color.Color) *Barber {
	return &Barber{
		name:        name,
		barberColor: color,
	}
}

var chairs chan *Customer

func CreateAndSendCustomers(leftCount *int) {

	for {
		customer := CreateCustomer(fake.FirstName(), 2+rand.Intn(5))

		select {
		case chairs <- customer:
		default:
			*leftCount = *leftCount + 1
			fmt.Printf("No space available, customer %s is leaving\n", customer.name)
		}
		time.Sleep(time.Millisecond * 500)
	}

}

func (barber *Barber) Serve() {
	barber.barberColor.Printf("%s is available to attend customer!\n", barber.name)

	for {
		customer := <-chairs
		barber.barberColor.Printf("Attending customer %s \n", customer.name)
		time.Sleep(time.Second * time.Duration(customer.hairLength))
		barber.barberColor.Printf("Customer %s is all groomed\n", customer.name)
	}
}

func main() {
	var leftCount = 0
	chairs = make(chan *Customer, 10)

	barber0 := CreateBarber(fake.FirstName(), color.New(color.FgGreen))
	barber1 := CreateBarber(fake.FirstName(), color.New(color.FgBlue))
	barber2 := CreateBarber(fake.FirstName(), color.New(color.FgCyan))

	go barber0.Serve()
	go barber1.Serve()
	go barber2.Serve()

	go CreateAndSendCustomers(&leftCount)

	for {
		time.Sleep(10 * time.Second)
		colorPrint := color.New(color.FgHiRed).PrintfFunc()
		colorPrint("%d customer left\n", leftCount)
	}
}
