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

package generator

import (
	"context"
	"io"

	"github.com/wealdtech/graphd/types"
)

// Service defines the interface for a graph generator service.
type Service interface {
	// Graph generates a graph from its definition.
	Graph(ctx context.Context, graph *types.Graph) (io.Reader, error)
}
