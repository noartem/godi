package godi

import (
	"fmt"
	"reflect"
)

// ErrorType name of type error
const ErrorType = "error"

// DepOptionsType name of type DepOptions
const DepOptionsType = "godi.DepOptions"

// Register add dependencies fatories to DI container
func (container *Container) Register(deps ...interface{}) error {
	for _, dep := range deps {
		err := container.RegisterOne(dep)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterOne add dependency fatory to DI container
func (container *Container) RegisterOne(dep interface{}) error {
	container.log.Printf("Register: %v", dep)

	depType := reflect.TypeOf(dep)

	err := checkDepOut(depType)
	if err != nil {
		return err
	}

	name := depType.Out(0).Name()
	if container.deps[name] == nil {
		container.deps[name] = []interface{}{}
	}

	container.deps[name] = append(container.deps[name], dep)

	return nil
}

func checkDepOut(depType reflect.Type) error {
	if depType.Kind() != reflect.Func {
		return fmt.Errorf("invalid dependency type: %s", depType)
	}

	numOut := depType.NumOut()
	if numOut == 0 || numOut > 3 {
		return fmt.Errorf("invalid dependency NumOut excepted: [1, 3], got: %d", depType.NumOut())
	}

	if numOut == 1 {
		return nil
	}

	var hasError, hasOpts bool

	switch depType.Out(1).Name() {
	case ErrorType:
		hasError = true
	case DepOptionsType:
		hasOpts = true
	default:
		return fmt.Errorf("invalid second out")
	}

	if numOut == 2 {
		return nil
	}

	out2Name := depType.Out(2).Name()

	if out2Name == ErrorType && hasError {
		return fmt.Errorf("has two error out")
	}

	if out2Name == DepOptionsType && hasOpts {
		return fmt.Errorf("has two options out")
	}

	return nil
}
