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

package rest_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	restdaemon "github.com/wealdtech/graphd/services/daemon/rest"
	mockgenerator "github.com/wealdtech/graphd/services/generator/mock"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	generator := mockgenerator.New()

	tests := []struct {
		name    string
		options []restdaemon.Parameter
		err     string
	}{
		{
			name: "GeneratorMissing",
			options: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithListenAddress("localhost:12345"),
			},
			err: "problem with parameters: no generator specified",
		},
		{
			name: "ListenAddressMissing",
			options: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithGenerator(generator),
			},
			err: "problem with parameters: no listen address specified",
		},
		{
			name: "Good",
			options: []restdaemon.Parameter{
				restdaemon.WithLogLevel(zerolog.Disabled),
				restdaemon.WithGenerator(generator),
				restdaemon.WithListenAddress("localhost:12345"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := restdaemon.New(ctx, test.options...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, s)
			}
		})
	}
}
