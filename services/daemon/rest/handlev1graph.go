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
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wealdtech/graphd/types"
)

// handleV1Graph handles version 1 of the graph call.
func (s *Service) handleV1Graph(c *gin.Context) {
	ctx := context.Background()

	body := &types.Graph{}
	if err := c.ShouldBindJSON(body); err != nil {
		log.Debug().Err(err).Msg("failed to read body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if e := log.Trace(); e.Enabled() {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to display decoded request body")
			return
		}
		e.Str("body", string(bodyBytes)).Msg("Decoded body")
	}

	graph, err := s.generator.Graph(ctx, body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to generate graph")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate graph"})
		return
	}
	if graph == nil {
		c.Status(http.StatusNoContent)
		return
	}

	bytes, err := io.ReadAll(graph)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read graph response")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate graph"})
	}

	format := "png"
	if body.Meta != nil && body.Meta.Format != "" {
		format = strings.ToLower(body.Meta.Format)
	}
	switch format {
	case "pdf":
		c.Data(http.StatusOK, "application/pdf", bytes)
	case "png":
		c.Data(http.StatusOK, "image/png", bytes)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to provide graph"})
	}
}
