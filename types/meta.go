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

package types

// Meta provides metadata for a graph.
type Meta struct {
	Theme  string  `json:"theme,omitempty"`
	Format string  `json:"format,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
	DPI    float64 `json:"dpi,omitempty"`
}
