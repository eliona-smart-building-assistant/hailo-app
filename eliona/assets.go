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
	"fmt"
	api "github.com/eliona-smart-building-assistant/go-eliona-api-client"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"hailo/conf"
	"hailo/hailo"
)

const (
	BinAssetType              = "Hailo FDS Bin"
	RecyclingStationAssetType = "Hailo FDS Recycling Station"
	DigitalHubAssetType       = "Hailo Digital Hub"
)

// CreateAssetsIfNecessary create all assets for specification including sub specification if not already exists
func CreateAssetsIfNecessary(config conf.Config, spec hailo.Spec) error {

	for _, projectId := range config.ProjectIds {
		assetId, err := createAssetIfNecessary(config, projectId, nil, spec)
		if err != nil {
			log.Error("Hailo", "Could not create assets for device %s: %v", spec.DeviceId, err)
			return err
		}
		for _, subSpec := range spec.DeviceTypeSpecific.ComponentIdList {
			_, err = createAssetIfNecessary(config, projectId, assetId, subSpec)
			if err != nil {
				log.Error("Hailo", "Could not create assets for sub device %s: %v", subSpec.DeviceId, err)
				return err
			}
		}
	}

	return nil
}

// createAssetIfNecessary create asset for specification if not already exists
func createAssetIfNecessary(config conf.Config, projectId string, parentAssetId *int32, spec hailo.Spec) (*int32, error) {

	// Get known asset id from configuration
	existingId := conf.GetAssetId(config.Id, projectId, spec.DeviceId)
	if existingId != nil {
		return existingId, nil
	}

	log.Debug("hailo", "Creating new asset for project %s and spec %s.", projectId, spec.DeviceId)

	// If no asset id exists for project and configuration, create a new one
	name := name(spec)
	description := description(spec)
	newId, err := asset.UpsertAsset(api.Asset{
		ProjectId:               projectId,
		GlobalAssetIdentifier:   spec.Generic.DeviceSerial,
		Name:                    common.Ptr(name),
		AssetType:               assetType(spec),
		Description:             common.Ptr(description),
		ParentLocationalAssetId: parentAssetId,
	})
	if err != nil {
		return nil, err
	}
	if newId == nil {
		return nil, fmt.Errorf("cannot create asset: %s", name)
	}

	// Remember the asset id for further usage
	err = conf.InsertAsset(config.Id, projectId, spec.DeviceId, *newId)
	if err != nil {
		return newId, err
	}

	return newId, nil
}

// assetType from Hailo FDS specification
func assetType(specification hailo.Spec) string {
	if specification.DeviceTypeSpecific.ComponentIdList != nil {
		return RecyclingStationAssetType
	}
	return BinAssetType
}

func name(specification hailo.Spec) string {
	return fmt.Sprintf("%s (%s)", specification.DeviceId, specification.Generic.Model)
}

func description(specification hailo.Spec) string {
	if assetType(specification) == RecyclingStationAssetType || specification.DeviceTypeSpecific.Channel == "" || specification.DeviceTypeSpecific.ContentCategory == "" {
		return fmt.Sprintf("%s", specification.Generic.Model)
	} else {
		return fmt.Sprintf("%s (%s - %s)", specification.Generic.Model, specification.DeviceTypeSpecific.Channel, specification.DeviceTypeSpecific.ContentCategory)
	}
}
