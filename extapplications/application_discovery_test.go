// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package extapplications

import (
	"context"
	"errors"
	"github.com/steadybit/extension-instana/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type instanaApiMock struct {
	mock.Mock
}

func (m *instanaApiMock) GetApplicationPerspectives(ctx context.Context, page int, pageSize int) (*types.ApplicationPerspectiveResponse, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ApplicationPerspectiveResponse), args.Error(1)
}

func TestIterateThroughMonitorsResponses(t *testing.T) {
	// Given
	mockedApi := new(instanaApiMock)
	page1 := types.ApplicationPerspectiveResponse{
		Items: []types.ApplicationPerspective{
			{
				Id:    "id1",
				Label: "name1",
			},
		},
		Links: map[string]string{
			"next": "next",
		},
	}
	page2 := types.ApplicationPerspectiveResponse{
		Items: []types.ApplicationPerspective{
			{
				Id:    "id2",
				Label: "name2",
			},
		},
		Links: map[string]string{},
	}

	mockedApi.On("GetApplicationPerspectives", mock.Anything, 1, mock.Anything).Return(&page1, nil)
	mockedApi.On("GetApplicationPerspectives", mock.Anything, 2, mock.Anything).Return(&page2, nil)

	// When
	monitors := getAllApplicationPerspectives(context.Background(), mockedApi)

	// Then
	require.Len(t, monitors, 2)
	require.Equal(t, "id1", monitors[0].Id)
	require.Equal(t, "name1", monitors[0].Label)
	require.Equal(t, "id2", monitors[1].Id)
	require.Equal(t, "name2", monitors[1].Label)
	mockedApi.AssertNumberOfCalls(t, "GetApplicationPerspectives", 2)
}

func TestErrorResponseReturnsIntermediateResult(t *testing.T) {
	// Given
	mockedApi := new(instanaApiMock)
	page1 := types.ApplicationPerspectiveResponse{
		Items: []types.ApplicationPerspective{
			{
				Id:    "id1",
				Label: "name1",
			},
		},
		Links: map[string]string{
			"next": "next",
		},
	}

	mockedApi.On("GetApplicationPerspectives", mock.Anything, 1, mock.Anything).Return(&page1, nil)
	mockedApi.On("GetApplicationPerspectives", mock.Anything, 2, mock.Anything).Return(nil, errors.New("oops"))

	// When
	monitors := getAllApplicationPerspectives(context.Background(), mockedApi)

	// Then
	require.Len(t, monitors, 1)
	require.Equal(t, "id1", monitors[0].Id)
	require.Equal(t, "name1", monitors[0].Label)
	mockedApi.AssertNumberOfCalls(t, "GetApplicationPerspectives", 2)
}
