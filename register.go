package godi

import (
	"fmt"
	"reflect"
)

// ErrorInterface interface of error
var ErrorInterface = reflect.TypeOf((*error)(nil)).Elem()

// BeanOptionsType reflect type of BeanOptions
var BeanOptionsType = reflect.TypeOf(&BeanOptions{})

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

	factoryName := genFactoryName(factoryType)
	factories := container.factories[factoryName]
	if factories == nil {
		factories = []interface{}{}
	}

	container.factories[factoryName] = append(factories, factory)

	container.log.Printf("Registered: %s = %v", factoryName, factory)

	return nil
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

	out1 := factoryType.Out(1)
	if out1.Implements(ErrorInterface) {
		hasError = true
	} else if out1 == BeanOptionsType {
		hasOpts = true
	} else {
		return fmt.Errorf("invalid second out: %s", factoryType.Out(1).String())
	}

	if numOut == 2 {
		return nil
	}

	out2 := factoryType.Out(2)

	if out2.Implements(ErrorInterface) && hasError {
		return fmt.Errorf("has two error out")
	}

	if out2 == BeanOptionsType && hasOpts {
		return fmt.Errorf("has two bean options out")
	}

	return nil
}
