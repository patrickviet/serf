package serf2

import (
	"github.com/hashicorp/memberlist"
	"time"
)

// Config is the configuration for creating a Serf instance.
type Config struct {
	// The name of this node. This must be unique in the cluster.
	NodeName string

	// The role for this node, if any. This is used to differentiate
	// between perhaps different members of a Serf. For example, you might
	// have a "load-balancer" role and a "web" role part of the same cluster.
	// When new nodes are added, the load balancer wants to know (so it
	// must be part of the cluster), but it doesn't want to add other load
	// balancers to the rotation, so it checks if the added nodes are "web".
	Role string

	// EventCh is a channel that receives all the Serf events. The events
	// are sent on this channel in proper ordering. Care must be taken that
	// this channel doesn't block, either by processing the events quick
	// enough or buffering the channel, otherwise it can block state updates
	// within Serf itself. If no EventCh is specified, no events will be fired,
	// but point-in-time snapshots of members can still be retrieved by
	// calling Members on Serf.
	EventCh chan<- Event

	// LeaveTimeout is the amount of time to wait for a node to gracefully
	// leave. If this is not set, a timeout of 120 seconds will be set.
	LeaveTimeout time.Duration

	// BroadcastTimeout is the amount of time to wait for a broadcast
	// message to be sent to the cluster. Broadcast messages are used for
	// things like leave messages and force remove messages. If this is not
	// set, a timeout of 5 seconds will be set.
	BroadcastTimeout time.Duration

	// The settings below relate to Serf's event coalescence feature. Serf
	// is able to coalesce multiple events into single events in order to
	// reduce the amount of noise that is sent along the EventCh. For example
	// if five nodes quickly join, the EventCh will be sent one EventMemberJoin
	// containing the five nodes rather than five individual EventMemberJoin
	// events. Coalescence can mitigate potential flapping behavior.
	//
	// Coalescence is disabled by default and can be enabled by setting
	// CoalescePeriod.
	//
	// CoalescePeriod specifies the time duration to coalesce events.
	// For example, if this is set to 5 seconds, then all events received
	// within 5 seconds that can be coalesced will be.
	//
	// QuiescentPeriod specifies the duration of time where if no events
	// are received, coalescence immediately happens. For example, if
	// CoalscePeriod is set to 10 seconds but QuiscentPeriod is set to 2
	// seconds, then the events will be coalesced and dispatched if no
	// new events are received within 2 seconds of the last event. Otherwise,
	// every event will always be delayed by at least 10 seconds.
	CoalescePeriod  time.Duration
	QuiescentPeriod time.Duration

	// The settings below relate to Serf keeping track of recently
	// failed/left nodes and attempting reconnects.
	//
	// ReapInterval is the interval when the reaper runs. If this is not
	// set (it is zero), it will be set to a reasonable default.
	//
	// ReconnectInterval is the interval when we attempt to reconnect
	// to failed nodes. If this is not set (it is zero), it will be set
	// to a reasonable default.
	//
	// ReconnectTimeout is the amount of time to attempt to reconnect to
	// a failed node before giving up and considering it completely gone.
	//
	// TombstoneTimeout is the amount of time to keep around nodes
	// that gracefully left as tombstones for syncing state with other
	// Serf nodes.
	ReapInterval      time.Duration
	ReconnectInterval time.Duration
	ReconnectTimeout  time.Duration
	TombstoneTimeout  time.Duration

	// MemberlistConfig is the memberlist configuration that Serf will
	// use to do the underlying membership management and gossip. Some
	// fields in the MemberlistConfig will be overwritten by Serf no
	// matter what:
	//
	//   * Name - This will always be set to the same as the NodeName
	//     in this configuration.
	//
	//   * Events - Serf uses a custom event delegate.
	//
	//   * Delegate - Serf uses a custom delegate.
	//
	MemberlistConfig *memberlist.Config
}

// DefaultConfig returns a Config struct that contains reasonable defaults
// for most of the configurations.
func DefaultConfig() *Config {
	return &Config{
		MemberlistConfig: memberlist.DefaultConfig(),
	}
}
