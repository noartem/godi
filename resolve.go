package godi

import (
	"fmt"
	"reflect"
)

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
