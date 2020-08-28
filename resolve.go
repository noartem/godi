package godi

import (
	"fmt"
	"reflect"
)

// resolveFactory get factory dependencies and return bean from factory
func (container *Container) resolveFactory(factoryName string, factory interface{}) (res interface{}, err error) {
	container.log.Printf("Resolve: %s = %v", factoryName, factory)

	if container.beanSingletons[factoryName] != nil {
		return container.beanSingletons[factoryName], nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("resolveFactory: %v", r)
		}
	}()

	factoryVal := reflect.ValueOf(factory)
	factoryType := factoryVal.Type()

	if factoryType.Kind() != reflect.Func {
		return factory, nil
	}

	numIn := factoryType.NumIn()
	var in []reflect.Value
	for i := 0; i < numIn; i++ {
		inType := factoryType.In(i)

		if inType.Kind() == reflect.Slice {
			inValues, err := container.GetAll(inType.Elem().Name())
			if err != nil {
				return nil, fmt.Errorf("cannot get %s: %v", inType.Name(), err)
			}

			typedValues, err := convertToTypedSlice(inType, inValues)
			if err != nil {
				return nil, err
			}

			in = append(in, typedValues)
		} else {
			inValue, err := container.Get(inType.Name())
			if err != nil {
				return nil, fmt.Errorf("cannot get %s: %v", inType.Name(), err)
			}

			in = append(in, reflect.ValueOf(inValue))
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

	if !optionsIsNil(options) && options.Type == Singleton {
		container.beanSingletons[factoryName] = resolved
	}

	container.log.Printf("Resolved: %s = %v (options: %v)", factoryName, resolved, options)

	return resolved, nil
}

func convertToTypedSlice(valuesType reflect.Type, values []interface{}) (val reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("convertToTypedSlice: %v", r)
		}
	}()

	valueType := valuesType.Elem()
	typedValues := reflect.MakeSlice(valuesType, 0, 0)

	for _, value := range values {
		typedValue := reflect.ValueOf(value).Convert(valueType)
		typedValues = reflect.Append(typedValues, typedValue)
	}

	return typedValues, nil
}

func parseFactoryOut(factoryOut []reflect.Value) (bean interface{}, options BeanOptions, factoryErr error, err error) {
	if len(factoryOut) >= 2 {
		switch factoryOut[1].Type().Name() {
		case ErrorType:
			factoryErr = factoryOut[1].Interface().(error)
		case BeanOptionsType:
			options = factoryOut[1].Interface().(BeanOptions)
		default:
			err = fmt.Errorf("invalid first factory out: %v", factoryOut[1])
			return
		}
	}

	if len(factoryOut) == 3 {
		switch factoryOut[2].Type().Name() {
		case ErrorType:
			if factoryErr != nil {
				err = fmt.Errorf("invalid factory out values: Already has error")
				return
			}

			factoryErr = factoryOut[2].Interface().(error)
		case BeanOptionsType:
			if optionsIsNil(options) {
				err = fmt.Errorf("invalid factory out values: Already has options")
				return
			}

			options = factoryOut[2].Interface().(BeanOptions)
		default:
			err = fmt.Errorf("invalid factory out values: %v", factoryOut)
			return
		}
	}

	bean = factoryOut[0].Interface()

	return
}

func optionsIsNil(options BeanOptions) bool {
	return options.Type == 0
}
