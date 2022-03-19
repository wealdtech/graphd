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

package matplotlib_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/graphd/services/generator/matplotlib"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		options []matplotlib.Parameter
		err     string
	}{
		{
			name: "ResourcesDirMissing",
			options: []matplotlib.Parameter{
				matplotlib.WithLogLevel(zerolog.Disabled),
			},
			err: "problem with parameters: no resources directory specified",
		},
		{
			name: "Good",
			options: []matplotlib.Parameter{
				matplotlib.WithLogLevel(zerolog.Disabled),
				matplotlib.WithResourcesDir("dir"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := matplotlib.New(ctx, test.options...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, s)
			}
		})
	}
}
