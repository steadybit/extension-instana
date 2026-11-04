// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extevents

import (
	"testing"

	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/stretchr/testify/assert"
)

func TestParameterGroups(t *testing.T) {
	desc := NewEventCheckAction().Describe()

	var grouped []string
	lastOrderOutside, firstOrderInside := -1, 1<<30
	for _, p := range desc.Parameters {
		if p.Group == nil {
			// No advanced parameter is left alone in "Additional Settings" on platforms with groups.
			assert.False(t, p.Advanced != nil && *p.Advanced, "parameter %s", p.Name)
			lastOrderOutside = max(lastOrderOutside, *p.Order)
			continue
		}
		assert.Equal(t, groupEventFilter, *p.Group)
		// Platforms without groups keep showing the filters in "Additional Settings".
		assert.True(t, p.Advanced != nil && *p.Advanced, "parameter %s", p.Name)
		grouped = append(grouped, p.Name)
		firstOrderInside = min(firstOrderInside, *p.Order)
	}
	assert.Equal(t, []string{"eventSeverityFilter", "eventTypeFilters"}, grouped)
	assert.Less(t, lastOrderOutside, firstOrderInside, "the filters come after the condition settings")
	assert.Equal(t, []action_kit_api.ParameterGroup{eventFilterGroup}, *desc.ParameterGroups)
}
