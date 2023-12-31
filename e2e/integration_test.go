// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 Steadybit GmbH

package e2e

import (
	"fmt"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_test/e2e"
	"github.com/steadybit/discovery-kit/go/discovery_kit_test/validate"
	"github.com/steadybit/extension-instana/extevents"
	"github.com/steadybit/extension-instana/extmaintenance"
	"github.com/steadybit/extension-kit/extlogging"
	"github.com/steadybit/extension-kit/extutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestWithMinikube(t *testing.T) {
	extlogging.InitZeroLog()
	server := createMockInstanaServer()
	defer server.Close()
	split := strings.SplitAfter(server.URL, ":")
	port := split[len(split)-1]

	extFactory := e2e.HelmExtensionFactory{
		Name: "extension-instana",
		Port: 8090,
		ExtraArgs: func(m *e2e.Minikube) []string {
			return []string{
				"--set", "logging.level=debug",
				"--set", "instana.apiKey=api-key-123",
				"--set", fmt.Sprintf("instana.baseUrl=http://host.minikube.internal:%s", port),
			}
		},
	}

	e2e.WithDefaultMinikube(t, &extFactory, []e2e.WithMinikubeTestCase{
		{
			Name: "validate discovery",
			Test: validateDiscovery,
		},
		{
			Name: "event check",
			Test: testEventCheck,
		},
		{
			Name: "create maintenance window",
			Test: testCreateMaintenanceWindow,
		},
	})
}

func validateDiscovery(t *testing.T, _ *e2e.Minikube, e *e2e.Extension) {
	assert.NoError(t, validate.ValidateEndpointReferences("/", e.Client))
}

func testEventCheck(t *testing.T, m *e2e.Minikube, e *e2e.Extension) {
	defer func() { Requests = []string{} }()

	target := &action_kit_api.Target{
		Name: "Application Perspective 1",
		Attributes: map[string][]string{
			"instana.application.id":   {"application-id-1"},
			"instana.application.name": {"application-name-1"},
		},
	}

	config := struct {
		Duration            int      `json:"duration"`
		EventSeverityFilter string   `json:"eventSeverityFilter"`
		EventTypeFilters    []string `json:"eventTypeFilters"`
		Condition           string   `json:"condition"`
		ConditionCheckMode  string   `json:"conditionCheckMode"`
	}{Duration: 1000, EventSeverityFilter: "info", EventTypeFilters: []string{"ISSUE", "INCIDENT", "CHANGE"}, Condition: "showOnly", ConditionCheckMode: "allTheTime"}

	executionContext := &action_kit_api.ExecutionContext{}

	action, err := e.RunAction(extevents.EventCheckActionId, target, config, executionContext)
	defer func() { _ = action.Cancel() }()
	require.NoError(t, err)
	err = action.Wait()
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		metrics := action.Metrics()
		if metrics == nil {
			return false
		}
		return len(metrics) > 0
	}, 5*time.Second, 500*time.Millisecond)
	metrics := action.Metrics()

	for _, metric := range metrics {
		assert.Equal(t, "XoDkHMssTRGiyEpPz337hQ", metric.Metric["id"])
		assert.Equal(t, "Condition [Ready]: Pod containers are not ready - containers with unready status: [platform platform-port-splitter]", metric.Metric["title"])
	}
}

func testCreateMaintenanceWindow(t *testing.T, m *e2e.Minikube, e *e2e.Extension) {
	defer func() { Requests = []string{} }()

	target := &action_kit_api.Target{
		Name: "Application Perspective 1",
		Attributes: map[string][]string{
			"instana.application.id":   {"application-id-1"},
			"instana.application.name": {"application-name-1"},
		},
	}

	config := struct {
		Duration int `json:"duration"`
	}{Duration: 10000}

	executionContext := &action_kit_api.ExecutionContext{
		ExperimentKey: extutil.Ptr("TST-1"),
		ExecutionId:   extutil.Ptr(47),
	}

	action, err := e.RunAction(extmaintenance.MaintenanceWindowActionId, target, config, executionContext)
	defer func() { _ = action.Cancel() }()
	require.NoError(t, err)
	err = action.Wait()
	require.NoError(t, err)
	require.Contains(t, Requests, "PUT-/api/settings/v2/maintenance/TST-1-47")
	require.Contains(t, Requests, "DELETE-/api/settings/v2/maintenance/TST-1-47")
}
