package main

import (
	"fmt"

	"github.com/fatih/color"
)

type Consumer struct {
	Id           int
	Name         string
	Requirements []ProductEnum
}

func (c *Consumer) Consume(channels []chan Product) {

	fmt.Printf("%s is hungry!\n", c.Name)

	for {

		select {
		case prod := <-channels[c.Requirements[0]]:
			colorPrint := color.New(prod.ProductColor).PrintfFunc()
			colorPrint("%s recieved %s (%d) !\n", c.Name, prod.Name, prod.Id)
		case prod := <-channels[c.Requirements[1]]:
			colorPrint := color.New(prod.ProductColor).PrintfFunc()
			colorPrint("%s recieved %s (%d) !\n", c.Name, prod.Name, prod.Id)

		}
	}
}
