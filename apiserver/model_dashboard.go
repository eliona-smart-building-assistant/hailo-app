/*
 * Hailo app API
 *
 * API to access and configure the Hailo app
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// Dashboard - A frontend dashboard
type Dashboard struct {

	// The internal Id of dashboard
	Id *int32 `json:"id,omitempty"`

	// The name for this dashboard
	Name string `json:"name"`

	// ID of the project to which the dashboard belongs
	ProjectId string `json:"projectId"`

	// ID of the user who owns the dashboard
	UserId string `json:"userId"`

	// The sequence of the dashboard
	Sequence *int32 `json:"sequence,omitempty"`

	// List of widgets on this dashboard (order matches the order of widgets on the dashboard)
	Widgets *[]Widget `json:"widgets,omitempty"`
}

// AssertDashboardRequired checks if the required fields are not zero-ed
func AssertDashboardRequired(obj Dashboard) error {
	elements := map[string]interface{}{
		"name":      obj.Name,
		"projectId": obj.ProjectId,
		"userId":    obj.UserId,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if obj.Widgets != nil {
		for _, el := range *obj.Widgets {
			if err := AssertWidgetRequired(el); err != nil {
				return err
			}
		}
	}
	return nil
}

// AssertRecurseDashboardRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Dashboard (e.g. [][]Dashboard), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDashboardRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDashboard, ok := obj.(Dashboard)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDashboardRequired(aDashboard)
	})
}
