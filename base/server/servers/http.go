package servers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/labstack/echo/v4"
)

type Http struct{
	server *echo.Echo
	Modules []any
}

func (h *Http) Start() {
	h.server = echo.New()
}

func (h *Http) Handler() {
		for _, module := range h.Modules {
		methods := reflect.TypeOf(module)

		for i := 0; i < methods.NumMethod(); i++ {
			method := methods.Method(i)
			handlers := reflect.ValueOf(module).MethodByName(method.Name).Call([]reflect.Value{
				reflect.ValueOf(h.server),
			})
			
			for _, handler := range handlers[1].Interface().([]any) {
				reflect.ValueOf(handler).MethodByName("Boot").Call([]reflect.Value{
					reflect.ValueOf(handlers[0].Interface()),
				})
			}
		}
	}
}

func (h *Http) Run(address string) {
	data, err := json.MarshalIndent(h.server.Routes(), "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("routes.json", data, 0644)

	fmt.Println("http server running on " + address)
	if err := h.server.Start(address); err != nil {
		app.Log.Error(err.Error())
	}
}
