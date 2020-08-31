package main

import (
	"fmt"

	"github.com/noartem/godi"
)

// IGreeter greeter service interface
type IGreeter interface {
	Greet() string
}

// Greeter service implementation
type Greeter struct {
	deps deps
}

type deps struct {
	godi.InStruct

	Name IName
	Pass IPassword
}

// NewGreeter is a IGreeter factory
func NewGreeter(deps deps) (IGreeter, *godi.BeanOptions) {
	fmt.Println("(Password: \"" + deps.Pass + "\")")

	greeter := &Greeter{deps}

	options := &godi.BeanOptions{
		Type: godi.Singleton,
	}

	return greeter, options
}

// Greet generate greeting
func (h *Greeter) Greet() string {
	return "Hello, " + h.deps.Name.Gen() + "!"
}
