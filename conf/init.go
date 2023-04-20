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
	"context"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"hailo/apiserver"
)

// InitConfiguration creates a default configuration to demonstrate how the eliona app should be configured. This configuration
// points to a not existing endpoint and have to be changed.
func InitConfiguration(connection db.Connection) error {
	_, err := InsertConfig(context.Background(), apiserver.Configuration{
		Enable:          common.Ptr(false),
		Description:     common.Ptr("Hailo FDS demo configuration. Please change to your hailo server endpoints and authentication."),
		IntervalSec:     60,           // 1 minute
		InactiveTimeout: 12 * 60 * 60, // 12 hours
		Username:        common.Ptr("username"),
		Password:        common.Ptr("password"),
		AuthServer:      common.Ptr("https://foo.execute-api.eu-central-1.amazonaws.com"),
		FdsServer:       common.Ptr("https://bar.execute-api.eu-central-1.amazonaws.com/hailo/v1"),
	})
	return err
}
