package godi

import (
	"fmt"
	"reflect"
)

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
		return fmt.Errorf("Register: %v", err)
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
		return fmt.Errorf("Invalid dependency type: %s", depType)
	}

	numOut := depType.NumOut()
	if numOut == 0 || numOut > 3 {
		return fmt.Errorf("Invalid dependency NumOut excepted: [1, 3], got: %d", depType.NumOut())
	}

	if numOut >= 2 {
		out1Name := depType.Out(1).Name()

		var hasError, hasOpts bool

		if out1Name == "error" {
			hasError = true
		} else if out1Name == "godi.DepOptions" {
			hasOpts = true
		}

		if numOut == 3 {
			out2Name := depType.Out(2).Name()

			if out2Name == "error" && hasError {
				return fmt.Errorf("Has two error out")
			}

			if out2Name == "godi.DepOptions" && hasOpts {
				return fmt.Errorf("Has two options out")
			}
		}
	}

	return nil
}
