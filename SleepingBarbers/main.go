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

func CreateAndSendCustomers(chairs chan Customer, leftCount *int) {

	for {
		customer := Customer{
			name:       fake.FirstName(),
			hairLength: 2 + rand.Intn(5),
		}

		select {
		case chairs <- customer:
		default:
			*leftCount = *leftCount + 1
			fmt.Printf("No space available, customer %s is leaving\n", customer.name)
		}
		time.Sleep(time.Millisecond * 500)
	}

}

func ServeCustomers(chairs chan Customer, barberColor color.Attribute) {
	colorPrint := color.New(barberColor).PrintfFunc()
	for {
		customer := <-chairs
		colorPrint("Attending customer %s \n", customer.name)
		time.Sleep(time.Second * time.Duration(customer.hairLength))
		colorPrint("Customer %s is all groomed\n", customer.name)
	}
}

func main() {
	var leftCount = 0
	chairs := make(chan Customer, 10)

	go ServeCustomers(chairs, color.FgBlue)
	go ServeCustomers(chairs, color.FgGreen)
	go ServeCustomers(chairs, color.FgCyan)
	go ServeCustomers(chairs, color.FgHiMagenta)

	go CreateAndSendCustomers(chairs, &leftCount)

	for {
		time.Sleep(10 * time.Second)
		colorPrint := color.New(color.FgHiRed).PrintfFunc()
		colorPrint("%d customer left\n", leftCount)
	}
}
