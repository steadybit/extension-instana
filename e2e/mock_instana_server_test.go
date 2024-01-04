package e2e

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

var Requests []string

func createMockInstanaServer() *httptest.Server {
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(fmt.Sprintf("httptest: failed to listen: %v", err))
	}
	server := httptest.Server{
		Listener: listener,
		Config: &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info().Str("path", r.URL.Path).Str("method", r.Method).Str("query", r.URL.RawQuery).Msg("Request received")
			Requests = append(Requests, fmt.Sprintf("%s-%s", r.Method, r.URL.Path))
			if strings.HasPrefix(r.URL.Path, "/api/events") && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(events())
			} else if strings.HasPrefix(r.URL.Path, "/api/infrastructure-monitoring/snapshots") && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(snapshots())
			} else if strings.HasPrefix(r.URL.Path, "/api/settings/v2/maintenance") && r.Method == http.MethodPut {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(maintenanceWindowCreated())
			} else if strings.HasPrefix(r.URL.Path, "/api/settings/v2/maintenance") && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})},
	}
	server.Start()
	log.Info().Str("url", server.URL).Msg("Started Mock-Server")
	return &server
}
func snapshots() []byte {
	return []byte(`{
    "items": [
        {
            "snapshotId": "snapshot-id-4711",
            "plugin": "kubernetesDeployment",
            "from": 1703284764000,
            "to": null,
            "tags": [
                "app.kubernetes.io/managed-by=Helm"
            ],
            "label": "test/test",
            "host": ""
        }
    ]
}`)
}

func events() []byte {
	return []byte(`[
    {
        "eventId": "XoDkHMssTRGiyEpPz337hQ",
        "start": 1703075293000,
        "end": 1703075923000,
        "type": "issue",
        "state": "open",
        "problem": "Condition [Ready]: Pod containers are not ready",
        "detail": "containers with unready status: [platform platform-port-splitter]",
        "severity": 10,
        "entityName": "Kubernetes Pod",
        "entityLabel": "platform/platform-d94f66f69-r72lf (pod)",
        "metrics": [],
        "entityType": "INFRASTRUCTURE",
        "fixSuggestion": "containers with unready status: [platform platform-port-splitter]",
        "snapshotId": "snapshot-id-4711"
    }
  ]`)
}

func maintenanceWindowCreated() []byte {
	return []byte(`{
        "id": "TST-1-47",
        "name": "Dev Deployment",
        "query": "entity.zone:dev-eu-central-1-eks",
        "scheduling": {
            "start": 1678954777825,
            "duration": {
                "amount": 20,
                "unit": "MINUTES"
            },
            "type": "ONE_TIME"
        },
        "paused": false,
        "lastUpdated": 1702373873467,
        "state": "EXPIRED",
        "validVersion": 1,
        "tagFilterExpression": null,
        "tagFilterExpressionEnabled": false,
        "occurrence": {
            "start": 1678954777825,
            "end": 1678955977825
        },
        "invalid": false,
        "applicationNames": []
    }`)
}
