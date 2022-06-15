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
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
)

const DefaultInactiveTimeout = 60 * 60 * 24 // time until set a container to inactive (sec)

type FdsConfig struct {
	Name       string `json:"username"`
	Password   string `json:"password"`
	FdsServer  string `json:"fds_server"`
	AuthServer string `json:"auth_server"`
}

type Config struct {
	Id              int
	FdsConfig       FdsConfig
	Enable          bool
	Active          bool
	IntervalSec     int
	AuthTimeout     int
	RequestTimeout  int
	InactiveTimeout int
	ProjectIds      []string
}

// GetConfigs reads all configured endpoints to a Hailo Digital Hub
func GetConfigs() []Config {
	channel := make(chan Config)
	var configurations []Config
	go func() {
		_ = db.Query(db.Pool(), "select app_id, config, enable, active, interval_sec, auth_timeout, request_timeout, inactive_timeout, proj_ids from hailo.config", channel)
	}()
	for configuration := range channel {
		if configuration.InactiveTimeout == 0 {
			configuration.InactiveTimeout = DefaultInactiveTimeout
		}
		configurations = append(configurations, configuration)
	}
	return configurations
}

// BuildFdsConfig create a config object with the given parameters and default values
func BuildFdsConfig(authServer string, username string, password string, fdsEndpoint string) Config {
	config := Config{
		FdsConfig: FdsConfig{
			Name:       username,
			Password:   password,
			FdsServer:  fdsEndpoint,
			AuthServer: authServer,
		},
		Enable:         true,
		AuthTimeout:    5,
		RequestTimeout: 60,
	}
	return config
}

func GetAssetId(confId int, projId string, deviceId string) *int {
	assetId, err := db.QuerySingleRow[*int](db.Pool(), "select asset_id from hailo.asset where config_id = $1 and proj_id = $2 and device_id = $3", confId, projId, deviceId)
	if err != nil {
		log.Error("Hailo", "Error getting asset id: %v", err)
	}
	return assetId
}

func InsertAsset(confId int, projId string, deviceId string, assetId int) error {
	return db.Exec(db.Pool(), "insert into hailo.asset (config_id, device_id, proj_id, asset_id) values ($1, $2, $3, $4)",
		confId,
		deviceId,
		projId,
		assetId)
}

func SetConfigActive(appId int, state bool) {
	err := db.Exec(db.Pool(), "update hailo.config set active = $1 where app_id = $2", state, appId)
	if err != nil {
		log.Error("Hailo", "Error during setting configuration inactive: %v", err)
	}
}

func SetAllConfigsInactive() {
	err := db.Exec(db.Pool(), "update hailo.config set active = false")
	if err != nil {
		log.Error("Hailo", "Error during setting configuration inactive: %v", err)
	}
}
