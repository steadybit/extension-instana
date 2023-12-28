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
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		})},
	}
	server.Start()
	log.Info().Str("url", server.URL).Msg("Started Mock-Server")
	return &server
}

func events() []byte {
	return []byte(`[
    {
        "eventId": "95wKHid1Qh2o8tBjGB0wcA",
        "start": 1703677881000,
        "end": 1703677881000,
        "type": "change",
        "state": "closed",
        "problem": "offline",
        "detail": "JVM on Host ip-10-10-81-117.eu-central-1.compute.internal",
        "severity": -1,
        "entityName": "JVM",
        "entityLabel": "Unknown",
        "entityType": "INFRASTRUCTURE",
        "fixSuggestion": "JVM on Host ip-10-10-81-117.eu-central-1.compute.internal",
        "snapshotId": "-QEDb2D3jtvz7vYJOMYcxTSEyXQ"
    }
]`)
}
