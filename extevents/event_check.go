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
	"github.com/steadybit/extension-instana/extapplications"
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
	End                   time.Time
	EventSeverityFilter   int
	EventTypeFilters      []string
	Condition             string
	ConditionCheckMode    string
	ConditionCheckSuccess bool
	SnapshotIds           map[string]bool
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
		Icon:        new(eventCheckActionIcon),
		TargetSelection: new(action_kit_api.TargetSelection{
			TargetType:          extapplications.ApplicationPerspectiveTargetId,
			QuantityRestriction: extutil.Ptr(action_kit_api.QuantityRestrictionAll),
			SelectionTemplates: new([]action_kit_api.TargetSelectionTemplate{
				{
					Label: "application perspective label",
					Query: "instana.application.label=\"\"",
				},
			}),
		}),
		Technology:  new("Instana"),
		Kind:        action_kit_api.Check,
		TimeControl: action_kit_api.TimeControlInternal,
		Parameters: []action_kit_api.ActionParameter{
			{
				Name:         "duration",
				Label:        "Duration",
				Description:  new(""),
				Type:         action_kit_api.ActionParameterTypeDuration,
				DefaultValue: new("30s"),
				Order:        new(1),
				Required:     new(true),
			},
			{
				Name:        "condition",
				Label:       "Condition",
				Description: new(""),
				Type:        action_kit_api.ActionParameterTypeString,
				Options: new([]action_kit_api.ParameterOption{
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
				DefaultValue: new(conditionShowOnly),
				Order:        new(2),
				Required:     new(true),
			},
			{
				Name:         "conditionCheckMode",
				Label:        "Condition Check Mode",
				Description:  new("Should the step succeed if the condition is met at least once or all the time?"),
				Type:         action_kit_api.ActionParameterTypeString,
				DefaultValue: new(conditionCheckModeAllTheTime),
				Options: new([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "All the time",
						Value: conditionCheckModeAllTheTime,
					},
					action_kit_api.ExplicitParameterOption{
						Label: "At least once",
						Value: conditionCheckModeAtLeastOnce,
					},
				}),
				Required: new(true),
				Order:    new(3),
			},
			{
				Name:        "eventSeverityFilter", //-1 = INFO, 5 = WARN, 10 = CRITICAL
				Label:       "Event Severity Filter",
				Description: new("Filter Problems by minimum severity."),
				Type:        action_kit_api.ActionParameterTypeString,
				Order:       new(4),
				Required:    new(true),
				Advanced:    new(true),
				Options: new([]action_kit_api.ParameterOption{
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
				DefaultValue: new(severityWarning),
			},
			{
				Name:        "eventTypeFilters",
				Label:       "Event Type Filter",
				Description: new("Filter Problems by an event type."),
				Type:        action_kit_api.ActionParameterTypeStringArray,
				Order:       new(5),
				Required:    new(true),
				Advanced:    new(true),
				Options: new([]action_kit_api.ParameterOption{
					action_kit_api.ExplicitParameterOption{
						Label: "Incident",
						Value: "INCIDENT",
					},
					action_kit_api.ExplicitParameterOption{
						Label: "Issue",
						Value: "ISSUE",
					},
				}),
				DefaultValue: new("[\"INCIDENT\",\"ISSUE\"]"),
			},
		},
		Widgets: new([]action_kit_api.Widget{
			action_kit_api.StateOverTimeWidget{
				Type:  action_kit_api.ComSteadybitWidgetStateOverTime,
				Title: "Instana Events",
				Identity: action_kit_api.StateOverTimeWidgetIdentityConfig{
					From: "id",
				},
				Label: action_kit_api.StateOverTimeWidgetLabelConfig{
					From: "title",
				},
				State: action_kit_api.StateOverTimeWidgetStateConfig{
					From: "state",
				},
				Tooltip: action_kit_api.StateOverTimeWidgetTooltipConfig{
					From: "tooltip",
				},
				Url: new(action_kit_api.StateOverTimeWidgetUrlConfig{
					From: new("url"),
				}),
				Value: new(action_kit_api.StateOverTimeWidgetValueConfig{
					Hide: new(true),
				}),
			},
		}),
		Prepare: action_kit_api.MutatingEndpointReference{},
		Start:   action_kit_api.MutatingEndpointReference{},
		Status: new(action_kit_api.MutatingEndpointReferenceWithCallInterval{
			CallInterval: new("5s"),
		}),
	}
}

func (m *EventCheckAction) Prepare(ctx context.Context, state *EventCheckState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	duration := request.Config["duration"].(float64)
	state.Start = time.Now()
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

	applicationPerspectiveId := request.Target.Attributes["instana.application.id"][0]
	snapshotIds, err := config.Config.GetSnapshotIds(ctx, applicationPerspectiveId)
	if err != nil {
		return nil, extension_kit.ToError("Failed to get snapshot-ids from Instana.", err)
	}
	state.SnapshotIds = make(map[string]bool)
	for _, snapshotId := range snapshotIds {
		state.SnapshotIds[snapshotId] = true
	}
	log.Debug().Int("count", len(state.SnapshotIds)).Msg("Initialized snapshot ids.")

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
	GetSnapshotIds(ctx context.Context, applicationPerspectiveId string) ([]string, error)
}

func EventCheckStatus(ctx context.Context, state *EventCheckState, api EventsApi) (*action_kit_api.StatusResult, error) {
	now := time.Now()
	events, err := api.GetEvents(ctx, state.Start, now, state.EventTypeFilters)
	if err != nil {
		return nil, extension_kit.ToError("Failed to get events from Instana.", err)
	}

	filteredEvents := make([]types.Event, 0)
	for _, event := range events {
		_, snapshotIdMatchesFilter := state.SnapshotIds[event.SnapshotId]
		if event.Severity >= state.EventSeverityFilter && snapshotIdMatchesFilter && event.State != "closed" {
			filteredEvents = append(filteredEvents, event)
		}
	}

	completed := now.After(state.End)
	var checkError *action_kit_api.ActionKitError
	if state.ConditionCheckMode == conditionCheckModeAllTheTime {
		if state.Condition == conditionNoEvents && len(filteredEvents) > 0 {
			checkError = new(action_kit_api.ActionKitError{
				Title:  fmt.Sprintf("No event expected, but %d events found.", len(filteredEvents)),
				Status: extutil.Ptr(action_kit_api.Failed),
			})
		}
		if state.Condition == conditionAtLeastOneEvent && len(filteredEvents) == 0 {
			checkError = new(action_kit_api.ActionKitError{
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
				checkError = new(action_kit_api.ActionKitError{
					Title:  "No event expected, but events found.",
					Status: extutil.Ptr(action_kit_api.Failed),
				})
			} else if state.Condition == conditionAtLeastOneEvent {
				checkError = new(action_kit_api.ActionKitError{
					Title:  "At least one event expected, but no events found.",
					Status: extutil.Ptr(action_kit_api.Failed),
				})
			}
		}
	}

	return &action_kit_api.StatusResult{
		Completed: completed,
		Error:     checkError,
		Metrics:   eventsToMetrics(filteredEvents, now),
	}, nil
}
func eventsToMetrics(events []types.Event, now time.Time) *action_kit_api.Metrics {
	var metrics []action_kit_api.Metric
	for _, event := range events {
		tooltip := fmt.Sprintf("Event Problem: %s\nEvent Detail: %s\nEvent Type: %s\nEvent Severity: %d\nEntity Name: %s\nEntity Label: %s\nEntity Type: %s", event.Problem, event.Detail, event.Type, event.Severity, event.EntityName, event.EntityLabel, event.EntityType)
		metrics = append(metrics,
			action_kit_api.Metric{
				Name: new("instana_events"),
				Metric: map[string]string{
					"id":      event.EventId,
					"title":   event.Problem + " - " + event.Detail,
					"state":   getState(event.Severity),
					"tooltip": tooltip,
					"url":     fmt.Sprintf("%s/#/events;eventId=%s", config.Config.BaseUrl, event.EventId),
				},
				Timestamp: now,
				Value:     0,
			},
		)
	}
	return new(metrics)
}

func getState(severity int) string {
	if severity == -1 {
		return "info"
	} else if severity == 5 {
		return "warn"
	} else if severity == 10 {
		return "danger"
	}
	return "danger"
}
