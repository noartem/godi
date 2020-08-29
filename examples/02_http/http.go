package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/noartem/godi"
)

// IPort server port
type IPort int

// PortDefault default IPort
const PortDefault IPort = 8080

// IHttp http server
type IHttp interface {
	StartServer() error
}

// HTTP IHttp implementation
type HTTP struct {
	controllers []IController
	port        IPort
}

// NewHTTP create new IHttp
func NewHTTP(controllers []IController, port IPort) (IHttp, *godi.BeanOptions) {
	httpServer := &HTTP{
		controllers: controllers,
		port:        port,
	}

	options := &godi.BeanOptions{Type: godi.Singleton}

	return httpServer, options
}

// StartServer register controllers and start http server
func (h *HTTP) StartServer() error {
	for i := 0; i < len(h.controllers); i++ {
		controller := h.controllers[i]
		http.HandleFunc(controller.Pattern(), func(w http.ResponseWriter, r *http.Request) {
			err := controller.Handler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(fmt.Sprintf("Error: %v!", err)))
				if err != nil {
					log.Printf("error in Write: %v", err)
				}
			}
		})
	}

	portString := fmt.Sprintf(":%d", h.port)

	return http.ListenAndServe(portString, nil)
}
