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
	"github.com/eliona-smart-building-assistant/go-eliona/assets"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
)

const (
	BinAssetType              = "Hailo FDS Bin"
	RecyclingStationAssetType = "Hailo FDS Recycling Station"
	DigitalHubAssetType       = "Hailo Digital Hub"
)

// InitAssetTypes creates the necessary asset types
func InitAssetTypes(connection db.Connection) error {

	err := assets.UpsertAssetType(connection, assets.AssetType{
		// Hailo FDS Bin               | f      | Hailo             | {"de": "FDS Web-API from Halio", "en": "Recycling station using FDS Web-API from Halio"}                        | eliona.de | trash
		Id:               BinAssetType,
		Custom:           false,
		Vendor:           "Hailo",
		Translation:      &assets.Translation{German: "FDS Web-API from Hailo", English: "Recycling station using FDS Web-API from Hailo"},
		DocumentationUrl: "https://www.hailodigitalhub.de/",
		Icon:             "trash",
		Attributes: []assets.AssetTypeAttribute{
			// Hailo FDS Bin               | device-status   | openings      |         | t      | {"de": "Öffnungen seit Leerung", "en": "Openings since emptying"}      | {device_type_specific,last_empty_count}       |
			{
				Id:             "openings",
				AttributeType:  "device-status",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Öffnungen seit Leerung", English: "Openings since emptying"},
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
				Precision:      common.Ptr(int16(0)),
			},
			// Hailo FDS Bin               | battery-voltage | bat_level     |         | t      | {"de": "Batteriestand", "en": "Battery Level"}                       | {device_type_specific,battery_level}         | %
			{
				Id:             "bat_level",
				AttributeType:  "battery-voltage",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Batteriestand", English: "Battery Level"},
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo FDS Bin               | device-info     | volume        | info    | t      | {"de": "Gesamtvolumen", "en": "Total Volume"}                 | {device_type_specific,bin_volume}               | l
			{
				Id:            "volume",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Gesamtvolumen", English: "Total Volume"},
				Unit:          "l",
			},
			// Hailo FDS Bin               | device-status   | alarm         |         | t      | {"de": "Alarm", "en": "Alarm"}                                        | {device_type_specific,bin_alarm}       |
			{
				Id:            "alarm",
				AttributeType: "device-status",
				Subtype:       assets.InputSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Alarm", English: "Alarm"},
			},
			// Hailo FDS Bin               | device-status   | totalopenings |         | t      | {"de": "Tolal Öffnungen", "en": "Total Openings"}                      | {device_type_specific,inputs_count}       |
			{
				Id:             "totalopenings",
				AttributeType:  "device-status",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Total Öffnungen", English: "Total Openings"},
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo FDS Bin               | level           | volumepercent |         | t      | {"de": "Füllstand", "en": "Fill Level"}                                  | {device_type_specific,filling_level,level}     | %
			{
				Id:             "volumepercent",
				AttributeType:  "level",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Füllstand", English: "Fill Level"},
				Unit:           "%",
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo FDS Bin               | device-info     | exp_percent   | status  | t      | {"de": "Erwarteter Füllstand Leerung", "en": "Expected fill level emptying"} | {device_type_specific,expected_filling_level} | %
			{
				Id:            "exp_percent",
				AttributeType: "device-info",
				Subtype:       assets.StatusSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Erwarteter Füllstand Leerung", English: "Expected fill level emptying"},
				Unit:          "%",
			},
			// Hailo FDS Bin               | device-info     | reg_date      | info    | t      | {"de": "Registrationsdatum", "en": "Registration Date"}        | {generic,registration_date}              |
			{
				Id:            "reg_date",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Registrationsdatum", English: "Registration Date"},
			},
			// Hailo FDS Bin               | device-info     | lastclean     |         | t      | {"de": "Letzte Leerung", "en": "Last emptying"}                 | {device_type_specific,last_service}              | d
			{
				Id:            "lastclean",
				AttributeType: "device-info",
				Subtype:       assets.InputSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Letzte Leerung", English: "Last emptying"},
				Unit:          "d",
				Precision:     common.Ptr(int16(0)),
			},
			// Hailo FDS Bin               | device-info     | time          |         | t      | {"de": "Nächste Leerung", "en": "Next emptying"}             | {device_type_specific,expected_next_service}                | d
			{
				Id:            "time",
				AttributeType: "device-info",
				Subtype:       assets.InputSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Nächste Leerung", English: "Next emptying"},
				Unit:          "d",
				Precision:     common.Ptr(int16(0)),
			},
			// Hailo FDS Bin               | device-info     | last_contact  |         | t      | {"de": "Letzter Kontakt", "en": "Last Contact"}            | {generic,last_contact}                  | h
			{
				Id:            "last_contact",
				AttributeType: "device-info",
				Subtype:       assets.InputSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Letzter Kontakt", English: "Last Contact"},
				Unit:          "h",
				Precision:     common.Ptr(int16(0)),
			},
		},
	})
	if err != nil {
		return err
	}

	// Hailo FDS Recycling Station | f      | Hailo             | {"de": "Recycling Station über FDS Web-API from Halio", "en": "Recycling station using FDS Web-API from Halio"} | eliona.de | trash
	err = assets.UpsertAssetType(connection, assets.AssetType{
		Id:               RecyclingStationAssetType,
		Custom:           false,
		Vendor:           "Hailo",
		Translation:      &assets.Translation{German: "Recycling Station über FDS Web-API from Hailo", English: "Recycling station using FDS Web-API from Hailo"},
		DocumentationUrl: "https://www.hailodigitalhub.de/",
		Icon:             "trash",
		Attributes: []assets.AssetTypeAttribute{
			// Hailo FDS Recycling Station | device-info     | reg_date      | info    | t      | {"de": "Registrationsdatum", "en": "Registration Date"}        | {generic,registration_date}               |
			{
				Id:            "reg_date",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Registrationsdatum", English: "Registration Date"},
			},
			// Hailo FDS Recycling Station | battery-voltage | bat_level     |         | t      | {"de": "Durchschnittlicher Batteriestand", "en": "Average Battery Level"}  | {device_type_specific,average_battery_level}   | %
			{
				Id:             "bat_level",
				AttributeType:  "battery-voltage",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Durchschnittlicher Batteriestand", English: "Average Battery Level"},
				Unit:           "%",
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo FDS Recycling Station | device-info     | last_contact  |         | t      | {"de": "Letzter Kontakt", "en": "Last Contact"}           | {generic,last_contact}                   | h
			{
				Id:            "last_contact",
				AttributeType: "device-info",
				Subtype:       assets.InputSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Letzter Kontakt", English: "Last Contact"},
				Unit:          "h",
				Precision:     common.Ptr(int16(0)),
			},
			// Hailo FDS Recycling Station | device-info     | volume        | info    | t      | {"de": "Kombiniertes Gesamtvolumen", "en": "Total Combined Volume"}    | {device_type_specific,total_combined_volume}      | l
			{
				Id:            "volume",
				AttributeType: "device-info",
				Subtype:       assets.InputSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Kombiniertes Gesamtvolumen", English: "Total Combined Volume"},
				Unit:          "l",
			},
			// Hailo FDS Recycling Station | device-status   | totalopenings |         | t      | {"de": "Total Öffnungen kombiniert", "en": "Total Combined Openings"}    | {device_type_specific,total_inputs_count}    |
			{
				Id:             "totalopenings",
				AttributeType:  "device-status",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Total Öffnungen kombiniert", English: "Total Combined Openings"},
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo FDS Recycling Station | level           | volumepercent |         | t      | {"de": "Durchschnittlicher Füllstand", "en": "Average Fill Level"}    | {device_type_specific,average_filling_level}        | %
			{
				Id:             "volumepercent",
				AttributeType:  "level",
				Subtype:        assets.InputSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Durchschnittlicher Füllstand", English: "Average Fill Level"},
				Unit:           "%",
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
		},
	})
	if err != nil {
		return err
	}

	// Hailo Digital Hub           | f      | Hailo Digital Hub | {"de": "The Web-API from Hailo Digital Hub"}                                                                    |          | trash
	err = assets.UpsertAssetType(connection, assets.AssetType{
		Id:               DigitalHubAssetType,
		Custom:           false,
		Vendor:           "Hailo",
		Translation:      &assets.Translation{German: "Web-API von Hailo Digital Hub", English: "The Web-API from Hailo Digital Hub"},
		DocumentationUrl: "https://www.hailodigitalhub.de/",
		Icon:             "trash",
		Attributes: []assets.AssetTypeAttribute{
			// Hailo Digital Hub           | device-info     | lastclean     |         | t      |                                       | {lastclean}                                       |
			{
				Id:            "lastclean",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
			},
			// Hailo Digital Hub           | device-info     | time          |         | t      |                                       | {forecast,time}                                       |
			{
				Id:            "time",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
			},
			// Hailo Digital Hub           | device-info     | rssi          | rfstat  | t      |                                           | {rssi}                                   |
			{
				Id:            "rssi",
				AttributeType: "device-info",
				Subtype:       assets.StatusSubtype,
				Enable:        true,
			},
			// Hailo Digital Hub           | battery-voltage | voltage       | status  | t      | {"de": "Batteriespannung", "en": "Battery Voltage"}      | {voltage}                    | V
			{
				Id:             "voltage",
				AttributeType:  "battery-voltage",
				Subtype:        assets.InfoSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Batteriespannung", English: "Battery Voltage"},
				Unit:           "V",
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
				Precision:      common.Ptr(int16(2)),
			},
			// Hailo Digital Hub           | device-info     | volume        |         | t      | {"de": "Volumen", "en": "Volume"}                | {volume}                            |
			{
				Id:            "volume",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Volumen", English: "Volume"},
				Precision:     common.Ptr(int16(2)),
			},
			// Hailo Digital Hub           | level           | volumepercent |         | t      | {"de": "Volumen", "en": "Volume"}          | {volumepercent}                                  | %
			{
				Id:             "volumepercent",
				AttributeType:  "level",
				Subtype:        assets.InfoSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Volumen", English: "Volume"},
				Unit:           "%",
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo Digital Hub           | device-status   | totalopenings |         | t      | {"de": "Öffnungen", "en": "Openings"}      | {totalopenings}                                  |
			{
				Id:             "totalopenings",
				AttributeType:  "device-status",
				Subtype:        assets.InfoSubtype,
				Enable:         true,
				Translation:    &assets.Translation{German: "Öffnungen", English: "Openings"},
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo Digital Hub           | device-info     | percent       |         | t      |                              | {forecast,percent}                                                |
			{
				Id:             "percent",
				AttributeType:  "device-info",
				Subtype:        assets.InfoSubtype,
				Enable:         true,
				PipelineMode:   assets.AveragePipelineMode,
				PipelineRaster: "{M15,H1,DAY}",
			},
			// Hailo Digital Hub           | device-info     | closed        |         | t      | {"de": "Geschlossen", "en": "Closed"}  | {closed}                                      |
			{
				Id:            "closed",
				AttributeType: "device-info",
				Subtype:       assets.InfoSubtype,
				Enable:        true,
				Translation:   &assets.Translation{German: "Geschlossen", English: "Closed"},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
