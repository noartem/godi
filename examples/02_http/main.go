package main

import (
	"fmt"
	"log"
	"os"

	"github.com/noartem/godi"
)

func exec() error {
	c, err := godi.NewContainerWithLogger(
		log.New(os.Stdout, "", 0),
		NewApp, PortDefault, NewHTTP, NewHelloController, NewErrorController,
	)
	if err != nil {
		return fmt.Errorf("godi.NewContainer: %v", err)
	}

	appI, err := c.Get("IApp")
	if err != nil {
		return fmt.Errorf("container.Get IApp: %v", err)
	}

	app, ok := appI.(IApp)
	if !ok {
		return fmt.Errorf("invalid appI type: %T", appI)
	}

	err = app.Start()
	if err != nil {
		return fmt.Errorf("app.Start: %v", err)
	}

	return nil
}

func main() {
	if err := exec(); err != nil {
		panic(err)
	}
}
