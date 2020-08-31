package godi

import (
	"fmt"
	"reflect"
)

// resolveFactory get factory dependencies and return bean from factory
func (container *Container) resolveFactory(factory interface{}) (res interface{}, err error) {
	factoryVal := reflect.ValueOf(factory)
	factoryType := factoryVal.Type()
	factoryName := genFactoryName(factoryType)

	container.log.Printf("Resolve: %s = %v", factoryName, factory)

	if container.beanSingletons[factoryName] != nil {
		return container.beanSingletons[factoryName], nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("resolveFactory: %v", r)
		}
	}()

	if factoryType.Kind() != reflect.Func {
		return factory, nil
	}

	numIn := factoryType.NumIn()
	var in []reflect.Value
	for i := 0; i < numIn; i++ {
		inValue, err := container.switchInKind(factoryType.In(i))
		if err != nil {
			return nil, err
		}

		in = append(in, inValue)
	}

	out := factoryVal.Call(in)
	factoryOut, err := parseFactoryOutValues(out)
	if err != nil {
		return nil, err
	}

	if factoryOut.err != nil {
		return nil, fmt.Errorf("error from factory: %v", factoryOut.err)
	}

	if factoryOut.options != nil && factoryOut.options.Type == Singleton {
		container.beanSingletons[factoryName] = factoryOut.bean
	}

	container.log.Printf("Resolved: %s = %v (options: %v)", factoryType, factoryOut.bean, factoryOut.options)

	return factoryOut.bean, nil
}

func (container *Container) switchInKind(inType reflect.Type) (reflect.Value, error) {
	switch inType.Kind() {
	case reflect.Slice:
		inValues, err := container.GetAll(inType.Elem().String())
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot get %s: %v", inType.String(), err)
		}

		typedValues, err := convertToTypedSlice(inType, inValues)
		if err != nil {
			return reflect.Value{}, err
		}

		return typedValues, nil
	case reflect.Struct:
		inStruct, ok := inType.FieldByName(InStructType.Name())
		if !(ok && inStruct.Type.String() == InStructType.String()) {
			return container.switchInKindDefault(inType)
		}

		filledStruct, err := container.fillStruct(inType)
		if err != nil {
			return reflect.Value{}, err
		}

		return filledStruct, nil
	default:
		return container.switchInKindDefault(inType)
	}
}

func (container *Container) fillStruct(structType reflect.Type) (reflect.Value, error) {
	newStruct := reflect.New(structType).Elem()

	for i := 0; i < newStruct.NumField(); i++ {
		structField := newStruct.Field(i)
		if structField.Type().Name() == InStructType.Name() {
			continue
		}

		fieldValue, err := container.Get(structField.Type().String())
		if err != nil {
			return reflect.Value{}, err
		}

		structField.Set(reflect.ValueOf(fieldValue))
	}

	return newStruct, nil
}

func (container *Container) switchInKindDefault(inType reflect.Type) (reflect.Value, error) {
	inValue, err := container.Get(inType.String())
	if err != nil {
		return reflect.Value{}, fmt.Errorf("cannot get %s: %v", inType.String(), err)
	}

	return reflect.ValueOf(inValue), nil
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

type factoryOut struct {
	bean    interface{}
	options *BeanOptions
	err     error
}

func parseFactoryOutValues(values []reflect.Value) (*factoryOut, error) {
	out := &factoryOut{}
	if len(values) >= 2 {
		out1 := values[1].Type()
		if out1.Implements(ErrorInterface) {
			out.err = values[1].Interface().(error)
		} else if out1 == BeanOptionsType {
			out.options = values[1].Interface().(*BeanOptions)
		} else {
			return nil, fmt.Errorf("invalid first factory out type: %v", values[1].Type().String())
		}
	}

	if len(values) == 3 {
		out2 := values[2].Type()
		if out2.Implements(ErrorInterface) {
			if out.err != nil {
				return nil, fmt.Errorf("invalid factory out values: Already has error")
			}

			out.err = values[2].Interface().(error)
		} else if out2 == BeanOptionsType {
			if out.options != nil {
				return nil, fmt.Errorf("invalid factory out values: Already has options")
			}

			out.options = values[2].Interface().(*BeanOptions)
		} else {
			return nil, fmt.Errorf("invalid factory out values: %v", values)
		}
	}

	out.bean = values[0].Interface()

	return out, nil
}
