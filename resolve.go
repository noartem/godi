package godi

import (
	"fmt"
	"reflect"
)

// resolveFactory get factory dependencies and return bean from factory
func (container *Container) resolveFactory(factoryName string, factory interface{}) (interface{}, error) {
	container.log.Printf("Resolve: %v", factory)

	if container.beanSingletons[factoryName] != nil {
		return container.beanSingletons[factoryName], nil
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic error: %v", r)
		}
	}()

	factoryVal := reflect.ValueOf(factory)
	factoryType := factoryVal.Type()

	numIn := factoryType.NumIn()
	var in []reflect.Value
	for i := 0; i < numIn; i++ {
		inType := factoryType.In(i)

		switch inType.Kind() {
		case reflect.Slice:
			inValues, err := container.GetAll(inType.Elem().Name())
			if err != nil {
				return nil, fmt.Errorf("cannot get %s: %v", inType.Name(), err)
			}

			for _, inValue := range inValues {
				in = append(in, reflect.ValueOf(inValue))
			}
		case reflect.Struct, reflect.Interface: // TODO: Maybe remove "reflect.Struct"?
			inValue, err := container.Get(inType.Name())
			if err != nil {
				return nil, fmt.Errorf("cannot get %s: %v", inType.Name(), err)
			}

			in = append(in, reflect.ValueOf(inValue))
		default:
			return nil, fmt.Errorf("invalid dependency of %s: %v", factoryName, inType)
		}
	}

	out := factoryVal.Call(in)
	resolved, options, factoryErr, err := parseFactoryOut(out)
	if err != nil {
		return nil, err
	}

	if factoryErr != nil {
		return nil, fmt.Errorf("error from factory: %v", factoryErr)
	}

	if options != nil && options.Type == Singleton {
		container.beanSingletons[factoryName] = resolved
	}

	container.log.Printf("Dep Options: %v", options)

	return resolved, nil
}

func parseFactoryOut(factoryOut []reflect.Value) (interface{}, *BeanOptions, error, error) {
	var options *BeanOptions
	var factoryErr error

	if len(factoryOut) >= 2 {
		switch factoryOut[1].Type().Name() {
		case ErrorType:
			factoryErr = factoryOut[1].Interface().(error)
		case BeanOptionsType:
			options = factoryOut[1].Interface().(*BeanOptions)
		default:
			return nil, nil, nil, fmt.Errorf("invalid first factory out: %v", factoryOut[1])
		}
	}

	if len(factoryOut) == 3 {
		switch factoryOut[2].Type().Name() {
		case ErrorType:
			if factoryErr != nil {
				return nil, nil, nil, fmt.Errorf("invalid factory out values: Already has error")
			}

			factoryErr = factoryOut[2].Interface().(error)
		case BeanOptionsType:
			if options != nil {
				return nil, nil, nil, fmt.Errorf("invalid factory out values: Already has options")
			}

			options = factoryOut[2].Interface().(*BeanOptions)
		default:
			return nil, nil, nil, fmt.Errorf("invalid factory out values: %v", factoryOut)
		}
	}

	bean := factoryOut[0].Interface()

	return bean, options, factoryErr, nil
}
