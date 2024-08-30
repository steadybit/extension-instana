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
	"github.com/steadybit/extension-kit/extruntime"
	_ "net/http/pprof" //allow pprof
)

func main() {
	extlogging.InitZeroLog()
	extbuild.PrintBuildInformation()
	extruntime.LogRuntimeInformation(zerolog.DebugLevel)
	config.ParseConfiguration()

	exthealth.SetReady(false)
	exthealth.StartProbes(8091)

	exthttp.RegisterHttpHandler("/", exthttp.GetterAsHandler(getExtensionList))
	discovery_kit_sdk.Register(extapplications.NewApplicationPerspectiveDiscovery())
	action_kit_sdk.RegisterAction(extevents.NewEventCheckAction())
	action_kit_sdk.RegisterAction(extmaintenance.NewCreateMaintenanceWindowAction())
	//extevents.RegisterEventListenerHandlers()

	action_kit_sdk.InstallSignalHandler()
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
