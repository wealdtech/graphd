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
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"github.com/wealdtech/graphd/services/generator"
)

// Service is basic reporter service.
type Service struct {
	generator     generator.Service
	listenAddress string
}

// module-wide log.
var log zerolog.Logger

// New creates a new service.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "daemon").Str("impl", "rest").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	s := &Service{
		generator:     parameters.generator,
		listenAddress: parameters.listenAddress,
	}

	// Start up the daemon.
	if err := s.startDaemon(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to start daemon")
	}

	return s, nil
}
