package log

import (
	"fmt"
	"runtime/debug"
	"slices"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"github.com/sirupsen/logrus"
)

type DebugStackHook struct {
	GetValue func() string
}

func (h *DebugStackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *DebugStackHook) Fire(e *logrus.Entry) error {
	// e.Data["aDefaultField"] = h.GetValue()
	environments := []string{"local", "development"}
	levels := []string{"error", "trace"}

	if slices.Contains(environments, app.Config.GetString("app.env")) {
		if slices.Contains(levels, e.Level.String()) {
			fmt.Println(string(debug.Stack()))
		}
	}

	return nil
}