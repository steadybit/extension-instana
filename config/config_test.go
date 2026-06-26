// Copyright 2025 steadybit GmbH. All rights reserved.

package config

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/steadybit/extension-instana/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSnapshotIds_EscapesApplicationPerspectiveId(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer srv.Close()

	spec := Specification{BaseUrl: srv.URL, ApiToken: "X"}
	// An id that tries to inject an extra query parameter (override the size limit).
	_, err := spec.GetSnapshotIds(context.Background(), "app-1&size=1")
	require.NoError(t, err)

	// Escaped properly, the injected text stays inside the query value...
	assert.Equal(t, "entity.application.id:app-1&size=1", gotQuery.Get("query"))
	// ...and does not override the size parameter.
	assert.Equal(t, "20000", gotQuery.Get("size"))
}

func TestGetEvents_EscapesEventTypeFilter(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	spec := Specification{BaseUrl: srv.URL, ApiToken: "X"}
	_, err := spec.GetEvents(context.Background(), time.Time{}, time.Time{}, []string{"incident&to=0"})
	require.NoError(t, err)

	assert.Equal(t, []string{"incident&to=0"}, gotQuery["eventTypeFilters"])
}

func TestCreateMaintenanceWindow_EscapesIdInPath(t *testing.T) {
	var gotEscapedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotEscapedPath = r.URL.EscapedPath()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"x"}`))
	}))
	defer srv.Close()

	spec := Specification{BaseUrl: srv.URL, ApiToken: "X"}
	// An id (derived from the experiment key) trying to traverse the path.
	_, _, err := spec.CreateMaintenanceWindow(context.Background(), types.CreateMaintenanceWindowRequest{Id: "exp/../../evil"})
	require.NoError(t, err)

	// The '/' must be percent-encoded so the id stays a single path segment.
	assert.Equal(t, "/api/settings/v2/maintenance/exp%2F..%2F..%2Fevil", gotEscapedPath)
}
