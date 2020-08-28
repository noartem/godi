package godi

import (
	"io/ioutil"
	"log"
	"os"
)

// Container simple DI container
type Container struct {
	factories      map[string][]interface{}
	beanSingletons map[string]interface{}
	log            *log.Logger
}

// BeanType type of bean
type BeanType int

const (
	// Prototype a new instance every time bean is requested
	Prototype BeanType = iota

	// Singleton only one instance of bean per Container
	Singleton
)

// BeanOptions factory options
type BeanOptions struct {
	Type BeanType

	// Hooks?
}

// NewContainer create new DI container and register dependencies
func NewContainer(factories ...interface{}) (*Container, error) {
	logger := log.New(ioutil.Discard, "", 0)

	return NewContainerWithLogger(logger, factories...)
}

// NewContainerWithLogging create new DI container with default logger
func NewContainerWithLogging(factories ...interface{}) (*Container, error) {
	logger := log.New(os.Stdout, "GODI:", 0)

	return NewContainerWithLogger(logger, factories...)
}

// NewContainerWithLogger create new DI container with custom logger and register dependencies
func NewContainerWithLogger(logger *log.Logger, factories ...interface{}) (*Container, error) {
	container := &Container{
		factories:      make(map[string][]interface{}),
		beanSingletons: make(map[string]interface{}),
		log:            logger,
	}

	err := container.Register(factories...)
	if err != nil {
		return nil, err
	}

	return container, nil
}
