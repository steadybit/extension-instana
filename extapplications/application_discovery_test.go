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

func (m *instanaApiMock) GetApplicationPerspectives(ctx context.Context, page int, pageSize int) ([]types.ApplicationPerspective, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.ApplicationPerspective), args.Error(1)
}

func TestIterateThroughMonitorsResponses(t *testing.T) {
	// Given
	mockedApi := new(instanaApiMock)
	page1 := []types.ApplicationPerspective{
		{
			Id:    "id1",
			Label: "name1",
		},
	}
	page2 := []types.ApplicationPerspective{
		{
			Id:    "id2",
			Label: "name2",
		},
	}
	var page3 []types.ApplicationPerspective

	mockedApi.On("GetApplicationPerspectives", mock.Anything, 0, mock.Anything).Return(page1, nil)
	mockedApi.On("GetApplicationPerspectives", mock.Anything, 1, mock.Anything).Return(page2, nil)
	mockedApi.On("GetApplicationPerspectives", mock.Anything, 2, mock.Anything).Return(page3, nil)

	// When
	monitors := getAllApplicationPerspectives(context.Background(), mockedApi)

	// Then
	require.Len(t, monitors, 2)
	require.Equal(t, "id1", monitors[0].Id)
	require.Equal(t, "name1", monitors[0].Label)
	require.Equal(t, "id2", monitors[1].Id)
	require.Equal(t, "name2", monitors[1].Label)
	mockedApi.AssertNumberOfCalls(t, "GetApplicationPerspectives", 3)
}

func TestErrorResponseReturnsIntermediateResult(t *testing.T) {
	// Given
	mockedApi := new(instanaApiMock)
	page1 := []types.ApplicationPerspective{
		{
			Id:    "id1",
			Label: "name1",
		},
	}

	mockedApi.On("GetApplicationPerspectives", mock.Anything, 0, mock.Anything).Return(page1, nil)
	mockedApi.On("GetApplicationPerspectives", mock.Anything, 1, mock.Anything).Return(nil, errors.New("oops"))

	// When
	monitors := getAllApplicationPerspectives(context.Background(), mockedApi)

	// Then
	require.Len(t, monitors, 1)
	require.Equal(t, "id1", monitors[0].Id)
	require.Equal(t, "name1", monitors[0].Label)
	mockedApi.AssertNumberOfCalls(t, "GetApplicationPerspectives", 2)
}
