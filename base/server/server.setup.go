package server

import (
    "github.com/backend-timedoor/gtimekeeper-framework/app"
    "github.com/backend-timedoor/gtimekeeper-framework/base/contracts"
    "github.com/backend-timedoor/gtimekeeper-framework/base/server/servers"
    "github.com/backend-timedoor/gtimekeeper-framework/container"
    "reflect"
)

const ContainerName string = "server"

func New(modules []any) contracts.Server {
    serverBag := map[string]any{
        "Grpc": &servers.Grpc{},
        "Http": &servers.Http{},
    }

    for _, module := range modules {
        setServerPreparation(module, serverBag)
    }

    s := &Server{}
    s.Servers = serverBag
    s.Validation = app.Validation

    container.Set(ContainerName, s)

    return s
}

func setServerPreparation(module any, serverBag map[string]any) {
    refModule := reflect.TypeOf(module)
    //fmt.Println(refModule)

    for i := 0; i < refModule.NumMethod(); i++ {
        method := refModule.Method(i)
        methodType := method.Type
        countReturn := methodType.NumOut()

        if countReturn == 2 { // 2 is grpc setup
            grpcServer := serverBag["Grpc"].(*servers.Grpc)
            if grpcServer.Server == nil {
                grpcServer.Start()
            }

            execMethod := reflect.ValueOf(module).MethodByName(method.Name).Call(nil)
            for _, mod := range execMethod[1].Interface().([]any) {
                grpcServer.RegisterHandler(mod)
            }
        } else if countReturn == 3 {
            httpServer := serverBag["Http"].(*servers.Http)
            if httpServer.Server == nil {
                httpServer.Start()
            }

            route := reflect.ValueOf(httpServer.Server.Group(""))
            execMethod := reflect.ValueOf(module).MethodByName(method.Name).Call([]reflect.Value{route})
            for _, mod := range execMethod[2].Interface().([]any) {
                httpServer.RegisterHandler(mod, execMethod[1].Interface())
            }
        } else {
            execMethod := reflect.ValueOf(module).MethodByName(method.Name).Call(nil)

            for _, instance := range execMethod[0].Interface().([]any) {
                setServerPreparation(instance, serverBag)
            }
        }
    }
}
