package godi

import (
	"io/ioutil"
	"log"
)

// Container simple DI container
type Container struct {
	deps       map[string][]interface{}
	singletons map[string]interface{}
	log        *log.Logger
}

type DepType int

const (
	// Prototype a new instance every time dep is requested
	Prototype DepType = iota

	// Singleton only one instance of dep per Container
	Singleton
)

// DepOptions factory options
type DepOptions struct {
	Type DepType

	// Hooks?
}

// NewContainer create new DI container and register dependecies
func NewContainer(deps ...interface{}) (*Container, error) {
	logger := log.New(ioutil.Discard, "", 0)

	return NewContainerWithLogger(logger, deps...)
}

// NewContainerWithLogger create new DI container with custom logger and register dependecies
func NewContainerWithLogger(logger *log.Logger, deps ...interface{}) (*Container, error) {
	container := &Container{
		deps:       make(map[string][]interface{}),
		singletons: make(map[string]interface{}),
		log:        logger,
	}

	err := container.Register(deps...)
	if err != nil {
		return nil, err
	}

	return container, nil
}
