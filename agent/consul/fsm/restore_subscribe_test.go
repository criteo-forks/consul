// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package fsm

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/agent/consul/state"
	"github.com/hashicorp/consul/agent/consul/stream"
	"github.com/hashicorp/consul/sdk/testutil"
)

// Restore takes stateLock, then the publisher lock to refresh the topics.
// Subscribe takes the publisher lock, then stateLock when its snapshot handler
// calls State() on a cache miss. Force a restore in between: it must still
// complete.
func TestFSM_RestoreDuringSubscribeSnapshot(t *testing.T) {
	t.Parallel()

	publisher := stream.NewEventPublisher(10 * time.Second)
	fsm := NewFromDeps(Deps{
		Logger: testutil.Logger(t),
		NewStateStore: func() *state.Store {
			return state.NewStateStore(nil)
		},
		Publisher:      publisher,
		StorageBackend: newStorageBackend(t, nil),
	})

	snap, err := fsm.Snapshot()
	require.NoError(t, err)
	defer snap.Release()
	sink := &MockSink{bytes.NewBuffer(nil), false}
	require.NoError(t, snap.Persist(sink))

	// Like the FSM's own snapshot handlers, this one calls State() while
	// Subscribe holds the publisher lock, but only when told to.
	inHandler, readState := make(chan struct{}), make(chan struct{})
	topic := stream.StringTopic("restore-test")
	require.NoError(t, publisher.RegisterHandler(topic,
		func(stream.SubscribeRequest, stream.SnapshotAppender) (uint64, error) {
			close(inHandler)
			<-readState
			fsm.State()
			return 1, nil
		}, false))
	go publisher.Subscribe(&stream.SubscribeRequest{Topic: topic, Subject: stream.StringSubject("x")})
	<-inHandler

	restored := make(chan error, 1)
	go func() { restored <- fsm.Restore(sink) }()

	// Let the handler read the state once the restore holds stateLock, or
	// after a second if it only holds it briefly.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) && fsm.stateLock.TryRLock() {
		fsm.stateLock.RUnlock()
		time.Sleep(time.Millisecond)
	}
	close(readState)

	select {
	case err := <-restored:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("restore deadlocked with a subscribe building its snapshot")
	}
}
