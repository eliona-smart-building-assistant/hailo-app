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

package main

import (
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"hailo/conf"
	"hailo/eliona"
	"hailo/hailo"
	"time"
)

// collectData collects data based on the configured FDS endpoints in table hailo.config. For each FDS endpoint the
// data is collected in a separate thread. These threads wait until the configured interval time is over. Afterwards
// a new thread is started for this connection.
func collectData() {

	// Check if import or update of assets is requested

	// Load all configured configs from table hailo.config.
	configs := conf.GetConfigs()
	if len(configs) <= 0 {
		log.Fatal("Hailo", "Couldn't read config from configured database.")
	}

	// Start collection data for each config
	for _, config := range configs {

		// Skip config if disabled and set inactive
		if !config.Enable {
			if config.Active {
				conf.SetConfigActive(config.Id, false)
			}
			continue
		}

		// Signals, that this config is active
		if !config.Active {
			conf.SetConfigActive(config.Id, true)
			log.Info("Hailo", "Collecting %d initialized with config:\n"+
				"FDS Fds Endpoint: %s\n"+
				"FDS Fds Auth Server: %s\n"+
				"Auth Timeout: %d\n"+
				"Request Timeout: %d",
				config.Id,
				config.FdsConfig.FdsServer,
				config.FdsConfig.AuthServer,
				config.AuthTimeout,
				config.RequestTimeout)
		}

		// Runs the ReadNode. If the current node is currently running, skip the execution
		// After the execution sleeps the configured timeout. During this timeout no further
		// process for this config (appId) is started to read the data.
		common.RunOnce(func() {

			log.Info("Hailo", "Collecting %d started", config.Id)

			// Collect data for the config
			collectDataForConfig(config)

			log.Info("Hailo", "Collecting %d finished", config.Id)

			// Waits until the time is excited
			time.Sleep(time.Second * time.Duration(config.IntervalSec))

		}, config.Id)
	}
}

// collectDataForConfig reads specification of all devices in the given connection. For all devices found asset
// data is written. In case of stations (group multiple component devices) data for each component is read and
// written.
func collectDataForConfig(config conf.Config) {

	// Read specs from Hailo FDS
	specs, err := hailo.GetSpecs(config)
	if err != nil {
		log.Error("Hailo", "Could not read specs for config %d: %v", config.Id, err)
		return
	}

	// For each spec write asset data
	for _, spec := range specs.Data {

		// If necessary create assets in eliona
		err = eliona.CreateAssetsIfNecessary(config, spec)
		if err != nil {
			return
		}

		// Writing asset data for specification
		err = eliona.UpsertDataForDevices(config, spec)
		if err != nil {
			return
		}

		// Get Status
		status, err := hailo.GetStatus(config, spec.DeviceId)
		if err != nil {
			log.Error("Hailo", "Could not read status for config %d and device '%s': %v", config.Id, spec.DeviceId, err)
			return
		}

		// Decide if device is station or single container
		if status.IsStation() {

			// Upsert status for station
			err = eliona.UpsertDataForStation(config, status)
			if err != nil {
				return
			}

			// Process station components
			for _, compStatus := range status.DeviceTypeSpecific.CompStatuses {

				// Get diag for component
				diag, err := hailo.GetDiag(config, compStatus.DeviceId)
				if err != nil {
					log.Error("Hailo", "Could not read diag for config %d and component '%s': %v", config.Id, compStatus.DeviceId, err)
					return
				}

				// Upsert status and diag for station components
				err = eliona.UpsertDataForBin(config, compStatus, diag)
				if err != nil {
					return
				}

			}

		} else {

			// Get diag for single container
			diag, err := hailo.GetDiag(config, status.DeviceId)
			if err != nil {
				log.Error("Hailo", "Could not read diag for config %d and station '%s': %v", config.Id, status.DeviceId, err)
				return
			}

			// Upsert status and diag for station single container
			err = eliona.UpsertDataForBin(config, status, diag)
			if err != nil {
				return
			}
		}
	}

}
