package godi

import (
	"fmt"
	"reflect"
	"testing"
)

var valOf = reflect.ValueOf
var typOf = reflect.TypeOf

func TestConvertToTypedSlice(t *testing.T) {
	f := func(valueType reflect.Type, values []interface{}, hasError bool) {
		typedValues, err := convertToTypedSlice(valueType, values)
		if err != nil && hasError {
			return
		}

		if err != nil {
			if hasError {
				return
			}

			t.Errorf("unexcepted error: %v", err)
			return
		}

		if typedValues.Kind() != reflect.Slice {
			t.Errorf("invalid typedValues kind want: slice, got: %s", typedValues.Kind())
			return
		}

		if len(values) != typedValues.Len() {
			t.Errorf("invalid typedValues length want: %d, got: %d", len(values), typedValues.Len())
			return
		}

		for i := 0; i < typedValues.Len(); i++ {
			typedValuesValue := typedValues.Index(i)
			typedValuesType := typedValuesValue.Type()

			if typedValuesType.Kind() != valueType.Elem().Kind() && typedValuesType.String() != valueType.Elem().String() {
				t.Errorf("invalid type of typeValues element: %v", typedValuesType)
				return
			}
		}
	}

	f(typOf([]string{}), []interface{}{"a", 'b', "c"}, false)
	f(typOf([]int{}), []interface{}{1, 2, 3}, false)
	f(typOf([][]int{}), []interface{}{[]int{1, 2, 3}, []int{1, 3, 3, 7}}, false)
	f(typOf([][]int{}), []interface{}{[]int{1, 2, 3}, []int{1, 3, 3, 7}, 4, ""}, true)
}

func TestParseFactoryOutValues(t *testing.T) {
	f := func(hasError bool, values ...interface{}) {
		var rValues []reflect.Value
		for _, value := range values {
			rValues = append(rValues, valOf(value))
		}

		factoryOut, err := parseFactoryOutValues(rValues)
		if err != nil {
			if hasError {
				return
			}

			t.Errorf("unexcepted error: %v", err)
			return
		}

		if !reflect.DeepEqual(factoryOut.bean, values[0]) {
			t.Errorf("invalid factoryOut.bean excepted: %v, got: %v", values[0], factoryOut.bean)
			return
		}
	}

	bean := 1337
	err1 := fmt.Errorf("some error")
	err2 := fmt.Errorf("another error")
	opts1 := &BeanOptions{Type: Singleton}
	opts2 := &BeanOptions{Type: Prototype}

	f(false, bean)
	f(false, "wow")
	f(false, []int{1, 3, 3, 7})

	f(false, bean, err1)
	f(true, bean, err1, err2)

	f(false, bean, opts1)
	f(false, bean, err1, opts1)
	f(false, bean, opts1, err1)
	f(true, bean, err1, opts1, err2)
	f(true, bean, err1, opts1, opts2)
	f(true, bean, opts1, opts2)
	f(true, bean, opts1, opts2, err1)
}
