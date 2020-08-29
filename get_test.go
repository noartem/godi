package godi

import (
	"log"
	"reflect"
	"testing"
)

type iName interface{ Name() string }
type eName struct{}

func newName() iName {
	return &eName{}
}
func (name *eName) Name() string {
	return "wow"
}

type iHello interface{ Hello() }
type eHello struct{ name iName }

func newHello(name iName) iHello {
	return &eHello{name: name}
}
func (hello *eHello) Hello() {
	log.Println(hello.name.Name())
}

type iGreet interface{ Greet() }
type eGreet struct{}

func newGreet() iGreet {
	return &eGreet{}
}
func (greet *eGreet) Greet() {
	log.Println("greet")
}

func TestGet(t *testing.T) {
	c, err := NewContainer(newHello, newName)
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
		return
	}

	f := func(name string, hasError bool, original interface{}) {
		bean, err := c.Get(name)
		if err != nil {
			if hasError {
				return
			}

			t.Errorf("unexcepted error from Get: %v", err)
		}

		if !reflect.DeepEqual(bean, original) {
			t.Errorf("invalid bean: %v", err)
		}
	}

	f("godi.iName", false, &eName{})
	f("godi.iHello", false, &eHello{name: &eName{}})
	f("godi.iGreet", true, nil)

	c.factories["godi.lorem_ipsum"] = []interface{}{}
	f("godi.lorem_ipsum", true, nil)

	err = c.Register(newGreet)
	if err != nil {
		t.Errorf("unxecpted error: %v", err)
		return
	}
	f("godi.iGreet", false, &eGreet{})
}

func TestGetAll(t *testing.T) {
	c, err := NewContainer(newHello, newName)
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
		return
	}

	f := func(name string, hasError bool, original ...interface{}) {
		beans, err := c.GetAll(name)
		if err != nil {
			if hasError {
				return
			}

			t.Errorf("unexcepted error from GetAll: %v", err)
		}

		if len(original) != len(beans) || !reflect.DeepEqual(beans, original) {
			t.Errorf("invalid beans excepted: %v, got: %v", original, beans)
		}
	}

	f("godi.iName", false, &eName{})
	f("godi.iHello", false, &eHello{name: &eName{}})
	f("godi.iGreet", true)

	c.factories["godi.lorem_ipsum"] = []interface{}{}
	f("godi.lorem_ipsum", false)

	err = c.Register(newGreet)
	if err != nil {
		t.Errorf("unxecpted error: %v", err)
		return
	}
	f("godi.iGreet", false, &eGreet{})
}
