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
	"github.com/eliona-smart-building-assistant/go-eliona/api"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
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
		err := createAssetIfNecessary(config, projectId, spec)
		if err != nil {
			log.Error("Hailo", "Could not create assets for device %s: %v", spec.DeviceId, err)
			return err
		}
		for _, subSpec := range spec.DeviceTypeSpecific.ComponentIdList {
			err = createAssetIfNecessary(config, projectId, subSpec)
			if err != nil {
				log.Error("Hailo", "Could not create assets for sub device %s: %v", subSpec.DeviceId, err)
				return err
			}
		}
	}

	return nil
}

// createAssetIfNecessary create asset for specification if not already exists
func createAssetIfNecessary(config conf.Config, projectId string, specification hailo.Spec) error {

	// Get known asset id from configuration
	existingId := conf.GetAssetId(config.Id, projectId, specification.DeviceId)
	if existingId != nil {
		return nil
	}

	log.Debug("hailo", "Creating new asset for project %s and specification %s.", projectId, specification.DeviceId)

	// If no asset id exists for project and configuration, create a new one
	name := fmt.Sprintf("%s (%s)", specification.DeviceId, specification.Generic.Model)
	newId, err := asset.UpsertAsset(api.Asset{
		ProjectId:             projectId,
		GlobalAssetIdentifier: specification.Generic.DeviceSerial,
		Name:                  common.Ptr(name),
		AssetType:             assetType(specification),
		Description:           common.Ptr(specification.DeviceTypeSpecific.Channel + " - " + specification.DeviceTypeSpecific.ContentCategory),
	})
	if err != nil {
		return err
	}
	if newId == nil {
		return fmt.Errorf("cannot create asset: %s", name)
	}

	// Remember the asset id for further usage
	err = conf.InsertAsset(config.Id, projectId, specification.DeviceId, *newId)
	if err != nil {
		return err
	}

	return nil
}

// assetType from Hailo FDS specification
func assetType(specification hailo.Spec) string {
	if specification.DeviceTypeSpecific.ComponentIdList != nil {
		return RecyclingStationAssetType
	}
	return BinAssetType
}
