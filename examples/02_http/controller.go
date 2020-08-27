package main

import (
	"errors"
	"fmt"
	"net/http"
)

// IController simple http controller interface
type IController interface {
	Pattern() string
	Handler(w http.ResponseWriter, r *http.Request) error
}

// HelloController hello world IController
type HelloController struct{ port IPort }

// NewHelloController create new IController
func NewHelloController(port IPort) IController {
	return &HelloController{port: port}
}

// Pattern return HelloController route patter
func (c *HelloController) Pattern() string {
	return "/hello"
}

// Handler HelloController route handler
func (c *HelloController) Handler(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte(fmt.Sprintf("Hello, world! You are on localhost:%d", c.port)))
	return nil
}

// ErrorController error IController
type ErrorController struct{}

// NewErrorController create new IController
func NewErrorController() IController {
	return &ErrorController{}
}

// Pattern return HelloController route patter
func (c *ErrorController) Pattern() string {
	return "/error"
}

// Handler HelloController route handler
func (c *ErrorController) Handler(w http.ResponseWriter, r *http.Request) error {
	return errors.New("some error here")
}
