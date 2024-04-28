package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type Producer struct {
	Id       int
	Name     string
	Products []ProductEnum
}

func (p *Producer) Produce(channels []chan Product) {

	fmt.Printf("%s is online!\n", p.Name)

	for {
		product_id := p.Products[rand.Intn(len(p.Products))]
		productcolor := product_id.Color()
		colorPrint := color.New(productcolor).PrintfFunc()
		product_chan := channels[product_id]
		prep_time := 2 + rand.Intn(5)
		prodId := int(rand.Int31n(1000))
		time.Sleep(time.Duration(prep_time) * time.Second)
		colorPrint("Produced %s (%d) in time %d. Dispatching order\n", product_id.String(), prodId, prep_time)
		product := Product{
			Id:           prodId,
			Name:         product_id.String(),
			ProducerName: p.Name,
			ProductColor: productcolor,
		}
		product_chan <- product
	}
}
