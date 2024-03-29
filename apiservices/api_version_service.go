/*
 * App template API
 *
 * API to access and configure the app template
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiservices

import (
	"context"
	"encoding/json"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"hailo/apiserver"
	"net/http"
	"os"
)

// VersionApiService is a service that implements the logic for the VersionApiServicer
// This service should implement the business logic for every endpoint for the VersionApi API.
// Include any external packages or services that will be required by this service.
type VersionApiService struct {
}

// NewVersionApiService creates a default api service
func NewVersionApiService() apiserver.VersionApiServicer {
	return &VersionApiService{}
}

// GetOpenAPI - OpenAPI specification for this API version
func (s *VersionApiService) GetOpenAPI(ctx context.Context) (apiserver.ImplResponse, error) {
	bytes, err := os.ReadFile("openapi.json")
	if err != nil {
		bytes, err = os.ReadFile("apiserver/openapi.json")
		if err != nil {
			log.Error("services", "%s: %v", "GetOpenAPI", err)
			return apiserver.ImplResponse{Code: http.StatusNotFound}, err
		}
	}
	var body interface{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		log.Error("services", "%s: %v", "GetOpenAPI", err)
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, body), nil
}

var BuildTimestamp string // injected during linking, see Dockerfile
var GitCommit string      // injected during linking, see Dockerfile

// GetVersion - Version of the API
func (s *VersionApiService) GetVersion(ctx context.Context) (apiserver.ImplResponse, error) {
	return apiserver.Response(http.StatusOK, common.Ptr(version())), nil
}

func version() map[string]any {
	return map[string]any{
		"timestamp": BuildTimestamp,
		"commit":    GitCommit,
	}
}
