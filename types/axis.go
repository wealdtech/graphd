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

// Axis defines an axis for a plot.
type Axis struct {
	// Title is the title for the axis.
	Title *Text `json:"title,omitempty"`
	// Type can be 'numeric', 'category' or 'datetime'
	Type string `json:"type,omitempty"`
	// Min is the minimum value for the axis.
	Min interface{} `json:"min,omitempty"`
	// Max is the maximum value for the axis.
	Max interface{} `json:"max,omitempty"`
	// MinorInterval is the interval between axis ticks.
	MinorInterval interface{} `json:"minor_interval,omitempty"`
	// MajorInterval is the interval between axis lines.
	MajorInterval interface{} `json:"major_interval,omitempty"`
	// LabelFormat is the format for axis labels.
	LabelFormat *Text `json:"label_format,omitempty"`
	// ZOrder is the ordering of this axis relative to others.
	ZOrder interface{} `json:"z_order,omitempty"`
}
