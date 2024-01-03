// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package extapplications

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/steadybit/discovery-kit/go/discovery_kit_api"
	"github.com/steadybit/discovery-kit/go/discovery_kit_sdk"
	"github.com/steadybit/extension-instana/config"
	"github.com/steadybit/extension-instana/types"
	"github.com/steadybit/extension-kit/extbuild"
	"github.com/steadybit/extension-kit/extutil"
	"time"
)

type applicationPerspectiveDiscovery struct {
}

var (
	_ discovery_kit_sdk.TargetDescriber    = (*applicationPerspectiveDiscovery)(nil)
	_ discovery_kit_sdk.AttributeDescriber = (*applicationPerspectiveDiscovery)(nil)
)

func NewApplicationPerspectiveDiscovery() discovery_kit_sdk.TargetDiscovery {
	discovery := &applicationPerspectiveDiscovery{}
	return discovery_kit_sdk.NewCachedTargetDiscovery(discovery,
		discovery_kit_sdk.WithRefreshTargetsNow(),
		discovery_kit_sdk.WithRefreshTargetsInterval(context.Background(), 1*time.Minute),
	)
}
func (d *applicationPerspectiveDiscovery) Describe() discovery_kit_api.DiscoveryDescription {
	return discovery_kit_api.DiscoveryDescription{
		Id:         ApplicationPerspectiveTargetId,
		RestrictTo: extutil.Ptr(discovery_kit_api.LEADER),
		Discover: discovery_kit_api.DescribingEndpointReferenceWithCallInterval{
			CallInterval: extutil.Ptr("1m"),
		},
	}
}

func (d *applicationPerspectiveDiscovery) DescribeTarget() discovery_kit_api.TargetDescription {
	return discovery_kit_api.TargetDescription{
		Id:       ApplicationPerspectiveTargetId,
		Label:    discovery_kit_api.PluralLabel{One: "Instana application perspective", Other: "Instana application perspectives"},
		Category: extutil.Ptr("monitoring"),
		Version:  extbuild.GetSemverVersionStringOrUnknown(),
		Icon:     extutil.Ptr(applicationPerspectiveIcon),
		Table: discovery_kit_api.Table{
			Columns: []discovery_kit_api.Column{
				{Attribute: "steadybit.label"},
			},
			OrderBy: []discovery_kit_api.OrderBy{
				{
					Attribute: "steadybit.label",
					Direction: "ASC",
				},
			},
		},
	}
}

func (d *applicationPerspectiveDiscovery) DescribeAttributes() []discovery_kit_api.AttributeDescription {
	return []discovery_kit_api.AttributeDescription{
		{
			Attribute: "instana.application-perspective.name",
			Label: discovery_kit_api.PluralLabel{
				One:   "Instana application perspective name",
				Other: "Instana application perspective names",
			},
		},
	}
}

func (d *applicationPerspectiveDiscovery) DiscoverTargets(ctx context.Context) ([]discovery_kit_api.Target, error) {
	return getAllApplicationPerspectives(ctx, &config.Config), nil
}

type GetApplicationPerspectivesApi interface {
	GetApplicationPerspectives(ctx context.Context, page int, pageSize int) (*types.ApplicationPerspectiveResponse, error)
}

func getAllApplicationPerspectives(ctx context.Context, api GetApplicationPerspectivesApi) []discovery_kit_api.Target {
	result := make([]discovery_kit_api.Target, 0, 500)

	pageSize := 100
	page := 1

	start := time.Now()
	for {
		log.Debug().Int("page", page).Msg("Fetch application perspectives from Instana")
		response, err := api.GetApplicationPerspectives(ctx, page, pageSize)
		if err != nil {
			log.Err(err).Msgf("Failed to get application perspectives from Instana for page %d and page size %d.", page, pageSize)
			return result
		}

		perspectives := response.Items
		for _, perspective := range perspectives {
			result = append(result, toTarget(perspective))
		}

		_, hasMore := response.Links["next"]
		if len(perspectives) == 0 || !hasMore {
			// end of list reached
			break
		}
		page = page + 1
	}
	log.Debug().Msgf("Discovery took %s, returning %d application perspectives.", time.Since(start), len(result))
	return result
}

func toTarget(perspective types.ApplicationPerspective) discovery_kit_api.Target {
	id := perspective.Id
	label := perspective.Label

	attributes := make(map[string][]string)
	attributes["steadybit.label"] = []string{label}
	attributes["instana.application.label"] = []string{label}
	attributes["instana.application.id"] = []string{id}

	return discovery_kit_api.Target{
		Id:         id,
		Label:      label,
		TargetType: ApplicationPerspectiveTargetId,
		Attributes: attributes,
	}
}
