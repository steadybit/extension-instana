// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 Steadybit GmbH

package e2e

import (
	"fmt"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_test/e2e"
	"github.com/steadybit/extension-instana/extevents"
	"github.com/steadybit/extension-kit/extlogging"
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
			Name: "event check",
			Test: testEventCheck,
		},
	})
}

func testEventCheck(t *testing.T, m *e2e.Minikube, e *e2e.Extension) {
	defer func() { Requests = []string{} }()

	config := struct {
		Duration            int      `json:"duration"`
		EventSeverityFilter string   `json:"eventSeverityFilter"`
		EventTypeFilters    []string `json:"eventTypeFilters"`
		Condition           string   `json:"condition"`
		ConditionCheckMode  string   `json:"conditionCheckMode"`
	}{Duration: 1000, EventSeverityFilter: "info", EventTypeFilters: []string{"ISSUE", "INCIDENT", "CHANGE"}, Condition: "showOnly", ConditionCheckMode: "allTheTime"}

	executionContext := &action_kit_api.ExecutionContext{}

	action, err := e.RunAction(extevents.EventCheckActionId, nil, config, executionContext)
	defer func() { _ = action.Cancel() }()
	require.NoError(t, err)
	err = action.Wait()
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		messages := action.Messages()
		if messages == nil {
			return false
		}
		return len(messages) > 0
	}, 5*time.Second, 500*time.Millisecond)
	messages := action.Messages()

	for _, message := range messages {
		assert.Equal(t, "INSTANA", *message.Type)
		assert.Equal(t, action_kit_api.Info, *message.Level)
		assert.Equal(t, "offline - JVM on Host ip-10-10-81-117.eu-central-1.compute.internal", message.Message)
	}
}
