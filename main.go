package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"rout/router"
)

var (
	route *string
	subs  *string
)

func main() {

	route = flag.String("r", "", "route")
	subs = flag.String("t", "", "topic")
	flag.Parse()

	if *subs == "" || *route == "" {
		fmt.Println(errors.New("Invalid route or topic. "))
		os.Exit(1)
	}

	if p, e := router.NewRoutePattern("b", *route); e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	} else {
		if p.Match(*subs) {
			fmt.Println(true)
			os.Exit(0)
		}
		fmt.Println(false)

	}
	os.Exit(1)
}
