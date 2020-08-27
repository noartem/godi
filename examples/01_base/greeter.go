package main

import (
	"fmt"

	"github.com/noartem/godi"
)

// IGreeter greeter service interface
type IGreeter interface {
	Greet() string
}

// Greeter service implementaion
type Greeter struct {
	name IName
}

// NewGreeter is a IGreeter factory
func NewGreeter(name IName, password IPassword) (IGreeter, godi.BeanOptions) {
	fmt.Println("(Password: \"" + password + "\")")

	greeter := &Greeter{
		name: name,
	}

	options := godi.BeanOptions{
		Type: godi.Singleton,
	}

	return greeter, options
}

// Greet generate greeting
func (h *Greeter) Greet() string {
	return "Hello, " + h.name.Gen() + "!"
}
