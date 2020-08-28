package godi

import (
	"reflect"
	"testing"
)

func TestOptionsIsNil(t *testing.T) {
	if !optionsIsNil(BeanOptions{}) {
		t.Error("optionsIsNil got: false, want: true")
	}

	if optionsIsNil(BeanOptions{Type: Singleton}) {
		t.Error("optionsIsNil got: true, want: false")
	}
}

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

			t.Errorf("error in convertToTypedSlice: %v", err)
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

			if typedValuesType.Kind() != valueType.Elem().Kind() && typedValuesType.Name() != valueType.Elem().Name() {
				t.Errorf("invalid type of typeValues element: %v", typedValuesType)
				return
			}
		}
	}

	f(reflect.TypeOf([]string{}), []interface{}{"a", 'b', "c"}, false)
	f(reflect.TypeOf([]int{}), []interface{}{1, 2, 3}, false)
	f(reflect.TypeOf([][]int{}), []interface{}{[]int{1,2,3}, []int{1,3,3,7}}, false)
	f(reflect.TypeOf([][]int{}), []interface{}{[]int{1,2,3}, []int{1,3,3,7}, 4, ""}, true)
}