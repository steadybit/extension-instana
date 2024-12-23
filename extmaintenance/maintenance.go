// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package extmaintenance

import (
	"context"
	"fmt"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/extension-instana/config"
	"github.com/steadybit/extension-instana/extapplications"
	"github.com/steadybit/extension-instana/types"
	extension_kit "github.com/steadybit/extension-kit"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
	"math"
	"net/http"
	"time"
)

type CreateMaintenanceWindowAction struct{}

// Make sure action implements all required interfaces
var (
	_ action_kit_sdk.Action[CreateMaintenanceWindowState]         = (*CreateMaintenanceWindowAction)(nil)
	_ action_kit_sdk.ActionWithStop[CreateMaintenanceWindowState] = (*CreateMaintenanceWindowAction)(nil)
)

type CreateMaintenanceWindowState struct {
	ApplicationPerspectiveId string
	DurationInMillis         int64
	ExperimentKey            *string
	ExecutionId              *int
	MaintenanceWindowId      *string
}

func NewCreateMaintenanceWindowAction() action_kit_sdk.Action[CreateMaintenanceWindowState] {
	return &CreateMaintenanceWindowAction{}
}
func (m *CreateMaintenanceWindowAction) NewEmptyState() CreateMaintenanceWindowState {
	return CreateMaintenanceWindowState{}
}

func (m *CreateMaintenanceWindowAction) Describe() action_kit_api.ActionDescription {
	return action_kit_api.ActionDescription{
		Id:          MaintenanceWindowActionId,
		Label:       "Create Maintenance Window",
		Description: "Start a Maintenance Window for a given duration.",
		Version:     extbuild.GetSemverVersionStringOrUnknown(),
		Icon:        extutil.Ptr(maintenanceWindowActionIcon),
		TargetSelection: extutil.Ptr(action_kit_api.TargetSelection{
			TargetType:          extapplications.ApplicationPerspectiveTargetId,
			QuantityRestriction: extutil.Ptr(action_kit_api.All),
			SelectionTemplates: extutil.Ptr([]action_kit_api.TargetSelectionTemplate{
				{
					Label: "application perspective label",
					Query: "instana.application.label=\"\"",
				},
			}),
		}),
		Technology:  extutil.Ptr("Instana"),
		Kind:        action_kit_api.Other,
		TimeControl: action_kit_api.TimeControlExternal,
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
		},
		Stop: extutil.Ptr(action_kit_api.MutatingEndpointReference{}),
	}
}

func (m *CreateMaintenanceWindowAction) Prepare(_ context.Context, state *CreateMaintenanceWindowState, request action_kit_api.PrepareActionRequestBody) (*action_kit_api.PrepareResult, error) {
	state.ApplicationPerspectiveId = request.Target.Attributes["instana.application.id"][0]
	state.ExperimentKey = request.ExecutionContext.ExperimentKey
	state.ExecutionId = request.ExecutionContext.ExecutionId
	state.DurationInMillis = int64(request.Config["duration"].(float64))
	return nil, nil
}

func (m *CreateMaintenanceWindowAction) Start(ctx context.Context, state *CreateMaintenanceWindowState) (*action_kit_api.StartResult, error) {
	return CreateMaintenanceWindow(ctx, state, &config.Config)
}

func (m *CreateMaintenanceWindowAction) Stop(ctx context.Context, state *CreateMaintenanceWindowState) (*action_kit_api.StopResult, error) {
	return DeleteMaintenanceWindow(ctx, state, &config.Config)
}

type MaintenanceWindowApi interface {
	CreateMaintenanceWindow(ctx context.Context, maintenanceWindow types.CreateMaintenanceWindowRequest) (*string, *http.Response, error)
	DeleteMaintenanceWindow(ctx context.Context, maintenanceWindowId string) (*http.Response, error)
}

func CreateMaintenanceWindow(ctx context.Context, state *CreateMaintenanceWindowState, api MaintenanceWindowApi) (*action_kit_api.StartResult, error) {
	name := "Steadybit"
	if state.ExperimentKey != nil && state.ExecutionId != nil {
		name = fmt.Sprintf("Steadybit %s - %d", *state.ExperimentKey, *state.ExecutionId)
	}

	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	if state.ExperimentKey != nil && state.ExecutionId != nil {
		id = fmt.Sprintf("%s-%d", *state.ExperimentKey, *state.ExecutionId)
	}

	amount := int64(math.Ceil(float64(state.DurationInMillis) / 1000 / 60))

	createRequest := types.CreateMaintenanceWindowRequest{
		Id:    id,
		Name:  name,
		Query: fmt.Sprintf("entity.application.id:\"%s\"", state.ApplicationPerspectiveId),
		Scheduling: types.Schedule{
			Duration: types.Duration{
				Amount: amount,
				Unit:   "MINUTES",
			},
			Start: time.Now().UnixMilli(),
			Type:  "ONE_TIME",
		},
	}

	windowId, _, err := api.CreateMaintenanceWindow(ctx, createRequest)
	if err != nil {
		return nil, extension_kit.ToError("Failed to create maintenance window.", err)
	}

	state.MaintenanceWindowId = windowId

	return &action_kit_api.StartResult{
		Messages: &action_kit_api.Messages{
			action_kit_api.Message{Level: extutil.Ptr(action_kit_api.Info), Message: fmt.Sprintf("Maintenance window created. (id %s)", *state.MaintenanceWindowId)},
		},
	}, nil
}

func DeleteMaintenanceWindow(ctx context.Context, state *CreateMaintenanceWindowState, api MaintenanceWindowApi) (*action_kit_api.StopResult, error) {
	if state.MaintenanceWindowId == nil {
		return nil, nil
	}

	resp, err := api.DeleteMaintenanceWindow(ctx, *state.MaintenanceWindowId)
	if err != nil {
		return nil, extension_kit.ToError(fmt.Sprintf("Failed to delete maintenace window (id %s). Full response: %v", *state.MaintenanceWindowId, resp), err)
	}

	return &action_kit_api.StopResult{
		Messages: &action_kit_api.Messages{
			action_kit_api.Message{Level: extutil.Ptr(action_kit_api.Info), Message: fmt.Sprintf("Maintenance window deleted. (id %s)", *state.MaintenanceWindowId)},
		},
	}, nil
}
