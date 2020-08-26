package main

// IGreeter greeter service interface
type IGreeter interface {
	Greet() string
}

// Greeter service implementaion
type Greeter struct {
	name IName
}

// NewGreeter is a IGreeter factory
func NewGreeter(name IName) IGreeter {
	return &Greeter{
		name: name,
	}
}

// Greet generate greeting
func (h *Greeter) Greet() string {
	return "Hello, " + h.name.Gen() + "!"
}
