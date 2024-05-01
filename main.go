package main

import (
	"flag"
	"fmt"
)

func main() {
	day := flag.String("day", "today", "What day's weather do you want to see?")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Printf("Here is the weather for %v!\n", *day)
	} else if flag.Arg(0) == "day" {
		fmt.Printf("Here is the weather for %v!\n", *day)
	}
}
