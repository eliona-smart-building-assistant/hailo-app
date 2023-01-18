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
	"net/http"
	"strings"
)

// AssetMappingApiController binds http requests to an api service and writes the service results to the http response
type AssetMappingApiController struct {
	service      AssetMappingApiServicer
	errorHandler ErrorHandler
}

// AssetMappingApiOption for how the controller is set up.
type AssetMappingApiOption func(*AssetMappingApiController)

// WithAssetMappingApiErrorHandler inject ErrorHandler into controller
func WithAssetMappingApiErrorHandler(h ErrorHandler) AssetMappingApiOption {
	return func(c *AssetMappingApiController) {
		c.errorHandler = h
	}
}

// NewAssetMappingApiController creates a default api controller
func NewAssetMappingApiController(s AssetMappingApiServicer, opts ...AssetMappingApiOption) Router {
	controller := &AssetMappingApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the AssetMappingApiController
func (c *AssetMappingApiController) Routes() Routes {
	return Routes{
		{
			"GetAssetMappings",
			strings.ToUpper("Get"),
			"/v1/asset-mappings",
			c.GetAssetMappings,
		},
	}
}

// GetAssetMappings - List all mapped assets
func (c *AssetMappingApiController) GetAssetMappings(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	configIdParam, err := parseInt64Parameter(query.Get("configId"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.GetAssetMappings(r.Context(), configIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
