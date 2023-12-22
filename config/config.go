/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
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

func ValidateConfiguration() {
	// You may optionally validate the configuration here.
}
