package gozyre

/*
#include "zyre.h"
*/
import "C"

import (
	"unsafe"
)

// EventType describes network events
type EventType string

const (
	// JoinEvent indicates a node has joined a specific group
	JoinEvent EventType = "JOIN"

	// ExitEvent indicates a node has left the network
	ExitEvent = "EXIT"

	// EnterEvent indicates a node hash entered the network
	EnterEvent = "ENTER"

	// LeaveEvent indicates a peer has left a specific group
	LeaveEvent = "LEAVE"

	// StopEvent indicates that a node will go away
	StopEvent = "STOP"

	// EvasiveEvent indicates a node is being quiet for too long
	EvasiveEvent = "EVASIVE"

	// WhisperEvent indicates a peer has sent a message to this node
	WhisperEvent = "WHISPER"

	// ShoutEvent indicates a peer has sent a message to a group
	ShoutEvent = "SHOUT"

	// UnknownEvent indicates a message that is unsupported or unknown
	UnknownEvent = "UNKNOWN"
)

// Event describes a network event
type Event struct {
	czyreEvent *C.struct__zyre_event_t
}

// NewEvent creates a new Event object
func NewEvent(node *Zyre) (*Event, error) {
	cevent := C.zyre_event_new(node.czyre)
	if cevent == nil {
		return nil, ErrNodeInterrupted
	}
	return &Event{czyreEvent: cevent}, nil
}

// Type returns the event type
func (ev *Event) Type() EventType {
	evt := C.GoString(C.zyre_event_type(ev.czyreEvent))
	switch evt {
	case "ENTER":
		return EnterEvent
	case "JOIN":
		return JoinEvent
	case "EXIT":
		return ExitEvent
	case "LEAVE":
		return LeaveEvent
	case "STOP":
		return StopEvent
	case "EVASIVE":
		return EvasiveEvent
	case "WHISPER":
		return WhisperEvent
	case "SHOUT":
		return ShoutEvent
	default:
		return UnknownEvent
	}
}

// UUID returns the uuid of the peer responsible for this event
func (ev *Event) UUID() string { return C.GoString(C.zyre_event_peer_uuid(ev.czyreEvent)) }

// Peer returns the name of the peer sending the event
func (ev *Event) Peer() string { return C.GoString(C.zyre_event_peer_name(ev.czyreEvent)) }

// Group returns the group name for a ShoutEvent
func (ev *Event) Group() string { return C.GoString(C.zyre_event_group(ev.czyreEvent)) }

// Address returns the address of the sending peer ip address
func (ev *Event) Address() string { return C.GoString(C.zyre_event_peer_addr(ev.czyreEvent)) }

// Headers returns the event headers. These are the headers from the node that
// initiated this event.
// XXX: Do we need to lock the headers? The underlying implementation uses a
// cursor to iterate through the items and may introduce race conditions.
func (ev *Event) Headers() map[string]string {
	czev := ev.czyreEvent
	if (czev == nil) {
		return nil
	}
	czhdr := C.zyre_event_headers(czev)
	if (czhdr == nil) {
		return nil
	}
	return zHashToMap(czhdr, false)
}

// Header returns the header value from the header name
func (ev *Event) Header(name string) string {
	hname := C.CString(name)
	defer C.free(unsafe.Pointer(hname))
	return C.GoString(C.zyre_event_header(ev.czyreEvent, hname))
}

// Message returns the frames of data from the message content, each frame as a byte slice
// XXX: Do we need to lock the channel so that only one consumer can call it?
// The underlying implementation iterates through the message frames using a
// cursor and can introduce race conditions.
func (ev *Event) Message() <-chan []byte { return zMsgToBytes(C.zyre_event_msg(ev.czyreEvent)) }

// Destroy removes and frees up the event object memory. The event should not
// be used after callling this method.
func (ev *Event) Destroy() {
	C.zyre_event_destroy(&ev.czyreEvent)
	ev.czyreEvent = nil
}
