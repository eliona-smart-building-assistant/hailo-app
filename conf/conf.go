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

package conf

import (
	"context"
	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"hailo/apiserver"
	dbhailo "hailo/db/hailo"
)

const DefaultInactiveTimeout = 60 * 60 * 24 // time until set a container to inactive (sec)

type FdsConfig struct {
	Name       string `json:"username"`
	Password   string `json:"password"`
	FdsServer  string `json:"fds_server"`
	AuthServer string `json:"auth_server"`
}

// GetConfigs reads all configured endpoints for a Hailo Digital Hub
func GetConfigs(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := dbhailo.Configs().All(ctx, db.Database(app.AppName()))
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		apiConfigs = append(apiConfigs, *apiConfigFromDbConfig(dbConfig))
	}
	return apiConfigs, nil
}

func GetAssetMappings(ctx context.Context, configId int64) ([]apiserver.AssetMapping, error) {
	var mods []qm.QueryMod
	if configId > 0 {
		mods = append(mods, dbhailo.AssetWhere.ConfigID.EQ(configId))
	}
	dbAssetMappings, err := dbhailo.Assets(mods...).All(ctx, db.Database(app.AppName()))
	if err != nil {
		return nil, err
	}
	var apiAssetMappings []apiserver.AssetMapping
	for _, dbAssetMapping := range dbAssetMappings {
		apiAssetMappings = append(apiAssetMappings, *apiAssetMappingFromDbAssetMapping(dbAssetMapping))
	}
	return apiAssetMappings, nil
}

// InsertConfig inserts or updates
func InsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig := dbConfigFromApiConfig(&config)
	err := dbConfig.Insert(ctx, db.Database(app.AppName()), boil.Blacklist(dbhailo.ConfigColumns.AppID))
	if err != nil {
		return apiserver.Configuration{}, err
	}
	config.Id = &dbConfig.AppID
	return config, err
}

// UpsertConfigById inserts or updates
func UpsertConfigById(ctx context.Context, configId int64, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig := dbConfigFromApiConfig(&config)
	dbConfig.AppID = configId
	err := dbConfig.Upsert(ctx, db.Database(app.AppName()), true,
		[]string{dbhailo.ConfigColumns.AppID},
		boil.Blacklist(dbhailo.ConfigColumns.AppID),
		boil.Infer(),
	)
	config.Id = &dbConfig.AppID
	return config, err
}

func apiAssetMappingFromDbAssetMapping(dbAssetMapping *dbhailo.Asset) *apiserver.AssetMapping {
	var apiAssetMapping apiserver.AssetMapping
	apiAssetMapping.AssetId = dbAssetMapping.AssetID
	apiAssetMapping.DeviceId = dbAssetMapping.DeviceID
	apiAssetMapping.ConfigId = int32(dbAssetMapping.ConfigID)
	apiAssetMapping.ProjId = dbAssetMapping.ProjID
	return &apiAssetMapping
}

func apiConfigFromDbConfig(dbConfig *dbhailo.Config) *apiserver.Configuration {
	var apiConfig apiserver.Configuration
	apiConfig.Id = &dbConfig.AppID
	apiConfig.AssetId = dbConfig.AssetID.Ptr()
	apiConfig.Enable = dbConfig.Enable.Ptr()
	apiConfig.Description = dbConfig.Description.Ptr()
	apiConfig.InactiveTimeout = getInactiveTimeout(dbConfig)
	var fdsConfig FdsConfig
	_ = dbConfig.Config.Unmarshal(&fdsConfig)
	apiConfig.Username = &fdsConfig.Name
	apiConfig.Password = &fdsConfig.Password
	apiConfig.AuthServer = &fdsConfig.AuthServer
	apiConfig.FdsServer = &fdsConfig.FdsServer
	apiConfig.AuthTimeout = dbConfig.AuthTimeout
	apiConfig.IntervalSec = dbConfig.IntervalSec
	apiConfig.RequestTimeout = dbConfig.RequestTimeout
	apiConfig.ProjIds = common.Ptr[[]string](dbConfig.ProjIds)
	return &apiConfig
}

func dbConfigFromApiConfig(apiConfig *apiserver.Configuration) *dbhailo.Config {
	var dbConfig dbhailo.Config
	dbConfig.AppID = null.Int64FromPtr(apiConfig.Id).Int64
	dbConfig.AssetID = null.Int32FromPtr(apiConfig.AssetId)
	dbConfig.Enable = null.BoolFromPtr(apiConfig.Enable)
	dbConfig.Description = null.StringFromPtr(apiConfig.Description)
	dbConfig.InactiveTimeout = null.Int32From(apiConfig.InactiveTimeout)
	dbConfig.AuthTimeout = apiConfig.AuthTimeout
	dbConfig.IntervalSec = apiConfig.IntervalSec
	dbConfig.RequestTimeout = apiConfig.RequestTimeout
	if apiConfig.ProjIds != nil {
		dbConfig.ProjIds = *apiConfig.ProjIds
	}
	var fdsConfig types.JSON
	_ = fdsConfig.Marshal(FdsConfig{
		Name:       null.StringFromPtr(apiConfig.Username).String,
		Password:   null.StringFromPtr(apiConfig.Password).String,
		AuthServer: null.StringFromPtr(apiConfig.AuthServer).String,
		FdsServer:  null.StringFromPtr(apiConfig.FdsServer).String,
	})
	dbConfig.Config = fdsConfig
	return &dbConfig
}

func getInactiveTimeout(config *dbhailo.Config) int32 {
	if config.InactiveTimeout.Valid {
		return config.InactiveTimeout.Int32
	} else {
		return DefaultInactiveTimeout
	}
}

// GetConfig reads configured endpoints to a Hailo Digital Hub
func GetConfig(ctx context.Context, configId int64) (*apiserver.Configuration, error) {
	dbConfigs, err := dbhailo.Configs(dbhailo.ConfigWhere.AppID.EQ(configId)).All(ctx, db.Database(app.AppName()))
	if err != nil {
		return nil, err
	}
	if len(dbConfigs) == 0 {
		return nil, err
	}
	return apiConfigFromDbConfig(dbConfigs[0]), nil
}

// DeleteConfig reads configured endpoints to a Hailo Digital Hub
func DeleteConfig(ctx context.Context, configId int64) (int64, error) {
	return dbhailo.Configs(dbhailo.ConfigWhere.AppID.EQ(configId)).DeleteAll(ctx, db.Database(app.AppName()))
}

// BuildFdsConfig create a config object with the given parameters and default values
func BuildFdsConfig(authServer string, username string, password string, fdsEndpoint string) apiserver.Configuration {
	config := apiserver.Configuration{
		Username:       &username,
		Password:       &password,
		FdsServer:      &fdsEndpoint,
		AuthServer:     &authServer,
		Enable:         common.Ptr(true),
		AuthTimeout:    5,
		RequestTimeout: 60,
	}
	return config
}

func GetAssetId(ctx context.Context, config apiserver.Configuration, projId string, deviceId string) (*int32, error) {
	dbAssets, err := dbhailo.Assets(
		dbhailo.AssetWhere.ConfigID.EQ(null.Int64FromPtr(config.Id).Int64),
		dbhailo.AssetWhere.ProjID.EQ(projId),
		dbhailo.AssetWhere.DeviceID.EQ(deviceId),
	).All(ctx, db.Database(app.AppName()))
	if err != nil || len(dbAssets) == 0 {
		return nil, err
	}
	return common.Ptr(dbAssets[0].AssetID), nil
}

func InsertAsset(ctx context.Context, config apiserver.Configuration, projId string, deviceId string, assetId int32) error {
	var dbAsset dbhailo.Asset
	dbAsset.ConfigID = null.Int64FromPtr(config.Id).Int64
	dbAsset.ProjID = projId
	dbAsset.DeviceID = deviceId
	dbAsset.AssetID = assetId
	return dbAsset.Insert(ctx, db.Database(app.AppName()), boil.Infer())
}

func SetConfigActiveState(ctx context.Context, config apiserver.Configuration, state bool) (int64, error) {
	return dbhailo.Configs(
		dbhailo.ConfigWhere.AppID.EQ(null.Int64FromPtr(config.Id).Int64),
	).UpdateAll(ctx, db.Database(app.AppName()), dbhailo.M{
		dbhailo.ConfigColumns.Active: state,
	})
}

func ProjIds(config apiserver.Configuration) []string {
	if config.ProjIds == nil {
		return []string{}
	}
	return *config.ProjIds
}

func IsConfigActive(config apiserver.Configuration) bool {
	return config.Active == nil || *config.Active
}

func IsConfigEnabled(config apiserver.Configuration) bool {
	return config.Enable == nil || *config.Enable
}

func SetAllConfigsInactive(ctx context.Context) (int64, error) {
	return dbhailo.Configs().UpdateAll(ctx, db.Database(app.AppName()), dbhailo.M{
		dbhailo.ConfigColumns.Active: false,
	})
}
