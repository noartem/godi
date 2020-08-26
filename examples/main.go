package main

import (
	"fmt"

	"github.com/noartem/godi"
)

func main() {
	// Crate DI container and register factories
	c, err := godi.NewContainer(
		NewGreeter, // IGreeter
		NewRandom,  // IRandom
		NewName,    // IName
	)
	if err != nil {
		panic(err)
	}

	// Get bean from container
	greeterI, err := c.Get("IGreeter")
	if err != nil {
		panic(err)
	}

	greeter, ok := greeterI.(IGreeter)
	if !ok {
		panic("Invalid bean")
	}

	fmt.Println(greeter.Greet())
}
