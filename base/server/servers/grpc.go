package servers

import (
	"fmt"
	"net"
	"reflect"

	"github.com/backend-timedoor/gtimekeeper/app"
	"google.golang.org/grpc"
)

type Grpc struct{
	Server *grpc.Server
	Modules []any
}

func (g *Grpc) Start() {
	g.Server = grpc.NewServer()
}

func (g *Grpc) Handler() {
		for _, module := range g.Modules {
		methods := reflect.TypeOf(module)

		for i := 0; i < methods.NumMethod(); i++ {
			method := methods.Method(i)
			handlers := reflect.ValueOf(module).MethodByName(method.Name).Call(nil)
			
			for _, handler := range handlers[0].Interface().([]any) {
				reflect.ValueOf(handler).MethodByName("Boot").Call([]reflect.Value{
					reflect.ValueOf(g.Server),
				})
			}
		}
	}
}

func (g *Grpc) Run(address string) {
		listener, err := net.Listen("tcp", address)
	if err != nil {
		app.Log.Error(err.Error())
	}

	fmt.Println("grpc server running on " + listener.Addr().String())
	if err = g.Server.Serve(listener); err != nil {
		app.Log.Error(err.Error())
	}
}
