//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package apiservices

import (
	"context"
	"hailo/apiserver"
	"hailo/conf"
	"net/http"
)

// ConfigurationApiService is a service that implements the logic for the ConfigurationApiServicer
// This service should implement the business logic for every endpoint for the ConfigurationApi API.
// Include any external packages or services that will be required by this service.
type ConfigurationApiService struct {
}

// NewConfigurationApiService creates a default api service
func NewConfigurationApiService() apiserver.ConfigurationApiServicer {
	return &ConfigurationApiService{}
}

// DeleteConfigurationById - Deletes a FDS endpoint
func (s *ConfigurationApiService) DeleteConfigurationById(ctx context.Context, configId int64) (apiserver.ImplResponse, error) {
	count, err := conf.DeleteConfig(ctx, configId)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	if count == 0 {
		return apiserver.ImplResponse{Code: http.StatusNotFound}, err
	}
	return apiserver.ImplResponse{Code: http.StatusNoContent}, err
}

// GetConfigurationById - Get FDS endpoint
func (s *ConfigurationApiService) GetConfigurationById(ctx context.Context, configId int64) (apiserver.ImplResponse, error) {
	config, err := conf.GetConfig(context.Background(), configId)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	if config == nil {
		return apiserver.ImplResponse{Code: http.StatusNotFound}, err
	}
	return apiserver.Response(http.StatusOK, config), nil
}

// GetConfigurations - Get all FDS endpoints
func (s *ConfigurationApiService) GetConfigurations(ctx context.Context) (apiserver.ImplResponse, error) {
	configs, err := conf.GetConfigs(ctx)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusOK, configs), nil
}

// PostConfiguration - Creates an FDS endpoint
func (s *ConfigurationApiService) PostConfiguration(ctx context.Context, config apiserver.Configuration) (apiserver.ImplResponse, error) {
	insertedConfig, err := conf.InsertConfig(ctx, config)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, insertedConfig), nil
}

// PutConfigurationById - Upserts an FDS endpoint
func (s *ConfigurationApiService) PutConfigurationById(ctx context.Context, configId int64, config apiserver.Configuration) (apiserver.ImplResponse, error) {
	upsertedConfig, err := conf.UpsertConfigById(ctx, configId, config)
	if err != nil {
		return apiserver.ImplResponse{Code: http.StatusInternalServerError}, err
	}
	return apiserver.Response(http.StatusCreated, upsertedConfig), nil
}
