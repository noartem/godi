# Godi

Simple Golang DI container based on reflection

Example:
```go
package main

import (
  "fmt"
  "github.com/noartem/godi"
)

// Name service interface
type IName interface {
  Gen() string
}

// Name service implementaion
type Name struct {}

// Name serivice implementation factory
func NewName(rnd IRandom) IName {
  return &Name{}
}

func (n *Name) Gen() string {
  return "World"
}

// Greeter service interface
type IGreeter interface {
  Greet() string
}

// Greeter service implementaion
type Greeter struct {
  name IName
}

// Greeter service factory
func NewGreeter(name IName) IHello {
  return &Hello{
    name: name,
  }
}

func (h *Hello) Greet() string {
  return fmt.Sprintf("Hello, %s!", h.name.Gen())
}

func main() {
  // Crate DI container and register services factories
  c, err := godi.NewContainer(
      NewGreeter, // will be registered as IGreeter
      NewName, // will be registered as IName
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

  fmt.Println(greeter.Greete())
}
```

