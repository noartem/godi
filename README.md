# Godi

[![PkgGoDev](https://pkg.go.dev/badge/github.com/noartem/godi)](https://pkg.go.dev/github.com/noartem/godi)

Simple Golang DI container based on reflection

Example: See [examples folder](https://github.com/noartem/godi/tree/master/examples)

Godoc: [pkg.go.dev](https://pkg.go.dev/github.com/noartem/godi)

## Get started

1. Install godi `go get github.com/noartem/godi`
2. Create interfaces and services implementing these interfaces:

   ```go
   type IName interface {
       NewName() string
   }

   type Name struct {}

   func (name *Name) Generate() {
       return "Lorem Ipsumovich"
   }
   ```

3. Create interfaces implementations factories:

   ```go
   func NewName() IName {
       return &Name{ ... }
   }
   ```

   In factories you can import other registered services:

   ```go
   func NewName(db IDatabase, log ILogger, ...) IName {
       return &Name{
           db: db,
           log: log,
       }
   }
   ```

   Factories can also return DepOptions and/or errors:

   ```go
   func NewName() (IName, *godi.DepOptions, error) {
       err := someFunc()
       if err != nil {
           return nil, nil, err
       }

       options := &godi.DepOptions{
           Type: godi.Singleton, // Default: godi.Prototype
       }

       return &Name{}, options, nil
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

5. Get service from a DI container:

   ```go
       // get service by interface name
       nameServiceRaw, err := c.Get("IName")
       if err != nil {
           panic(err)
       }

       nameService, ok := nameServiceRaw.(IName)
       if !ok {
           panic("Invalid name service")
       }

       // now you can use name service
       fmt.Println(nameService.Generate())
   }
   ```

6. Have fun! full example in [examples folder](https://github.com/noartem/godi/tree/master/examples)
