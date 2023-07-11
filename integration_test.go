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
	"github.com/eliona-smart-building-assistant/app-integration-tests/app"
	"github.com/eliona-smart-building-assistant/app-integration-tests/assert"
	"github.com/eliona-smart-building-assistant/app-integration-tests/test"
	"testing"
)

func TestApp(t *testing.T) {
	app.StartApp()
	test.AppWorks(t)
	t.Run("TestAssetTypes", assetTypes)
	t.Run("TestWidgetTypes", widgetTypes)
	t.Run("TestSchema", schema)
	app.StopApp()
}

func assetTypes(t *testing.T) {
	t.Parallel()

	assert.AssetTypeExists(t, "Hailo FDS Bin", []string{})
	assert.AssetTypeExists(t, "Hailo Digital Hub", []string{})
	assert.AssetTypeExists(t, "Hailo FDS Recycling Station", []string{})
}

func widgetTypes(t *testing.T) {
	t.Parallel()

	assert.WidgetTypeExists(t, "Hailo")
	assert.WidgetTypeExists(t, "Hailo Station")
}

func schema(t *testing.T) {
	t.Parallel()

	assert.SchemaExists(t, "hailo", []string{"config", "asset"})
}
