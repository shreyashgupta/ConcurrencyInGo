package main

import "github.com/fatih/color"

type Product struct {
	Id           int
	Name         string
	ProducerName string
	ProductColor color.Attribute
}

func (e ProductEnum) String() string {
	switch e {
	case pizza:
		return "pizza"
	case pasta:
		return "pasta"
	case noodles:
		return "noodles"
	case soup:
		return "soup"
	default:
		return ""
	}
}

func (e ProductEnum) Color() color.Attribute {
	switch e {
	case pizza:
		return color.FgRed
	case pasta:
		return color.FgWhite
	case noodles:
		return color.FgYellow
	case soup:
		return color.FgCyan
	default:
		return color.BgHiWhite
	}
}
