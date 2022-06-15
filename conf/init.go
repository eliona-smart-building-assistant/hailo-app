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

// InitConfiguration creates a default configuration to demonstrate how the eliona app should be configured. This configuration
// points to a not existing endpoint and have to be changed.
func InitConfiguration(connection db.Connection) error {
	err := db.Exec(connection, "insert into hailo.config ("+
		"config, "+
		"enable, "+
		"description, "+
		"interval_sec, "+
		"inactive_timeout)"+
		" values ($1, $2, $3, $4, $5) on conflict do nothing",
		FdsConfig{
			Name:       "username",
			Password:   "password",
			AuthServer: "https://foo.execute-api.eu-central-1.amazonaws.com",
			FdsServer:  "https://bar.execute-api.eu-central-1.amazonaws.com/hailo/v1",
		},
		true,
		"Hailo FDS demo configuration. Please change to your hailo server endpoints and authentication.",
		60,       // 1 minute
		12*60*60, // 12 hours
	)
	if err != nil {
		log.Error("Hailo", "Error during init hailo config with a demo configuration: %v", err)
	}
	return err
}
