/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/extension-instana/types"
	"io"
	"net/http"
	"time"
)

// Specification is the configuration specification for the extension. Configuration values can be applied
// through environment variables. Learn more through the documentation of the envconfig package.
// https://github.com/kelseyhightower/envconfig
type Specification struct {
	// The Instana Base Url, like 'https://unit-example.instana.io'
	BaseUrl string `json:"baseUrl" split_words:"true" required:"true"`
	// The Instana API Token
	ApiToken string `json:"apiToken" split_words:"true" required:"true"`
}

var (
	Config Specification
)

func ParseConfiguration() {
	err := envconfig.Process("steadybit_extension", &Config)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse configuration from environment.")
	}
}
func (s *Specification) GetSnapshotIds(_ context.Context, applicationPerspectiveId string) ([]string, error) {
	url := fmt.Sprintf("%s/api/infrastructure-monitoring/snapshots?query=entity.application.id:%s&size=20000", s.BaseUrl, applicationPerspectiveId)

	responseBody, response, err := s.do(url, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get snapshot-ids from Instana. Full response %+v", string(responseBody))
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Error().Int("code", response.StatusCode).Err(err).Msgf("Unexpected response %+v", string(responseBody))
		return nil, errors.New("unexpected response code")
	}

	var result types.SnapshotSearchResponse
	if responseBody != nil {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			log.Error().Err(err).Str("body", string(responseBody)).Msgf("Failed to parse body")
			return nil, err
		} else {
			if len(result.Items) == 20000 {
				log.Warn().Msgf("There are more than 20000 snapshots for application perspective %s. Only the first 20000 will be considered. You might miss events.", applicationPerspectiveId)
			}
			snapshotIds := make([]string, 0, len(result.Items))
			for _, snapshot := range result.Items {
				snapshotIds = append(snapshotIds, snapshot.SnapshotId)
			}
			return snapshotIds, nil
		}
	} else {
		log.Error().Err(err).Msgf("Empty response body")
		return nil, errors.New("empty response body")
	}
}

func (s *Specification) GetEvents(_ context.Context, from time.Time, to time.Time, eventTypeFilters []string) ([]types.Event, error) {
	url := fmt.Sprintf("%s/api/events?from=%d&to=%d", s.BaseUrl, from.UnixMilli(), to.UnixMilli())
	for _, eventTypeFilter := range eventTypeFilters {
		url = fmt.Sprintf("%s&eventTypeFilters=%s", url, eventTypeFilter)
	}

	responseBody, response, err := s.do(url, "GET", nil)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get events from Instana. Full response %+v", string(responseBody))
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Error().Int("code", response.StatusCode).Err(err).Msgf("Unexpected response %+v", string(responseBody))
		return nil, errors.New("unexpected response code")
	}

	var result []types.Event
	if responseBody != nil {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			log.Error().Err(err).Str("body", string(responseBody)).Msgf("Failed to parse body")
			return nil, err
		}
	}

	return result, err
}

func (s *Specification) GetApplicationPerspectives(_ context.Context, page int, pageSize int) (*types.ApplicationPerspectiveResponse, error) {
	url := fmt.Sprintf("%s/api/application-monitoring/applications?page=%d&pageSize=%d", s.BaseUrl, page, pageSize)

	responseBody, response, err := s.do(url, "GET", nil)
	if err != nil {
		log.Error().Int("page", page).Int("pageSize", pageSize).Err(err).Msgf("Failed to get application perspectives from Instana. Full response %+v", string(responseBody))
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Error().Int("code", response.StatusCode).Int("page", page).Int("pageSize", pageSize).Err(err).Msgf("Unexpected response %+v", string(responseBody))
		return nil, errors.New("unexpected response code")
	}

	var result types.ApplicationPerspectiveResponse
	if responseBody != nil {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			log.Error().Int("page", page).Int("pageSize", pageSize).Err(err).Str("body", string(responseBody)).Msgf("Failed to parse body")
			return nil, err
		}
		return &result, err
	} else {
		log.Error().Int("page", page).Int("pageSize", pageSize).Err(err).Msgf("Empty response body")
		return nil, errors.New("empty response body")
	}
}

func (s *Specification) CreateMaintenanceWindow(_ context.Context, maintenanceWindow types.CreateMaintenanceWindowRequest) (*string, *http.Response, error) {
	b, err := json.Marshal(maintenanceWindow)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to marshal request")
		return nil, nil, err
	}

	responseBody, response, err := s.do(fmt.Sprintf("%s/api/settings/v2/maintenance/%s", s.BaseUrl, maintenanceWindow.Id), "PUT", b)
	if err != nil {
		return nil, response, err
	}

	if response.StatusCode != 200 {
		log.Error().Int("code", response.StatusCode).Err(err).Msgf("Unexpected response %+v", string(responseBody))
		return nil, response, errors.New("unexpected response code")
	}

	var result types.CreateMaintenanceWindowRequest
	if responseBody != nil {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			log.Error().Err(err).Str("body", string(responseBody)).Msgf("Failed to parse response body")
			return nil, response, err
		}
	}

	return &result.Id, response, err
}

func (s *Specification) DeleteMaintenanceWindow(_ context.Context, maintenanceWindowId string) (*http.Response, error) {
	_, response, err := s.do(fmt.Sprintf("%s/api/settings/v2/maintenance/%s", s.BaseUrl, maintenanceWindowId), "DELETE", nil)
	return response, err
}

func (s *Specification) do(url string, method string, body []byte) ([]byte, *http.Response, error) {
	log.Debug().Str("url", url).Str("method", method).Msg("Requesting Instana API")
	if body != nil {
		log.Debug().Int("len", len(body)).Str("body", string(body)).Msg("Request body")
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create request")
		return nil, nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", fmt.Sprintf("apiToken %s", s.ApiToken))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to execute request")
		return nil, response, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msgf("Failed to close response body")
		}
	}(response.Body)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read body")
		return nil, response, err
	}

	return responseBody, response, err
}
