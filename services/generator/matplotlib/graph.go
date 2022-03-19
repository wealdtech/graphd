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

package matplotlib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wealdtech/graphd/types"
)

// Graph generates a raph from its definition.
func (s *Service) Graph(ctx context.Context, in *types.Graph) (io.Reader, error) {
	err := validate(in)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	script := filepath.Join(s.resourcesDir, "scripts", "graph.py")
	log.Trace().Str("path", script).Msg("Script path")

	// Create a temporary directory to store the input and output files.
	dir, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	inFile := filepath.Join(dir, "input.json")
	if err := os.WriteFile(inFile, data, 0600); err != nil {
		return nil, err
	}

	outBase := filepath.Join(dir, "output")

	args := make([]string, 0)
	args = append(args, inFile)
	args = append(args, outBase)
	log.Trace().Strs("args", args).Msg("Parameters")
	cmd := exec.CommandContext(ctx, script, args...)
	if err := cmd.Run(); err != nil {
		if exitError, isExitError := err.(*exec.ExitError); isExitError {
			if len(exitError.Stderr) > 0 {
				stderr := make([]string, 0, len(exitError.Stderr))
				for _, line := range exitError.Stderr {
					stderr = append(stderr, string(line))
				}
				log.Warn().Strs("stderr", stderr).Msg("Errors")
			}
		}
		return nil, err
	}

	format := "png"
	if in.Meta != nil && in.Meta.Format != "" {
		format = strings.ToLower(in.Meta.Format)
	}
	outFile := fmt.Sprintf("%s.%s", outBase, format)
	output, err := os.ReadFile(outFile)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to read output")
		return nil, err
	}

	return bytes.NewReader(output), nil
}

// validate validates the graph structure.
func validate(graph *types.Graph) error {
	if graph == nil {
		return errors.New("no graph")
	}

	if len(graph.Plots) == 0 {
		return errors.New("no data")
	}

	for i := range graph.Plots {
		if len(graph.Plots[i].Data) == 0 {
			return fmt.Errorf("no data for plot %d", i+1)
		}

		for j := range graph.Plots[i].Data {
			if len(graph.Plots[i].Data[j]) != 2 {
				return fmt.Errorf("point %d of plot %d should have two co-ordinates; has %d", j+1, i+1, len(graph.Plots[i].Data[j]))
			}
		}
	}

	return nil
}
