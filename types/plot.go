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

// Plot is details of a single plot.
type Plot struct {
	// YAxis is the number of the Y axis in the graph to which this plot applies.
	YAxis interface{} `json:"yaxis,omitempty"`
	// Style can be 'line', 'bar' or 'stackedbar'.
	Style string `json:"style,omitempty"`
	// Zorder provides ordering for multiple plots.
	ZOrder interface{} `json:"z_order,omitempty"`
	// LineStyle can be 'solid', 'dashed' or 'dotted'.
	LineStyle string `json:"line_style,omitempty"`
	// Name is the name of the data series (shows up in legends).
	Name string `json:"name,omitempty"`
	// Width is the relative width of the bar [0,1]
	Width interface{} `json:"width,omitempty"`
	// Color is the color of the plot.
	Color string `json:"color,omitempty"`
	// Data is the data itself.
	Data [][]interface{} `json:"data"`
}
