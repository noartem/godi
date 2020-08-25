package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/noartem/godi"
)

type IRandom interface {
	RandomName() string
}

type Random struct{}

func (rnd *Random) RandomName() string {
	return randEl([]string{
		"Maxim",
		"Artem",
		"Andrey",
	})
}

func NewRandom() IRandom {
	return &Random{}
}

type IName interface {
	Name() string
}

type Name struct {
	rnd IRandom
}

func NewName(rnd IRandom) IName {
	return &Name{
		rnd: rnd,
	}
}

func (n *Name) Name() string {
	return "Noskov " + n.rnd.RandomName()
}

type IHello interface {
	Hello() string
}

type Hello struct {
	name IName
}

func NewHello(name IName) IHello {
	return &Hello{
		name: name,
	}
}

func (h *Hello) Hello() string {
	return fmt.Sprintf("Hello, %s!", h.name.Name())
}

func randEl(arr []string) string {
	return arr[rand.Intn(len(arr)-1)]
}

func main() {
	c, err := godi.NewContainerWithLogger(
		log.New(os.Stdout, "", 0),
		NewHello, NewRandom, NewName)
	if err != nil {
		panic(err)
	}

	hI, err := c.Get("IHello")
	if err != nil {
		panic(err)
	}

	h, ok := hI.(IHello)
	if !ok {
		panic("Invalid h")
	}

	fmt.Println(h.Hello())
}
