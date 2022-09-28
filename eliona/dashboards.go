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
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/client"
	"github.com/eliona-smart-building-assistant/go-utils/common"
)

func HailoSmartWasteDashboard() (api.Dashboard, error) {
	dashboard := api.Dashboard{}
	dashboard.Name = "Hailo Smart Waste"
	dashboard.Widgets = []api.Widget{}

	// Process bins
	bins, _, err := client.NewClient().AssetsApi.
		GetAssets(client.AuthenticationContext()).
		AssetTypeName("Hailo FDS Bin").
		Execute()
	if err != nil {
		return api.Dashboard{}, err
	}
	for _, bin := range bins {
		widget := api.Widget{
			WidgetTypeName: "Hailo",
			AssetId:        bin.Id,
			Details: map[string]interface{}{
				"size":     1,
				"timespan": 7,
			},
			Data: []api.WidgetData{
				{
					ElementSequence: nullableInt32(1),
					AssetId:         bin.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "volumepercent",
						"description":         "Level",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(2),
					AssetId:         bin.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "lastclean",
						"description":         "Last cleaning",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(3),
					AssetId:         bin.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "time",
						"description":         "Next cleaning",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(4),
					AssetId:         bin.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "openings",
						"description":         "Openings since last cleaning",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(4),
					AssetId:         bin.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "bat_level",
						"description":         "Battery level",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
			},
		}
		dashboard.Widgets = append(dashboard.Widgets, widget)
	}

	// Process stations
	stations, _, err := client.NewClient().AssetsApi.
		GetAssets(client.AuthenticationContext()).
		AssetTypeName("Hailo FDS Recycling Station").
		Expansions([]string{"Asset.childrenInfo"}).
		Execute()
	if err != nil {
		return api.Dashboard{}, err
	}
	for _, station := range stations {
		widget := api.Widget{
			WidgetTypeName: "Hailo Station",
			AssetId:        station.Id,
			Details: map[string]interface{}{
				"size":     1,
				"timespan": 30,
			},
			Data: []api.WidgetData{
				{
					ElementSequence: nullableInt32(1),
					AssetId:         station.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "volumepercent",
						"description":         "Average level",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
				{
					ElementSequence: nullableInt32(2),
					AssetId:         station.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "volumepercent",
						"description":         "Complete",
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				},
			},
		}

		// add child bins to widget
		if station.ChildrenInfo != nil {
			for _, childBin := range station.ChildrenInfo {
				childData := api.WidgetData{
					ElementSequence: nullableInt32(2),
					AssetId:         childBin.Id,
					Data: map[string]interface{}{
						"aggregatedDataField": nil,
						"aggregatedDataType":  "heap",
						"attribute":           "volumepercent",
						"description":         childBin.Description,
						"key":                 "",
						"seq":                 0,
						"subtype":             "input",
					},
				}
				widget.Data = append(widget.Data, childData)
			}
		}

		// add station widget to dashboard
		dashboard.Widgets = append(dashboard.Widgets, widget)
	}
	return dashboard, nil
}

func nullableInt32(val int32) api.NullableInt32 {
	return *api.NewNullableInt32(common.Ptr[int32](val))
}
