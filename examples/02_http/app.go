package main

import (
	"fmt"

	"github.com/noartem/godi"
)

// IApp main application bean
type IApp interface {
	// start application
	Start() error
}

// App default IApp
type App struct {
	http IHttp
}

// NewApp create new App
func NewApp(http IHttp) (IApp, godi.BeanOptions) {
	app := &App{
		http: http,
	}

	options := godi.BeanOptions{
		Type: godi.Singleton,
	}

	return app, options
}

// Start start application
func (app *App) Start() error {
	fmt.Println("Starting server...")
	return app.http.StartServer()
}
