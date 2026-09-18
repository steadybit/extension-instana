/*
 * Copyright 2023 steadybit GmbH. All rights reserved.
 */

package main

import (
	"github.com/rs/zerolog"
	"github.com/steadybit/action-kit/go/action_kit_api/v2"
	"github.com/steadybit/action-kit/go/action_kit_sdk"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/extension-instana/config"
	"github.com/steadybit/extension-instana/extapplications"
	"github.com/steadybit/extension-instana/extevents"
	"github.com/steadybit/extension-instana/extmaintenance"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/exthealth"
	"github.com/steadybit/extension-kit/exthttp"
	"github.com/steadybit/extension-kit/extlogging"
	"github.com/steadybit/extension-kit/extotel"
	"github.com/steadybit/extension-kit/extruntime"
	"github.com/steadybit/extension-kit/extsignals"
)

func main() {
	extlogging.InitZeroLog()

	// Export OpenTelemetry traces when an OTLP endpoint is configured, so an
	// operator debugging a slow or timing-out action can see what happened inside
	// this extension. Off, and free, until OTEL_EXPORTER_OTLP_ENDPOINT is set —
	// see the extension-kit README for the full set of variables.
	extotel.InitOpenTelemetry()
	extbuild.PrintBuildInformation()
	extruntime.LogRuntimeInformation(zerolog.DebugLevel)
	config.ParseConfiguration()
	config.ValidateConfiguration()

	exthealth.SetReady(false)
	exthealth.StartProbes(8091)

	discovery_kit_sdk.Register(extapplications.NewApplicationPerspectiveDiscovery())
	action_kit_sdk.RegisterAction(extevents.NewEventCheckAction())
	action_kit_sdk.RegisterAction(extmaintenance.NewCreateMaintenanceWindowAction())
	//extevents.RegisterEventListenerHandlers()

	exthttp.RegisterRevisionedHandler("/", getExtensionList)

	extsignals.ActivateSignalHandlers()
	action_kit_sdk.RegisterCoverageEndpoints()
	exthealth.SetReady(true)

	exthttp.Listen(exthttp.ListenOpts{
		Port: 8090,
	})
}

// ExtensionListResponse exists to merge the possible root path responses supported by the
// various extension kits. In this case, the response for ActionKit, DiscoveryKit and EventKit.
type ExtensionListResponse struct {
	action_kit_api.ActionList       `json:",inline"`
	discovery_kit_api.DiscoveryList `json:",inline"`
}

func getExtensionList() ExtensionListResponse {
	return ExtensionListResponse{
		ActionList:    action_kit_sdk.GetActionList(),
		DiscoveryList: discovery_kit_sdk.GetDiscoveryList(),
		//EventListenerList: extevents.GetEventListenerList(),
	}
}
