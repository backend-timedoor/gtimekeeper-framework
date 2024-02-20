package container

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

var (
	App map[string]any
)

type Item struct {
	Key   string
	Value any
}

func Set(key string, value any) {
	App[key] = value
}

func Get(key string) any {
	return App[key]
}

func ExecRef(key string, method string, args []reflect.Value) []reflect.Value {
	return reflect.ValueOf(App[key]).MethodByName(method).Call(args)
}

func Log() *logrus.Logger {
	return App["log"].(*logrus.Logger)
}
