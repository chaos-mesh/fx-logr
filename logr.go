// Copyright 2023 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fxlogr

import (
	"strings"

	"github.com/go-logr/logr"
	"go.uber.org/fx/fxevent"
)

type LogrLogger struct {
	Logger *logr.Logger

	logLevel   int
	errorLevel int
}

var _ fxevent.Logger = (*LogrLogger)(nil)

// UseLogLevel sets the log level for log events.
func (l *LogrLogger) UseLogLevel(level int) {
	l.logLevel = level
}

// UseErrorLevel sets the log level for error events.
func (l *LogrLogger) UseErrorLevel(level int) {
	l.errorLevel = level
}

func (l *LogrLogger) logEvent(msg string, keysAndValues ...interface{}) {
	l.Logger.V(l.logLevel).Info(msg, keysAndValues...)
}

func (l *LogrLogger) logError(err error, msg string, keysAndValues ...interface{}) {
	l.Logger.V(l.errorLevel).Error(err, msg, keysAndValues...)
}

// LogEvent logs an event to the provided Logr logger.
func (l *LogrLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logEvent("OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logError(e.Err, "OnStart hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
			)
		} else {
			l.logEvent("OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		l.logEvent("OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logError(e.Err, "OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
			)
		} else {
			l.logEvent("OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		if len(e.ModuleName) != 0 {
			if e.Err != nil {
				l.logError(e.Err, "error encountered while applying options",
					"type", e.TypeName,
					"module", e.ModuleName,
				)
			} else {
				l.logEvent("supplied",
					"type", e.TypeName,
					"module", e.ModuleName,
				)
			}
		} else {
			if e.Err != nil {
				l.logError(e.Err, "error encountered while applying options",
					"type", e.TypeName,
				)
			} else {
				l.logEvent("supplied",
					"type", e.TypeName,
				)
			}
		}
	case *fxevent.Provided:
		if len(e.ModuleName) != 0 {
			for _, rtype := range e.OutputTypeNames {
				if e.Private {
					l.logEvent("provided",
						"constructor", e.ConstructorName,
						"module", e.ModuleName,
						"type", rtype,
						"private", true,
					)
				} else {
					l.logEvent("provided",
						"constructor", e.ConstructorName,
						"module", e.ModuleName,
						"type", rtype,
					)
				}
			}
			if e.Err != nil {
				l.logError(e.Err, "error encountered while applying options",
					"module", e.ModuleName,
				)
			}
		} else {
			for _, rtype := range e.OutputTypeNames {
				if e.Private {
					l.logEvent("provided",
						"constructor", e.ConstructorName,
						"type", rtype,
						"private", true,
					)
				} else {
					l.logEvent("provided",
						"constructor", e.ConstructorName,
						"type", rtype,
					)
				}
			}
			if e.Err != nil {
				l.logError(e.Err, "error encountered while applying options")
			}
		}
	case *fxevent.Replaced:
		if len(e.ModuleName) != 0 {
			for _, rtype := range e.OutputTypeNames {
				l.logEvent("replaced",
					"module", e.ModuleName,
					"type", rtype,
				)
			}
			if e.Err != nil {
				l.logError(e.Err, "error encountered while replacing",
					"module", e.ModuleName,
				)
			}
		} else {
			for _, rtype := range e.OutputTypeNames {
				l.logEvent("replaced",
					"type", rtype,
				)
			}
			if e.Err != nil {
				l.logError(e.Err, "error encountered while replacing")
			}
		}
	case *fxevent.Decorated:
		if len(e.ModuleName) != 0 {
			for _, rtype := range e.OutputTypeNames {
				l.logEvent("decorated",
					"decorator", e.DecoratorName,
					"module", e.ModuleName,
					"type", rtype,
				)
			}
			if e.Err != nil {
				l.logError(e.Err, "error encountered while applying options",
					"module", e.ModuleName,
				)
			}
		} else {
			for _, rtype := range e.OutputTypeNames {
				l.logEvent("decorated",
					"decorator", e.DecoratorName,
					"type", rtype,
				)
			}
			if e.Err != nil {
				l.logError(e.Err, "error encountered while applying options")
			}
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		if len(e.ModuleName) != 0 {
			l.logEvent("invoking",
				"function", e.FunctionName,
				"module", e.ModuleName,
			)
		} else {
			l.logEvent("invoking",
				"function", e.FunctionName,
			)
		}
	case *fxevent.Invoked:
		if len(e.ModuleName) != 0 {
			if e.Err != nil {
				l.logError(e.Err, "invoke failed",
					"stack", e.Trace,
					"function", e.FunctionName,
					"module", e.ModuleName,
				)
			}
		} else {
			if e.Err != nil {
				l.logError(e.Err, "invoke failed",
					"stack", e.Trace,
					"function", e.FunctionName,
				)
			}
		}
	case *fxevent.Stopping:
		l.logEvent("received signal",
			"signal", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.logError(e.Err, "stop failed")
		}
	case *fxevent.RollingBack:
		l.logError(e.StartErr, "start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logError(e.Err, "rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.logError(e.Err, "start failed")
		} else {
			l.logEvent("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logError(e.Err, "custom logger initialization failed")
		} else {
			l.logEvent("initialized custom fxevent.Logger", "function", e.ConstructorName)
		}
	}

}

// WithLogr returns a function that returns a fxevent.Logger backed by a logr.Logger.
func WithLogr(l *logr.Logger) func() fxevent.Logger {
	return func() fxevent.Logger {
		return &LogrLogger{Logger: l}
	}
}
