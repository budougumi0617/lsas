package main

import (
	"flag"

	"github.com/budougumi0617/lsas"
)

var (
	regionFlag      string
	printHeaderFlag bool
)

func main() {
	flag.StringVar(&regionFlag, "region", "", "AWS region")
	flag.StringVar(&regionFlag, "r", "", "AWS region")
	flag.BoolVar(&printHeaderFlag, "print", false, "print result header")
	flag.BoolVar(&printHeaderFlag, "p", false, "print result header")
	flag.Parse()

	if err := lsas.Execute(regionFlag, printHeaderFlag); err != nil {
		panic(err)
	}
}
