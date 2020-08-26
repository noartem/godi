package main

import (
	"fmt"

	"github.com/noartem/godi"
)

func main() {
	// Crate DI container and register services factories
	c, err := godi.NewContainer(
		NewGreeter, // will be registered as IGreeter
		NewRandom,  // will be registered as IRandom
		NewName,    // will be registered as IName
	)
	if err != nil {
		panic(err)
	}

	// Get generated service from container
	greeterI, err := c.Get("IGreeter")
	if err != nil {
		panic(err)
	}

	greeter, ok := greeterI.(IGreeter)
	if !ok {
		panic("Invalid interface")
	}

	fmt.Println(greeter.Greet())
}
