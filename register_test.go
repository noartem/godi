package godi

import (
	"fmt"
	"testing"
)

func TestCheckFactoryOut(t *testing.T) {
	f := func(factory interface{}, hasError bool) {
		err := checkFactoryOut(typOf(factory))
		if err != nil {
			if hasError {
				return
			}

			t.Errorf("unexcepted checkFactoryOut error: %v", err)
			return
		}
	}

	bean := 1337
	err1 := fmt.Errorf("1")
	err2 := fmt.Errorf("2")
	options1 := &BeanOptions{Type: Singleton}
	options2 := &BeanOptions{Type: Prototype}

	f(func() int { return bean }, false)
	f(func() (int, error) { return bean, err1 }, false)
	f(func() (int, *BeanOptions, error) { return bean, options1, err1 }, false)

	f("", true)
	f(func() (int, error, error) { return bean, err1, err2 }, true)
	f(func() (int, *BeanOptions, *BeanOptions) { return bean, options1, options2 }, true)
	f(func() (int, *BeanOptions, error, error) { return bean, options1, err1, err2 }, true)
	f(func() (int, int, error) { return bean, bean, err1 }, true)
	f(func() (int, *BeanOptions, int) { return bean, options1, bean }, true)
	f(func() (int, int, int) { return bean, bean, bean }, true)
	f(func() (int, error, int) { return bean, err1, bean }, true)
}
