package fxlogr

import (
	"github.com/go-logr/logr"
	"go.uber.org/fx/fxevent"
)

type logger struct {
	Logger logr.Logger
}

func (l *logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Logger.Info("OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName)
	}
}

func WithLogr(l logr.Logger) func() fxevent.Logger {
	return func() fxevent.Logger {
		return &logger{Logger: l}
	}
}
