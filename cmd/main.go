package main

import (
	"flag"

	"github.com/budougumi0617/lsas"
)

var (
	regionFlag      string
	printHeaderFlag bool
	ignoreCaseFlag  bool
)

func main() {
	flag.StringVar(&regionFlag, "region", "", "AWS region")
	flag.StringVar(&regionFlag, "r", "", "AWS region")
	flag.BoolVar(&printHeaderFlag, "print", false, "print result header")
	flag.BoolVar(&printHeaderFlag, "p", false, "print result header")
	flag.BoolVar(&ignoreCaseFlag, "ignore-case", false, "Perform case insensitive matching. By default, grep is case sensitive.")
	flag.BoolVar(&ignoreCaseFlag, "i", false, "Perform case insensitive matching. By default, grep is case sensitive.")
	flag.Parse()

	if err := lsas.Execute(regionFlag, printHeaderFlag, ignoreCaseFlag); err != nil {
		panic(err)
	}
}
