package godi

import "fmt"

// Get return last registered factory by name (interface name)
func (container *Container) Get(name string) (interface{}, error) {
	container.log.Printf("Get: %s", name)

	deps := container.deps[name]
	if deps == nil {
		return nil, fmt.Errorf("dependencies with name %s is not found", name)
	}

	if len(deps) == 0 {
		return nil, fmt.Errorf("dependencies with name %s is empty", name)
	}

	// return last registered dependency of this type
	dep := deps[len(deps)-1]

	iniDep, err := container.Resolve(name, dep)
	if err != nil {
		return nil, err
	}

	return iniDep, nil
}

// GetAll return all registered factories by name (interface name)
func (container *Container) GetAll(name string) ([]interface{}, error) {
	container.log.Printf("GetAll: %s", name)

	deps := container.deps[name]
	if deps == nil {
		return nil, fmt.Errorf("dependecies with name %s are not found", name)
	}

	iniDeps := []interface{}{}
	for _, dep := range deps {
		iniDep, err := container.Resolve(name, dep)
		if err != nil {
			return iniDeps, err
		}

		iniDeps = append(iniDeps, iniDep)
	}

	return iniDeps, nil
}
