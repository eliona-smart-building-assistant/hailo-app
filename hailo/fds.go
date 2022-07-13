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

package hailo

import (
	"encoding/base64"
	"encoding/json"
	"github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"hailo/conf"
	"strings"
	"sync"
	"time"
)

const (
	AuthApiPath          = "/beta/v1/authentication"
	FdsSpecificationPath = "/specifications"
	FdsStatusPath        = "/status"
	FdsDiagnosticsPath   = "/diagnostics"
	FdsIdParam           = "/?ids="
)

type auth struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type jwtClaims struct {
	IssuedAt       int64 `json:"iat"`
	ExpirationTime int64 `json:"exp"`
}

type Specs struct {
	Data []Spec `json:"data"`
}

type Spec struct {
	DeviceId string `json:"device_id"`
	Generic  struct {
		DeviceType       string `json:"device_type"`
		RegistrationDate string `json:"registration_date"`
		DeviceSerial     string `json:"device_serial"`
		Manufacturer     string `json:"manufacturer"`
		Model            string `json:"model"`
	} `json:"generic"`
	DeviceTypeSpecific struct {
		BinVolume           int    `json:"bin_volume"`
		Category            string `json:"category"`
		Channel             string `json:"channel"`
		ComponentIdList     []Spec `json:"component_id_list"`
		TotalCombinedVolume int    `json:"total_combined_volume"`
		ContentCategory     string `json:"content_category"`
	} `json:"device_type_specific"`
}

type Statuses struct {
	Data []Status `json:"data"`
}

type Status struct {
	DeviceId string `json:"device_id"`
	Generic  struct {
		LastContact string `json:"last_contact"`
	} `json:"generic"`
	DeviceTypeSpecific struct {
		// Single Container
		BatteryLevel   float32 `json:"battery_level"`
		InputCount     int     `json:"inputs_count"`
		LastEmptyCount int     `json:"last_empty_count"`
		BinAlarm       bool    `json:"bin_alarm"`
		// Station
		AverageBatteryLevel float32  `json:"average_battery_level"`
		AverageFillingLevel float32  `json:"average_filling_level"`
		TotalInputsCount    int      `json:"total_inputs_count"`
		CompStatuses        []Status `json:"component_statuses"`
		// Both
		FillingLevel []struct {
			Level float32 `json:"level"`
		} `json:"filling_level"`
	} `json:"device_type_specific"`
}

type Diags struct {
	Data []Diag `json:"data"`
}

type Diag struct {
	DeviceId string `json:"device_id"`
	Generic  struct {
		LastService string `json:"last_service"`
	} `json:"generic"`
	DeviceTypeSpecific struct {
		// Single Container
		ExpectedNextService  string  `json:"expected_next_service"`
		ExpectedFillingLevel float32 `json:"expected_filling_level"`
		// Station
		StationExpectedNextService  string  `json:"station_expected_next_service"`
		AverageExpectedFillingLevel float32 `json:"average_expected_filling_level"`
	} `json:"device_type_specific"`
	ComponentDiagnostics []Diag `json:"component_diagnostics"`
}

// tokens holds generated tokens for further use until they come invalid
var tokens sync.Map

// GetSpecs reads the specification for all Hailo smart devices from eliona endpoint
func GetSpecs(configuration conf.Config) (Specs, error) {
	request, err := http.NewRequestWithBearer(
		configuration.FdsConfig.FdsServer+FdsSpecificationPath,
		getToken(configuration),
	)
	if err != nil {
		return Specs{}, err
	}

	specs, err := http.Read[Specs](request, time.Duration(configuration.RequestTimeout)*time.Second, true)
	if err != nil {
		return specs, err
	}

	return specs, nil
}

// GetDiag reads the diagnostic data for the given device id
func GetDiag(configuration conf.Config, deviceId string) (Diag, error) {
	request, err := http.NewRequestWithBearer(
		configuration.FdsConfig.FdsServer+FdsDiagnosticsPath+FdsIdParam+deviceId,
		getToken(configuration),
	)
	if err != nil {
		return Diag{}, err
	}

	diagnostics, err := http.Read[Diags](request, time.Duration(configuration.RequestTimeout)*time.Second, true)
	if err != nil {
		return Diag{}, err
	}
	return diagnostics.Data[0], nil
}

// GetStatus reads the status data for the given device id
func GetStatus(configuration conf.Config, deviceId string) (Status, error) {
	request, err := http.NewRequestWithBearer(
		configuration.FdsConfig.FdsServer+FdsStatusPath+FdsIdParam+deviceId,
		getToken(configuration),
	)
	if err != nil {
		return Status{}, err
	}

	statuses, err := http.Read[Statuses](request, time.Duration(configuration.RequestTimeout)*time.Second, true)
	if err != nil {
		return Status{}, err
	}

	return statuses.Data[0], nil
}

// IsStation returns true, if the status is from a station. A station contains multiple component statuses.
func (status Status) IsStation() bool {
	return len(status.DeviceTypeSpecific.CompStatuses) > 0
}

// isTokenValid checks if the given token is valid
func isTokenValid(token string) bool {
	currentTime := time.Now().Unix()

	start := strings.Index(token, ".")
	end := strings.LastIndex(token, ".")

	jwtBody := token[start+1 : end]
	jwtBody = decodeBase64(jwtBody)
	var claims jwtClaims
	err := json.Unmarshal([]byte(jwtBody), &claims)
	if err != nil {
		return false
	}

	if currentTime+240 > claims.ExpirationTime {
		log.Info("Hailo", "Token expired")
		return false
	}

	if currentTime+5 < claims.IssuedAt {
		log.Info("Hailo", "Seems the token will issue in the feature! :D")
		return false
	}

	return true
}

// getToken creates a new token or delivers a previous token until this token is valid
func getToken(configuration conf.Config) string {
	token, found := tokens.Load(configuration.Id)
	if found {
		if isTokenValid(token.(string)) {
			return token.(string)
		}
	}
	token, _ = authenticate(configuration)
	tokens.Store(configuration.Id, token)
	return token.(string)
}

// authenticate creates a new token
func authenticate(configuration conf.Config) (string, error) {

	log.Info("Hailo", "Create new Authentication token")
	request, err := http.NewPostRequest(
		configuration.FdsConfig.AuthServer+AuthApiPath,
		auth{
			UserName: configuration.FdsConfig.Name,
			Password: configuration.FdsConfig.Password,
		},
	)
	if err != nil {
		return "", err
	}

	token, err := http.Do(request, time.Duration(configuration.AuthTimeout)*time.Second, true)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(token), "\"", ""), nil
}

func decodeBase64(b64 string) string {
	plain, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return ""
	}
	return string(plain)
}

func encodeBase64(plain string) string {
	return base64.StdEncoding.EncodeToString([]byte(plain))
}
