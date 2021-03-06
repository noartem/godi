package godi

import (
	"fmt"
	"reflect"
)

// Get return bean from last registered by name (interface name) factory
func (container *Container) Get(name string) (interface{}, error) {
	container.log.Printf("Get: %s", name)

	factories := container.getFactories(name)
	if factories == nil {
		return nil, fmt.Errorf("factories with name %s is not found", name)
	}

	if len(factories) == 0 {
		return nil, fmt.Errorf("factories with name %s is empty", name)
	}

	// return last registered dependency of this type
	factory := factories[len(factories)-1]

	bean, err := container.resolveFactory(factory)
	if err != nil {
		return nil, err
	}

	return bean, nil
}

// GetAll return beans from all registered by name (interface name) factories
func (container *Container) GetAll(name string) ([]interface{}, error) {
	container.log.Printf("GetAll: %s", name)

	factories := container.getFactories(name)
	if factories == nil {
		return nil, fmt.Errorf("dependecies with name %s are not found", name)
	}

	var beans []interface{}
	for _, factory := range factories {
		bean, err := container.resolveFactory(factory)
		if err != nil {
			return beans, err
		}

		beans = append(beans, bean)
	}

	return beans, nil
}

func genFactoryName(factoryType reflect.Type) string {
	if factoryType.Kind() == reflect.Func {
		factoryType = factoryType.Out(0)
	}

	return factoryType.String()
}

func (container *Container) getFactories(name string) []interface{} {
	factories := container.factories[name]

	// try to get factories without prefix
	if factories == nil {
		factories = container.factories["main."+name]
	}

	return factories
}
