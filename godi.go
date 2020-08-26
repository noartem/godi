package godi

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
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

// NewContainer create new DI container with custom logger and register dependecies
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

// Resolve get factory dependencies and return factory instance
func (container *Container) Resolve(depName string, dep interface{}) (interface{}, error) {
	container.log.Printf("Resolve: %v", dep)

	if container.singletons[depName] != nil {
		return container.singletons[depName], nil
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic error: %v", r)
		}
	}()

	depVal := reflect.ValueOf(dep)
	depType := depVal.Type()

	numIn := depType.NumIn()

	var in []reflect.Value
	for i := 0; i < numIn; i++ {
		inType := depType.In(i)

		switch inType.Kind() {
		case reflect.Slice, reflect.Array:
			inValues, err := container.GetAll(inType.Elem().Name())
			if err != nil {
				return nil, fmt.Errorf("InitDep: Cannot get %s: %v", inType.Name(), err)
			}

			for _, inValue := range inValues {
				in = append(in, reflect.ValueOf(inValue))
			}
		case reflect.Struct, reflect.Interface:
			inValue, err := container.Get(inType.Name())
			if err != nil {
				return nil, fmt.Errorf("InitDep: Cannot get %s: %v", inType.Name(), err)
			}
			in = append(in, reflect.ValueOf(inValue))
		default:
			return nil, fmt.Errorf("InitDep: Invalid dependency of %s: %v", depName, inType)
		}
	}

	out := depVal.Call(in)
	resolved, options, depErr, err := parseDepInitOut(out)
	if err != nil {
		return nil, err
	}

	if depErr != nil {
		return nil, fmt.Errorf("Dep error: %v", depErr)
	}

	if options != nil && options.Type == Singleton {
		container.singletons[depName] = resolved
	}

	container.log.Printf("Dep Options: %v", options)

	return resolved, nil
}

func parseDepInitOut(rawDepOut []reflect.Value) (interface{}, *DepOptions, error, error) {
	var options *DepOptions
	var depErr error

	if len(rawDepOut) >= 2 {
		switch rawDepOut[1].Type().Name() {
		case "error":
			depErr = rawDepOut[1].Interface().(error)
		case "godi.DepOptions":
			options = rawDepOut[1].Interface().(*DepOptions)
		default:
			return nil, nil, nil, fmt.Errorf("Invalid first dep out: %v", rawDepOut[1])
		}
	}

	if len(rawDepOut) == 3 {
		switch rawDepOut[2].Type().Name() {
		case "error":
			if depErr != nil {
				return nil, nil, nil, fmt.Errorf("Invalid dep out values: Already has error")
			}

			depErr = rawDepOut[2].Interface().(error)
		case "godi.DepOptions":
			if options != nil {
				return nil, nil, nil, fmt.Errorf("Invalid dep out values: Already has options")
			}

			options = rawDepOut[2].Interface().(*DepOptions)
		default:
			return nil, nil, nil, fmt.Errorf("Invalid dep out values: %v", rawDepOut)
		}
	}

	return rawDepOut[0].Interface(), options, depErr, nil
}

// Get return last registered factory by name (interface name)
func (container *Container) Get(name string) (interface{}, error) {
	container.log.Printf("Get: %s", name)

	deps := container.deps[name]
	if deps == nil {
		return nil, fmt.Errorf("Dependencies with name %s is not found", name)
	}

	if len(deps) == 0 {
		return nil, fmt.Errorf("Dependencies with name %s is empty", name)
	}

	// return last registered dependency of this type
	dep := deps[len(deps)-1]

	iniDep, err := container.Resolve(name, dep)
	if err != nil {
		return nil, err
	}

	return iniDep, nil
}

// GetAll return all registered factories by name (interface name)
func (container *Container) GetAll(name string) ([]interface{}, error) {
	container.log.Printf("GetAll: %s", name)

	deps := container.deps[name]
	if deps == nil {
		return nil, fmt.Errorf("Dependecies with name %s are not found", name)
	}

	iniDeps := []interface{}{}
	for _, dep := range deps {
		iniDep, err := container.Resolve(name, dep)
		if err != nil {
			return iniDeps, err
		}

		iniDeps = append(iniDeps, iniDep)
	}

	return iniDeps, nil
}
