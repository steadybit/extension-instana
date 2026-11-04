// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package extevents

import "github.com/steadybit/action-kit/go/action_kit_api/v2"

const groupEventFilter = "Event Filter"

var eventFilterGroup = action_kit_api.ParameterGroup{
	Name:        groupEventFilter,
	Description: "Only matching events count for the condition and are shown.",
}
