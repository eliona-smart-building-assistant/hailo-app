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
	"github.com/eliona-smart-building-assistant/go-eliona/apps"
	"github.com/eliona-smart-building-assistant/go-eliona/common"
	"github.com/eliona-smart-building-assistant/go-eliona/db"
	"github.com/eliona-smart-building-assistant/go-eliona/log"
	"hailo/conf"
	"hailo/eliona"
	"os"
	"time"
)

// The main function starts the app by starting all services necessary for this app and waits
// until all services are finished. In most cases the services run infinite, except the app is stopped
// externally, e.g. during a shut-down of the eliona environment.
func main() {
	log.Info("Hailo", "Starting the app.")

	// Check program args. If test, then print output and exit
	args := determineArgs()
	if args.test {
		printData(args.authServer, args.userName, args.password, args.fdsEndpoint)
		os.Exit(0)
	}

	// Necessary to close used init resources, because db.Pool() is used in this app.
	defer db.ClosePool()

	// Init the app before the first run.
	apps.Init(db.Pool(), common.AppName(),
		apps.ExecSqlFile("conf/init.sql"),
		eliona.InitAssetTypes,
		conf.InitConfiguration,
	)

	// Patch the app
	apps.Patch(db.Pool(), common.AppName(), "020000",
		apps.ExecSqlFile("conf/v2.0.0.sql"))

	// Starting the service to collect the data for each configured Hailo Smart Hub.
	apps.WaitFor(
		apps.Loop(collectData, time.Second*11),
	)

	// At the end set all configuration inactive
	conf.SetAllConfigsInactive()

	log.Info("Hailo", "Terminate the app.")
}
