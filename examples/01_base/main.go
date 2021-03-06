package main

import (
	"fmt"
	"github.com/noartem/godi"
)

func main() {
	// Crate DI container and register factories
	c, err := godi.NewContainer(
		NewGreeter,      // IGreeter
		NewRandom,       // IRandom
		NewName,         // IName
		PasswordDefault, // IPassword
		PasswordTest,    // IPassword for tests
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

	for i := 0; i < 10; i++ {
		fmt.Println(greeter.Greet())
	}

	// will don't call NewGreeter and return greeter from cache,
	// because Greeter is singleton
	_, err = c.Get("IGreeter")
	if err != nil {
		panic(err)
	}
}
