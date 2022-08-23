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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeConvert(t *testing.T) {
	tm := time.Now().Unix()
	iso := parseToIsoTime(tm)
	unix := parseTime(iso).Unix()
	assert.Equal(t, tm, unix)
}

func TestTimeSinceTimeUntil(t *testing.T) {

	unix := time.Now().Unix()
	unix -= 2 * 60 * 60
	iso := parseToIsoTime(unix)
	result := parseTimeToHours(iso)
	assert.Equal(t, 2.0, result)

	unix += 10 * 60 * 60
	iso = parseToIsoTime(unix)
	result = parseTimeToHours(iso)
	assert.Equal(t, 8.0, result)

	unix = time.Now().Unix()
	unix += 3 * 24 * 60 * 60
	iso = parseToIsoTime(unix)
	result = parseTimeToDays(iso)
	assert.Equal(t, 3.0, result)

	unix -= 44 * 24 * 60 * 60
	iso = parseToIsoTime(unix)
	result = parseTimeToDays(iso)
	assert.Equal(t, 41.0, result)

}
