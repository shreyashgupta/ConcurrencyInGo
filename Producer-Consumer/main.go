package main

import (
	"time"
)

type ProductEnum int

const (
	pizza   ProductEnum = 0
	pasta   ProductEnum = 1
	soup    ProductEnum = 2
	noodles ProductEnum = 3
)

func main() {
	channels := []chan Product{
		make(chan Product, 10),
		make(chan Product, 10),
		make(chan Product, 10),
		make(chan Product, 10),
	}

	producerA := Producer{
		Id:       1,
		Name:     "House of Italy",
		Products: []ProductEnum{pizza, pasta},
	}

	producerB := Producer{
		Id:       1,
		Name:     "China Town",
		Products: []ProductEnum{noodles, soup},
	}
	consumerA := Consumer{
		Id:           1,
		Name:         "shreyash",
		Requirements: []ProductEnum{pizza, soup},
	}

	consumerB := Consumer{
		Id:           2,
		Name:         "rohit",
		Requirements: []ProductEnum{pasta, noodles},
	}

	consumerC := Consumer{
		Id:           3,
		Name:         "dzeto",
		Requirements: []ProductEnum{pasta, soup},
	}

	go producerA.Produce(channels)
	go producerB.Produce(channels)

	time.Sleep(time.Second * 5)

	go consumerA.Consume(channels)
	go consumerB.Consume(channels)
	go consumerC.Consume(channels)

	for {
		time.Sleep(time.Second)
	}
}
