package servers

import (
	"fmt"
	"net"
	"reflect"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"google.golang.org/grpc"
)

type Grpc struct {
	Server  *grpc.Server
	Modules []any
}

func (g *Grpc) Start() {
	g.Server = grpc.NewServer()
}

//func (g *Grpc) Handler() {
//	for _, module := range g.Modules {
//		g.registerHandler(module)
//	}
//}

func (g *Grpc) RegisterHandler(t any) {
	methods := reflect.TypeOf(t)
	_, isHandler := methods.MethodByName("Boot")

	if !isHandler {
		for i := 0; i < methods.NumMethod(); i++ {
			method := methods.Method(i)
			fmt.Println(method.Name)
			execMethod := reflect.ValueOf(t).MethodByName(method.Name).Call(nil)

			for _, instance := range execMethod[1].Interface().([]any) {
				g.RegisterHandler(instance)
			}
		}
	} else {
		reflect.ValueOf(t).MethodByName("Boot").Call([]reflect.Value{
			reflect.ValueOf(g.Server),
		})
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
