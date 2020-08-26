# Godi [![PkgGoDev](https://pkg.go.dev/badge/github.com/noartem/godi)](https://pkg.go.dev/github.com/noartem/godi)

Simple Golang [DI](https://en.wikipedia.org/wiki/Dependency_injection) container based on reflection

Example: [examples folder](https://github.com/noartem/godi/tree/master/examples)

Godoc: [pkg.go.dev](https://pkg.go.dev/github.com/noartem/godi)

## Get started

1. Install godi `go get github.com/noartem/godi`
2. Create interfaces and beans implementing these interfaces:

   ```go
   type IName interface {
       NewName() string
   }

   type Name struct {}

   func (name *Name) Generate() {
       return "Lorem Ipsumovich"
   }
   ```

3. Create bean factory:

   ```go
   func NewName() IName {
       return &Name{ ... }
   }
   ```

   In factory you can import other beans:

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
           return nil, nil, err
       }

       options := &godi.BeanOptions{
           Type: godi.Singleton, // Default: godi.Prototype
       }

       name := &Name{}

       return name, options, nil
   }
   ```

   or `func NewName() (IName, error) {}`, or `func NewName() (IName, *DepOptions)`

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

7. Profit! Full example in [examples folder](https://github.com/noartem/godi/tree/master/examples)
