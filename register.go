package godi

import (
	"fmt"
	"reflect"
)

// ErrorType name of type error
const ErrorType = "error"

// BeanOptionsType name of type BeanOptions
const BeanOptionsType = "BeanOptions"

// Registrar function registering factories
type Registrar func(c *Container) error

// RegisterCompose run all registrars
func (container *Container) RegisterCompose(registrars ...Registrar) error {
	for _, registrar := range registrars {
		err := registrar(container)
		if err != nil {
			return err
		}
	}

	return nil
}

// Register add beans factories to DI container
func (container *Container) Register(factories ...interface{}) error {
	for _, factory := range factories {
		err := container.RegisterOne(factory)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterOne add bean factory to DI container
func (container *Container) RegisterOne(factory interface{}) error {
	container.log.Printf("Register: %v", factory)

	factoryType := reflect.TypeOf(factory)

	err := checkFactoryOut(factoryType)
	if err != nil {
		return err
	}

	name := getFactoryName(factoryType)
	if container.factories[name] == nil {
		container.factories[name] = []interface{}{}
	}

	container.factories[name] = append(container.factories[name], factory)

	container.log.Printf("Registered: %s = %v", name, factory)

	return nil
}

func getFactoryName(factoryType reflect.Type) string {
	if factoryType.Kind() == reflect.Func {
		return factoryType.Out(0).Name()
	}

	return factoryType.Name()
}

func checkFactoryOut(factoryType reflect.Type) error {
	if factoryType.Kind() != reflect.Func {
		// if factory is already bean don't check
		return nil
	}

	numOut := factoryType.NumOut()
	if numOut == 0 || numOut > 3 {
		return fmt.Errorf("invalid factory NumOut excepted: [1, 3], got: %d", factoryType.NumOut())
	}

	if numOut == 1 {
		return nil
	}

	var hasError, hasOpts bool

	switch factoryType.Out(1).Name() {
	case ErrorType:
		hasError = true
	case BeanOptionsType:
		hasOpts = true
	default:
		return fmt.Errorf("invalid second out: %s", factoryType.Out(1).Name())
	}

	if numOut == 2 {
		return nil
	}

	out2Name := factoryType.Out(2).Name()

	if out2Name == ErrorType && hasError {
		return fmt.Errorf("has two error out")
	}

	if out2Name == BeanOptionsType && hasOpts {
		return fmt.Errorf("has two bean options out")
	}

	return nil
}
