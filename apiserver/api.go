/*
 * Hailo app API
 *
 * API to access and configure the Hailo app
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

import (
	"context"
	"net/http"
)

// AssetMappingApiRouter defines the required methods for binding the api requests to a responses for the AssetMappingApi
// The AssetMappingApiRouter implementation should parse necessary information from the http request,
// pass the data to a AssetMappingApiServicer to perform the required actions, then write the service results to the http response.
type AssetMappingApiRouter interface {
	GetAssetMappingsByConfig(http.ResponseWriter, *http.Request)
}

// ConfigurationApiRouter defines the required methods for binding the api requests to a responses for the ConfigurationApi
// The ConfigurationApiRouter implementation should parse necessary information from the http request,
// pass the data to a ConfigurationApiServicer to perform the required actions, then write the service results to the http response.
type ConfigurationApiRouter interface {
	GetConfiguration(http.ResponseWriter, *http.Request)
	GetConfigurations(http.ResponseWriter, *http.Request)
	PutConfiguration(http.ResponseWriter, *http.Request)
}

// AssetMappingApiServicer defines the api actions for the AssetMappingApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type AssetMappingApiServicer interface {
	GetAssetMappingsByConfig(context.Context, int32) (ImplResponse, error)
}

// ConfigurationApiServicer defines the api actions for the ConfigurationApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ConfigurationApiServicer interface {
	GetConfiguration(context.Context, int32) (ImplResponse, error)
	GetConfigurations(context.Context) (ImplResponse, error)
	PutConfiguration(context.Context, int32) (ImplResponse, error)
}
