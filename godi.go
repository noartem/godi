package godi

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

type Container struct {
	deps map[string][]interface{}
	log  *log.Logger
}

func NewContainer(deps ...interface{}) (*Container, error) {
	logger := log.New(ioutil.Discard, "", 0)

	return NewContainerWithLogger(logger, deps...)
}

func NewContainerWithLogger(logger *log.Logger, deps ...interface{}) (*Container, error) {
	container := &Container{
		deps: make(map[string][]interface{}),
		log:  logger,
	}

	err := container.Register(deps)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (container *Container) Register(deps []interface{}) error {
	for _, dep := range deps {
		err := container.RegisterOne(dep)
		if err != nil {
			return err
		}
	}

	return nil
}

func (container *Container) RegisterOne(dep interface{}) error {
	container.log.Printf("Register: %v", dep)

	depVal := reflect.ValueOf(dep)
	depType := depVal.Type()

	if depVal.Kind() != reflect.Func {
		return fmt.Errorf("Invalid dependency type: %s", dep)
	}

	if depType.NumOut() != 1 {
		return fmt.Errorf("Invalid dependency NumOut got: %d, excepted: 1", depType.NumOut())
	}

	name := depType.Out(0).Name()

	if container.deps[name] == nil {
		container.deps[name] = []interface{}{}
	}

	container.deps[name] = append(container.deps[name], dep)

	return nil
}

func (container *Container) InitDep(dep interface{}) (interface{}, error) {
	container.log.Printf("Resolve: %v", dep)

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
		inValue, err := container.Get(inType.Name())
		if err != nil {
			return nil, fmt.Errorf("InitDep: Cannot get %s: %v", inType.Name(), err)
		}

		in = append(in, reflect.ValueOf(inValue))
	}

	depOut := depVal.Call(in)

	return depOut[0].Interface(), nil
}

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

	iniDep, err := container.InitDep(dep)
	if err != nil {
		return nil, err
	}

	return iniDep, nil
}

func (container *Container) GetAll(name string) ([]interface{}, error) {
	container.log.Printf("GetAll: %s", name)

	deps := container.deps[name]
	if deps == nil {
		return nil, fmt.Errorf("Dependecies with name %s are not found", name)
	}

	iniDeps := []interface{}{}
	for _, dep := range deps {
		iniDep, err := container.InitDep(dep)
		if err != nil {
			return iniDeps, err
		}

		iniDeps = append(iniDeps, iniDep)
	}

	return iniDeps, nil
}
