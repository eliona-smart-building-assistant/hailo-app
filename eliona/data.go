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

package eliona

import (
	"context"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"hailo/apiserver"
	"hailo/conf"
	"hailo/hailo"
	"math"
	"time"
)

func UpsertDataForDevices(config apiserver.Configuration, spec hailo.Spec) error {

	for _, projectId := range conf.ProjIds(config) {

		err := upsertDataForDevice(config, projectId, spec)
		if err != nil {
			log.Error("Hailo", "Could not upsert data for device %s: %v", spec.DeviceId, err)
			return err
		}
		for _, subSpec := range spec.DeviceTypeSpecific.ComponentIdList {
			err = upsertDataForDevice(config, projectId, subSpec)
			if err != nil {
				log.Error("Hailo", "Could not upsert data for sub device %s: %v", subSpec.DeviceId, err)
				return err
			}
		}
	}
	return nil
}

type deviceDataPayload struct {
	RegistrationDate string `json:"reg_data"`
	Volume           int    `json:"volume"`
}

func upsertData(subtype api.DataSubtype, time time.Time, assetId int32, payload any) error {
	var statusData api.Data
	statusData.Subtype = subtype
	statusData.Timestamp = *api.NewNullableTime(&time)
	statusData.AssetId = assetId
	statusData.Data = common.StructToMap(payload)
	err := asset.UpsertDataIfAssetExists(statusData)
	if err != nil {
		log.Error("Hailo", "Error during writing data: %v", err)
		return err
	}
	return nil
}

func upsertDataForDevice(config apiserver.Configuration, projectId string, spec hailo.Spec) error {
	log.Debug("Hailo", "Upsert data for device: config %d and device '%s'", config.Id, spec.DeviceId)
	assetId, err := conf.GetAssetId(context.Background(), config, projectId, spec.DeviceId)
	if err != nil {
		return err
	}
	return upsertData(
		api.SUBTYPE_INFO,
		parseTime(spec.Generic.RegistrationDate),
		*assetId,
		deviceDataPayload{RegistrationDate: spec.Generic.RegistrationDate, Volume: binVolume(spec)},
	)
}

func binVolume(spec hailo.Spec) int {
	binVolume := spec.DeviceTypeSpecific.BinVolume
	if spec.DeviceTypeSpecific.TotalCombinedVolume != 0 {
		binVolume = spec.DeviceTypeSpecific.TotalCombinedVolume
	}
	return binVolume
}

type stationDataPayload struct {
	BatteryLevel     int     `json:"bat_level"`
	LastContact      float64 `json:"last_contact"`
	TotalOpenings    int     `json:"totalopenings"`
	VolumePercentage int     `json:"volumepercent"`
	Active           bool    `json:"active"`
}

func UpsertDataForStation(config apiserver.Configuration, status hailo.Status) error {
	for _, projectId := range conf.ProjIds(config) {
		log.Debug("Hailo", "Upsert data for station: config %d and station '%s'", config.Id, status.DeviceId)
		lastContact := parseTimeToHours(status.Generic.LastContact)
		assetId, err := conf.GetAssetId(context.Background(), config, projectId, status.DeviceId)
		if err != nil {
			return err
		}
		err = upsertData(
			api.SUBTYPE_INPUT,
			parseTime(status.Generic.LastContact),
			*assetId,
			stationDataPayload{
				int(status.DeviceTypeSpecific.AverageBatteryLevel * 100),
				lastContact,
				status.DeviceTypeSpecific.TotalInputsCount,
				int(status.DeviceTypeSpecific.AverageFillingLevel * 100),
				CheckActivity(config, lastContact),
			},
		)
		if err != nil {
			log.Error("Hailo", "Could not upsert data for station %s: %v", status.DeviceId, err)
			return err
		}
	}
	return nil
}

func CheckActivity(connection apiserver.Configuration, lastContact float64) bool {
	return lastContact < (float64)(connection.InactiveTimeout/3600)
}

type binDataPayload struct {
	BatteryLevel     int     `json:"bat_level"`
	Openings         int     `json:"openings"`
	LastContact      float64 `json:"last_contact"`
	Alarm            bool    `json:"alarm"`
	TotalOpenings    int     `json:"totalopenings"`
	VolumePercentage int     `json:"volumepercent"`
	Time             float64 `json:"time"`
	LastClean        float64 `json:"lastclean"`
	Active           bool    `json:"active"`
}

type statusDataPayload struct {
	ExpectedPercent int `json:"exp_percent"`
}

func UpsertDataForBin(config apiserver.Configuration, status hailo.Status, diag hailo.Diag) error {
	for _, projectId := range conf.ProjIds(config) {
		log.Debug("Hailo", "Upsert data for bin: config %d and bin '%s'", config.Id, status.DeviceId)
		lastContact := parseTimeToHours(status.Generic.LastContact)
		assetId, err := conf.GetAssetId(context.Background(), config, projectId, status.DeviceId)
		if err != nil {
			return err
		}
		err = upsertData(
			api.SUBTYPE_INPUT,
			parseTime(status.Generic.LastContact),
			*assetId,
			binDataPayload{
				int(status.DeviceTypeSpecific.BatteryLevel * 100),
				status.DeviceTypeSpecific.LastEmptyCount,
				lastContact,
				status.DeviceTypeSpecific.BinAlarm,
				status.DeviceTypeSpecific.InputCount,
				int(status.DeviceTypeSpecific.FillingLevel[0].Level * 100),
				parseTimeToDays(diag.DeviceTypeSpecific.ExpectedNextService),
				parseTimeToDays(diag.Generic.LastService),
				CheckActivity(config, lastContact),
			},
		)
		if err != nil {
			log.Error("Hailo", "Could not upsert data for bin %s: %v", status.DeviceId, err)
			return err
		}
		assetId, err = conf.GetAssetId(context.Background(), config, projectId, status.DeviceId)
		if err != nil {
			return err
		}
		err = upsertData(
			api.SUBTYPE_STATUS,
			parseTime(status.Generic.LastContact),
			*assetId,
			statusDataPayload{int(diag.DeviceTypeSpecific.ExpectedFillingLevel * 100)},
		)
		if err != nil {
			log.Error("Hailo", "Could not upsert data for bin %s: %v", status.DeviceId, err)
			return err
		}
	}

	return nil
}

// 2021-01-26T09:16:16.000Z
func parseTime(iso string) time.Time {
	if iso == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339Nano, iso)
	if err != nil {
		log.Warn("Hailo", "Error while converting ISO 8601 time to unix time %v, in: %s", err, iso)
		return time.Now()
	}
	return t
}

func parseToIsoTime(unix int64) string {
	t := time.Unix(unix, 0)
	return t.UTC().Format("2006-01-02T15:04:05.006Z")
}

func parseTimeToDays(iso string) float64 {
	return math.Round((parseTimeToHours(iso)*100)/24) / 100
}

func parseTimeToHours(iso string) float64 {
	now := time.Now().Unix()
	dst := parseTime(iso).Unix()
	rst := 0.0
	if now > dst {
		rst = (float64)((now - dst) / (60 * 60))
	} else {
		rst = (float64)((dst - now) / (60 * 60))
	}
	return math.Round(rst*100) / 100
}
