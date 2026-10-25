// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package submatview

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/proto/private/pbservice"
	"github.com/hashicorp/consul/proto/private/pbsubscribe"
)

func TestMaterializer_ResetKeepsViewUntilNewSnapshot(t *testing.T) {
	view := &fakeView{srvs: make(map[string]*pbservice.CheckServiceNode)}
	m := newMaterializer(hclog.NewNullLogger(), view, nil)
	require.NoError(t, m.updateView([]*pbsubscribe.Event{
		newEventServiceHealthRegister(4, 1, "srv1"),
		newEventServiceHealthRegister(5, 2, "srv1"),
	}, 5))

	m.reset()

	// The next subscription asks for a full snapshot...
	require.Equal(t, uint64(0), m.currentIndex())

	// ...while a query that times out before it arrives gets the previous view.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	result, err := m.query(ctx, 5)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, uint64(5), result.Index)
	require.Len(t, result.Value.(fakeResult).srvs, 2)

	// The new snapshot replaces the view.
	require.NoError(t, m.updateView([]*pbsubscribe.Event{
		newEventServiceHealthRegister(8, 1, "srv1"),
	}, 8))
	result, err = m.query(context.Background(), 5)
	require.NoError(t, err)
	require.Equal(t, uint64(8), result.Index)
	require.Len(t, result.Value.(fakeResult).srvs, 1)
}
