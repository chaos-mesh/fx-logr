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
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx/fxevent"
)

func TestLogrLogger(t *testing.T) {
	someError := errors.New("some error")

	tests := []struct {
		name        string
		give        fxevent.Event
		wantMessage string
	}{
		{
			name: "OnStartExecuting",
			give: &fxevent.OnStartExecuting{
				FunctionName: "hook.onStart",
				CallerName:   "bytes.NewBuffer",
			},
			wantMessage: "\"level\"=0 \"msg\"=\"OnStart hook executing\" \"callee\"=\"hook.onStart\" \"caller\"=\"bytes.NewBuffer\"",
		},
		{
			name: "OnStopExecuting",
			give: &fxevent.OnStopExecuting{
				FunctionName: "hook.onStop1",
				CallerName:   "bytes.NewBuffer",
			},
			wantMessage: "\"level\"=0 \"msg\"=\"OnStop hook executing\" \"callee\"=\"hook.onStop1\" \"caller\"=\"bytes.NewBuffer\"",
		},
		{

			name: "OnStopExecuted/Error",
			give: &fxevent.OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			wantMessage: "\"msg\"=\"OnStop hook failed\" \"error\"=\"some error\" \"callee\"=\"hook.onStart1\" \"caller\"=\"bytes.NewBuffer\"",
		},
		{
			name: "OnStopExecuted",
			give: &fxevent.OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Runtime:      time.Millisecond * 3,
			},
			wantMessage: "\"level\"=0 \"msg\"=\"OnStop hook executed\" \"callee\"=\"hook.onStart1\" \"caller\"=\"bytes.NewBuffer\" \"runtime\"=\"3ms\"",
		},
		{

			name: "OnStartExecuted/Error",
			give: &fxevent.OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			wantMessage: "\"msg\"=\"OnStart hook failed\" \"error\"=\"some error\" \"callee\"=\"hook.onStart1\" \"caller\"=\"bytes.NewBuffer\"",
		},
		{
			name: "OnStartExecuted",
			give: &fxevent.OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Runtime:      time.Millisecond * 3,
			},
			wantMessage: "\"level\"=0 \"msg\"=\"OnStart hook executed\" \"callee\"=\"hook.onStart1\" \"caller\"=\"bytes.NewBuffer\" \"runtime\"=\"3ms\"",
		},
		{
			name:        "Supplied",
			give:        &fxevent.Supplied{TypeName: "*bytes.Buffer"},
			wantMessage: "\"level\"=0 \"msg\"=\"supplied\" \"type\"=\"*bytes.Buffer\"",
		},
		{
			name:        "Supplied/Error",
			give:        &fxevent.Supplied{TypeName: "*bytes.Buffer", Err: someError},
			wantMessage: "\"msg\"=\"error encountered while applying options\" \"error\"=\"some error\" \"type\"=\"*bytes.Buffer\"",
		},
		{
			name: "Provide",
			give: &fxevent.Provided{
				ConstructorName: "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
				Private:         false,
			},
			wantMessage: "\"level\"=0 \"msg\"=\"provided\" \"constructor\"=\"bytes.NewBuffer()\" \"module\"=\"myModule\" \"type\"=\"*bytes.Buffer\"",
		},
		{
			name: "PrivateProvide",
			give: &fxevent.Provided{
				ConstructorName: "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
				Private:         true,
			},
			wantMessage: "\"level\"=0 \"msg\"=\"provided\" \"constructor\"=\"bytes.NewBuffer()\" \"module\"=\"myModule\" \"type\"=\"*bytes.Buffer\" \"private\"=true",
		},
		{
			name:        "Provide/Error",
			give:        &fxevent.Provided{Err: someError},
			wantMessage: "\"msg\"=\"error encountered while applying options\" \"error\"=\"some error\"",
		},
		{
			name: "Replace",
			give: &fxevent.Replaced{
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
			},
			wantMessage: "\"level\"=0 \"msg\"=\"replaced\" \"module\"=\"myModule\" \"type\"=\"*bytes.Buffer\"",
		},
		{
			name:        "Replace/Error",
			give:        &fxevent.Replaced{Err: someError},
			wantMessage: "\"msg\"=\"error encountered while replacing\" \"error\"=\"some error\"",
		},
		{
			name: "Decorate",
			give: &fxevent.Decorated{
				DecoratorName:   "bytes.NewBuffer()",
				ModuleName:      "myModule",
				OutputTypeNames: []string{"*bytes.Buffer"},
			},
			wantMessage: "\"level\"=0 \"msg\"=\"decorated\" \"decorator\"=\"bytes.NewBuffer()\" \"module\"=\"myModule\" \"type\"=\"*bytes.Buffer\"",
		},
		{
			name:        "Decorate/Error",
			give:        &fxevent.Decorated{Err: someError},
			wantMessage: "\"msg\"=\"error encountered while applying options\" \"error\"=\"some error\"",
		},
		{
			name:        "Invoking/Success",
			give:        &fxevent.Invoking{ModuleName: "myModule", FunctionName: "bytes.NewBuffer()"},
			wantMessage: "\"level\"=0 \"msg\"=\"invoking\" \"function\"=\"bytes.NewBuffer()\" \"module\"=\"myModule\"",
		},
		{
			name:        "Invoked/Error",
			give:        &fxevent.Invoked{FunctionName: "bytes.NewBuffer()", Err: someError},
			wantMessage: "\"msg\"=\"invoke failed\" \"error\"=\"some error\" \"stack\"=\"\" \"function\"=\"bytes.NewBuffer()\"",
		},
		{
			name:        "Start/Error",
			give:        &fxevent.Started{Err: someError},
			wantMessage: "\"msg\"=\"start failed\" \"error\"=\"some error\"",
		},
		{
			name:        "Stopping",
			give:        &fxevent.Stopping{Signal: os.Interrupt},
			wantMessage: "\"level\"=0 \"msg\"=\"received signal\" \"signal\"=\"INTERRUPT\"",
		},
		{
			name:        "Stopped/Error",
			give:        &fxevent.Stopped{Err: someError},
			wantMessage: "\"msg\"=\"stop failed\" \"error\"=\"some error\"",
		},
		{
			name:        "RollingBack/Error",
			give:        &fxevent.RollingBack{StartErr: someError},
			wantMessage: "\"msg\"=\"start failed, rolling back\" \"error\"=\"some error\"",
		},
		{
			name:        "RolledBack/Error",
			give:        &fxevent.RolledBack{Err: someError},
			wantMessage: "\"msg\"=\"rollback failed\" \"error\"=\"some error\"",
		},
		{
			name:        "Started",
			give:        &fxevent.Started{},
			wantMessage: "\"level\"=0 \"msg\"=\"started\"",
		},
		{
			name:        "LoggerInitialized/Error",
			give:        &fxevent.LoggerInitialized{Err: someError},
			wantMessage: "\"msg\"=\"custom logger initialization failed\" \"error\"=\"some error\"",
		},
		{
			name:        "LoggerInitialized",
			give:        &fxevent.LoggerInitialized{ConstructorName: "bytes.NewBuffer()"},
			wantMessage: "\"level\"=0 \"msg\"=\"initialized custom fxevent.Logger\" \"function\"=\"bytes.NewBuffer()\"",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			message := ""

			l := funcr.New(
				func(_, args string) {
					message = args
				},
				funcr.Options{},
			)

			WithLogr(&l)().LogEvent(tt.give)

			assert.Equal(t, tt.wantMessage, message)
		})
	}
}
