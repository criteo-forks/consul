// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package subscribe

import (
	"context"
	"errors"
	"sync"
	"time"
)

const (
	// maxLeaderContactAge is how long a follower may go without hearing from
	// the raft leader before it stops serving subscriptions.
	maxLeaderContactAge = 10 * time.Second

	// stepDownGrace is how long after losing leadership a server keeps
	// serving subscriptions regardless of its last contact time, which only
	// becomes meaningful once it hears from the new leader.
	stepDownGrace = 5 * time.Second

	leaderContactCheckInterval = time.Second
)

var errNoLeaderContact = errors.New("server has no recent contact with the raft leader")

// LeaderContactWatcher tracks whether this server has heard from the raft
// leader recently. A follower cut from the leader keeps serving stale data
// without any error, so subscriptions end once it is out of contact and the
// client resubscribes elsewhere. One watcher serves all the subscriptions of
// a server: each stream waits on the context returned by Context instead of
// polling on its own.
type LeaderContactWatcher struct {
	isLeader    func() bool
	lastContact func() time.Time

	lock     sync.Mutex
	ctx      context.Context
	cancel   context.CancelCauseFunc
	leaderAt time.Time // last check that found this server leader
}

// NewLeaderContactWatcher returns a watcher for a server whose raft state is
// read through isLeader and lastContact. Run must be called to keep it
// current.
func NewLeaderContactWatcher(isLeader func() bool, lastContact func() time.Time) *LeaderContactWatcher {
	w := &LeaderContactWatcher{isLeader: isLeader, lastContact: lastContact}
	w.ctx, w.cancel = context.WithCancelCause(context.Background())
	w.check(time.Now())
	return w
}

// Context returns a context that is canceled with errNoLeaderContact once
// this server loses contact with the raft leader. It is already canceled
// while contact is lost, and replaced by a new one when contact is back.
func (w *LeaderContactWatcher) Context() context.Context {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.ctx
}

// Run checks leader contact every second until ctx is canceled.
func (w *LeaderContactWatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(leaderContactCheckInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			w.check(now)
		}
	}
}

func (w *LeaderContactWatcher) check(now time.Time) {
	w.lock.Lock()
	defer w.lock.Unlock()

	fresh := true
	switch {
	case w.isLeader():
		w.leaderAt = now
	case now.Sub(w.leaderAt) < stepDownGrace:
	default:
		fresh = now.Sub(w.lastContact()) <= maxLeaderContactAge
	}

	switch {
	case !fresh && w.ctx.Err() == nil:
		w.cancel(errNoLeaderContact)
	case fresh && w.ctx.Err() != nil:
		w.ctx, w.cancel = context.WithCancelCause(context.Background())
	}
}
