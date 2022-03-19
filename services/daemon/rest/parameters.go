// Copyright © 2022 Weald Technology Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rest

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/wealdtech/graphd/services/generator"
)

type parameters struct {
	logLevel      zerolog.Level
	generator     generator.Service
	listenAddress string
}

// Parameter is the interface for service parameters.
type Parameter interface {
	apply(*parameters)
}

type parameterFunc func(*parameters)

func (f parameterFunc) apply(p *parameters) {
	f(p)
}

// WithLogLevel sets the log level for the module.
func WithLogLevel(logLevel zerolog.Level) Parameter {
	return parameterFunc(func(p *parameters) {
		p.logLevel = logLevel
	})
}

// WithGenerator sets the generator for this module.
func WithGenerator(generator generator.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.generator = generator
	})
}

// WithListenAddress sets the listen address for this module.
func WithListenAddress(address string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.listenAddress = address
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel: zerolog.GlobalLevel(),
	}
	for _, p := range params {
		if params != nil {
			p.apply(&parameters)
		}
	}

	if parameters.generator == nil {
		return nil, errors.New("no generator specified")
	}
	if parameters.listenAddress == "" {
		return nil, errors.New("no listen address specified")
	}

	return &parameters, nil
}
