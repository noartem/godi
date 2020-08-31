# Godi [![PkgGoDev](https://pkg.go.dev/badge/github.com/noartem/godi)](https://pkg.go.dev/github.com/noartem/godi)

Simple Golang [DI](https://en.wikipedia.org/wiki/Dependency_injection) container based on reflection

Examples: [examples folder](https://github.com/noartem/godi/tree/master/examples)

Godoc: [pkg.go.dev](https://pkg.go.dev/github.com/noartem/godi)

## Get started

1. Install godi `go get github.com/noartem/godi`
2. Create interfaces and beans implementing these interfaces:

   ```go
   type IName interface {
       NewName() string
   }

   type Name struct {}

   func (name *Name) Generate() string {
       return "Lorem Ipsumovich"
   }
   ```

3. Create bean factory:

   ```go
   func NewName() IName {
       return &Name{ ... }
   }
   ```

   In factory, you can import other beans:

   ```go
   func NewName(db IDatabase, log ILogger, ...) IName {
       return &Name{
           db: db,
           log: log,
       }
   }
   ```

   Factories can also return BeanOptions and/or errors:

   ```go
   func NewName() (IName, *godi.BeanOptions, error) {
       err := someFunc()
       if err != nil {
           return &Name{}, nil, err
       }

       options := &godi.BeanOptions{
           Type: godi.Singleton, // Default: godi.Prototype
       }

       name := &Name{}

       return name, options, nil
   }
   ```

   or `func NewName() (IName, error) {}`, or `func NewName() (IName, *godi.BeanOptions)`

4. Create DI container and register factories:

   ```go
   func main() {
       c, err := godi.NewContainer(NewName, NewRandom, NewFoo, ...)
       if err != nil {
           panic(err)
       }

       ...
   ```

5. Get bean from a container:

   ```go
       // get bean by interface name
       nameBean, err := c.Get("IName")
       if err != nil {
           panic(err)
       }

       name, ok := nameBean.(IName)
       if !ok {
           panic("Invalid name bean")
       }

       // now you can use IName
       fmt.Println(name.Generate())
   }
   ```

6. Build your architecture based on [IOC](https://en.wikipedia.org/wiki/Inversion_of_control) with [DI](https://en.wikipedia.org/wiki/Dependency_injection)

7. Profit! See another examples in [examples folder](https://github.com/noartem/godi/tree/master/examples)

## Other features

1. Static beans. 
    Can be used for sharing constants (global config, secrets, e.t.c.)

   1. Create a custom type

      ```go
      type IPassword string
      ```

   2. Create implementations

      ```go
      var defaultPassword IPassword = "qwerty123"
      ```

   3. Register

      ```go
      godi.NewContainer(defaultPassword)
      ```

      Constant will be registered as `IPassword`

   4. Use it

      ```go
      func NewName(password IPassword) IName {...}
      ```

2. Structures with dependencies in Input. 
   If you don't want to write boilerplate code creating structure with all input dependency, you can write them 
   in structure with `InStruct` field and require this as input in factory.

   1. Create structure based on godi.InStruct and add your dependencies:

       ```go
       import "github.com/noartem/godi"
    
       type deps struct {
          godi.InStruct
       
          Name IName
          Config IConfig
          Random IRandom
       }
       ```

   2. Require this structure in factory:
   
       ```go
       func NewHello(deps deps) IHello { ... }
       ```

   3. Use your dependencies:

       ```go
       func NewHello(deps deps) IHello {
          log.Println(deps.Name.GenerateName())
          log.Println(deps.Config)
          log.Println(deps.Random.Intn(1337))
       }
       ```
