package servers

import (
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Http struct {
	Server  *echo.Echo
	Modules []any
}

func (h *Http) Start() {
	h.Server = echo.New()
	// h.Server.Use(middleware.Recover())
}

func (h *Http) Handler() {
	for _, module := range h.Modules {
		h.registerHandler(module, h.Server.Group(""))
	}
}

func (h *Http) registerHandler(t any, args ...any) {
	var route reflect.Value
	if len(args) >= 1 {
		route = reflect.ValueOf(args[0])
	}

	methods := reflect.TypeOf(t)
	_, isHandler := methods.MethodByName("Boot")

	if !isHandler {
		for i := 0; i < methods.NumMethod(); i++ {
			method := methods.Method(i)

			execMethod := reflect.ValueOf(t).MethodByName(method.Name).Call([]reflect.Value{route})

			for _, instance := range execMethod[1].Interface().([]any) {
				h.registerHandler(instance, execMethod[0].Interface())
			}
		}
	} else {
		reflect.ValueOf(t).MethodByName("Boot").Call([]reflect.Value{route})
	}
}

func (h *Http) Run(address string) {
	data, err := json.MarshalIndent(h.Server.Routes(), "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("routes.json", data, 0644)

	if err := h.Server.Start(address); err != nil {
		log.Fatal(err.Error())
	}
}
