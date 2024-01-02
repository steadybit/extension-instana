// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package extevents

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/extension-instana/config"
	"github.com/steadybit/extension-instana/types"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
	"time"
)

type EventCheckAction struct{}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[EventCheckState]           = (*EventCheckAction)(nil)
	_ action_kit_sdk.ActionWithStatus[EventCheckState] = (*EventCheckAction)(nil)
)

type EventCheckState struct {
	Start                 time.Time
	LastFetch             time.Time
	End                   time.Time
	EventSeverityFilter   int
	EventTypeFilters      []string
	Condition             string
	ConditionCheckMode    string
	ConditionCheckSuccess bool
}

func NewEventCheckAction() action_kit_sdk.Action[EventCheckState] {
	return &EventCheckAction{}
}

func (m *EventCheckAction) NewEmptyState() EventCheckState {
	return EventCheckState{}
}

func (m *EventCheckAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          EventCheckActionId,
		Label:       "Event Check",
		Description: "Checks for the existence of certain events in Instana.",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		Icon:        extutil.Ptr(eventCheckActionIcon),
		Category:    extutil.Ptr("monitoring"),
		Kind:        action_kit_api.Check,
		TimeControl: action_kit_api.TimeControlInternal,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "duration",
				Label:        "Duration",
				Description:  extutil.Ptr(""),
				Type:         action_kit_api.Duration,
				DefaultValue: extutil.Ptr("30s"),
				Order:        extutil.Ptr(1),
				Required:     extutil.Ptr(true),
			},
			{
				Name:        "condition",
				Label:       "Condition",
				Description: extutil.Ptr(""),
				Type:        action_kit_api.String,
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "No check, only show events",
						Value: conditionShowOnly,
					},
					action_kit_api.ExplicitParameterOption{
						Label: "No event expected",
						Value: conditionNoEvents,
					},
					action_kit_api.ExplicitParameterOption{
						Label: "At least one event expected",
						Value: conditionAtLeastOneEvent,
					},
				}),
				DefaultValue: extutil.Ptr(conditionShowOnly),
				Order:        extutil.Ptr(2),
				Required:     extutil.Ptr(true),
			},
			{
				Name:         "conditionCheckMode",
				Label:        "Condition Check Mode",
				Description:  extutil.Ptr("Should the step succeed if the condition is met at least once or all the time?"),
				Type:         action_kit_api.String,
				DefaultValue: extutil.Ptr(conditionCheckModeAllTheTime),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "All the time",
						Value: conditionCheckModeAllTheTime,
					},
					action_kit_api.ExplicitParameterOption{
						Label: "At least once",
						Value: conditionCheckModeAtLeastOnce,
					},
				}),
				Required: extutil.Ptr(true),
				Order:    extutil.Ptr(3),
			},
			{
				Name:        "eventSeverityFilter", //-1 = INFO, 5 = WARN, 10 = CRITICAL
				Label:       "Event Severity Filter",
				Description: extutil.Ptr("Filter Problems by minimum severity."),
				Type:        action_kit_api.String,
				Order:       extutil.Ptr(4),
				Required:    extutil.Ptr(true),
				Advanced:    extutil.Ptr(true),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "Info",
						Value: severityInfo,
					},
					action_kit_api.ExplicitParameterOption{
						Label: "Warning",
						Value: severityWarning,
					},
					action_kit_api.ExplicitParameterOption{
						Label: "Critical",
						Value: severityCritical,
					},
				}),
				DefaultValue: extutil.Ptr(severityWarning),
			},
			{
				Name:        "eventTypeFilters",
				Label:       "Event Type Filter",
				Description: extutil.Ptr("Filter Problems by an event type."),
				Type:        action_kit_api.StringArray,
				Order:       extutil.Ptr(5),
				Required:    extutil.Ptr(true),
				Advanced:    extutil.Ptr(true),
				Options: extutil.Ptr([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "Incident",
						Value: "INCIDENT",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "Issue",
						Value: "ISSUE",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "Change",
						Value: "CHANGE",
					},
				}),
				DefaultValue: extutil.Ptr("[\"INCIDENT\",\"ISSUE\"]"),
			},
		},
		Widgets: extutil.Ptr([]action_kit_api.Widget{
			action_kit_api.LogWidget{
				Type:    action_kit_api.ComSteadybitWidgetLog,
				Title:   "Instana Events",
				LogType: logType,
			},
		}),
		Prepare: action_kit_api.MutatingEndpointReference{},
		Start:   action_kit_api.MutatingEndpointReference{},
		Status: extutil.Ptr(action_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("5s"),
		}),
	}
}

func (m *EventCheckAction) Prepare(_ context.Context, state *EventCheckState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	duration := request.Config["duration"].(float64)
	state.Start = time.Now()
	state.LastFetch = state.Start
	state.End = time.Now().Add(time.Millisecond * time.Duration(duration))

	if request.Config["eventSeverityFilter"] != nil {
		severityFilter := fmt.Sprintf("%v", request.Config["eventSeverityFilter"])
		if severityFilter == severityInfo {
			state.EventSeverityFilter = -1
		} else if severityFilter == severityWarning {
			state.EventSeverityFilter = 5
		} else if severityFilter == severityCritical {
			state.EventSeverityFilter = 10
		} else {
			return nil, extension_kit.ToError(fmt.Sprintf("Unknown Event Severity Filter: '%s'.", severityFilter), nil)
		}
	} else {
		return nil, extension_kit.ToError("Event Severity Filter is required.", nil)
	}

	state.EventTypeFilters = extutil.ToStringArray(request.Config["eventTypeFilters"])

	if request.Config["condition"] != nil {
		state.Condition = fmt.Sprintf("%v", request.Config["condition"])
	}
	if request.Config["conditionCheckMode"] != nil {
		state.ConditionCheckMode = fmt.Sprintf("%v", request.Config["conditionCheckMode"])
	}

	return nil, nil
}

func (m *EventCheckAction) Start(ctx context.Context, state *EventCheckState) (*action_kit_api.StartResult, error) {
	statusResult, err := EventCheckStatus(ctx, state, &config.Config)
	if statusResult == nil {
		return nil, err
	}
	startResult := action_kit_api.StartResult{
		Artifacts: statusResult.Artifacts,
		Error:     statusResult.Error,
		Messages:  statusResult.Messages,
		Metrics:   statusResult.Metrics,
	}
	return &startResult, err
}

func (m *EventCheckAction) Status(ctx context.Context, state *EventCheckState) (*action_kit_api.StatusResult, error) {
	return EventCheckStatus(ctx, state, &config.Config)
}

type EventsApi interface {
	GetEvents(ctx context.Context, from time.Time, to time.Time, eventTypeFilters []string) ([]types.Event, error)
}

func EventCheckStatus(ctx context.Context, state *EventCheckState, api EventsApi) (*action_kit_api.StatusResult, error) {
	now := time.Now()
	events, err := api.GetEvents(ctx, state.LastFetch, now, state.EventTypeFilters)
	state.LastFetch = now
	if err != nil {
		return nil, extension_kit.ToError("Failed to get events from Instana.", err)
	}

	filteredEvents := make([]types.Event, 0)
	for _, event := range events {
		if event.Severity >= state.EventSeverityFilter {
			filteredEvents = append(filteredEvents, event)
		}
	}

	completed := now.After(state.End)
	var checkError *action_kit_api.ActionKitError
	if state.ConditionCheckMode == conditionCheckModeAllTheTime {
		if state.Condition == conditionNoEvents && len(filteredEvents) > 0 {
			checkError = extutil.Ptr(action_kit_api.ActionKitError{
				Title:  fmt.Sprintf("No event expected, but %d events found.", len(filteredEvents)),
				Status: extutil.Ptr(action_kit_api.Failed),
			})
		}
		if state.Condition == conditionAtLeastOneEvent && len(filteredEvents) == 0 {
			checkError = extutil.Ptr(action_kit_api.ActionKitError{
				Title:  "At least one event expected, but no events found.",
				Status: extutil.Ptr(action_kit_api.Failed),
			})
		}

	} else if state.ConditionCheckMode == conditionCheckModeAtLeastOnce {
		if state.Condition == conditionNoEvents && len(filteredEvents) == 0 {
			state.ConditionCheckSuccess = true
		}
		if state.Condition == conditionAtLeastOneEvent && len(filteredEvents) > 0 {
			state.ConditionCheckSuccess = true
		}
		if completed && !state.ConditionCheckSuccess {
			if state.Condition == conditionNoEvents {
				checkError = extutil.Ptr(action_kit_api.ActionKitError{
					Title:  "No event expected, but events found.",
					Status: extutil.Ptr(action_kit_api.Failed),
				})
			} else if state.Condition == conditionAtLeastOneEvent {
				checkError = extutil.Ptr(action_kit_api.ActionKitError{
					Title:  "At least one event expected, but no events found.",
					Status: extutil.Ptr(action_kit_api.Failed),
				})
			}
		}
	}

	return &action_kit_api.StatusResult{
		Completed: completed,
		Error:     checkError,
		Messages:  eventsToMessages(filteredEvents),
	}, nil
}

func eventsToMessages(events []types.Event) *action_kit_api.Messages {
	var messages []action_kit_api.Message
	for _, event := range events {
		messages = append(messages, action_kit_api.Message{
			Message:         event.Problem + " - " + event.Detail,
			Type:            extutil.Ptr(logType),
			Level:           convertToLevel(event.Severity),
			Timestamp:       extutil.Ptr(time.UnixMilli(event.Start)),
			TimestampSource: extutil.Ptr(action_kit_api.TimestampSourceExternal),
			Fields: extutil.Ptr(action_kit_api.MessageFields{
				"Type":         event.Type,
				"Entity name":  event.EntityName,
				"Entity label": event.EntityLabel,
				"Entity type":  event.EntityType,
			}),
		})
	}
	log.Debug().Int("count", len(messages)).Msg("Found events, send to platform.")
	return extutil.Ptr(messages)
}

func convertToLevel(severity int) *action_kit_api.MessageLevel {
	if severity == -1 {
		return extutil.Ptr(action_kit_api.Info)
	} else if severity == 5 {
		return extutil.Ptr(action_kit_api.Warn)
	} else if severity == 10 {
		return extutil.Ptr(action_kit_api.Error)
	}
	return nil
}
