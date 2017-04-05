// Package gozyre provides ...
package gozyre

/*
#cgo !windows pkg-config: libczmq libzmq libzyre
#cgo windows LDFLAGS: -lws2_32 -liphlpapi -lrpcrt4 -lzmq -lczmq -lzyre
#cgo windows CFLAGS: -Wno-pedantic-ms-format -DLIBCZMQ_EXPORTS -DZMQ_DEFINED_STDINT -DLIBCZMQ_EXPORTS -DZMQ_BUILD_DRAFT_API -DZYRE_EXPORTS

#include "czmq.h"
#include "zyre.h"
#include "zyre_event.h"
*/

import "C"

import (
	"errors"
	"net"
	"time"
	"unsafe"
)

type SendType uint8

const (
	PeerMessage SendType = iota
	GroupMessage
)

type Message struct {
	type_    SendType
	headers  map[string]string
	content  []byte
	endpoint string
	size     uint64
	peer     string
	group    string
}

// Zyre struct wraps the C zyre_t struct
type Zyre struct {
	zyre *C.struct__zyre_t
}

// New constructs a new node for peer-to-peer discovery
func New(name string, port uint16, iface *net.Interface) *Zyre {
}

func (self *Zyre) Start() error {
}

func (self *Zyre) Destroy() {
}

func (self *Zyre) Name() string {
}

func (self *Zyre) Uuid() string {
}

// Version returns the underlying zyre version information in major.minor.patch
func (self *Zyre) Version() []byte {
}

func (self *Zyre) SetVerbose() {
}

func (self *Zyre) SetEvasiveTimeout(timeout time.Duration) {
}

func (self *Zyre) SetExpiredTimeout(expired time.Duration) {
}

func (self *Zyre) SetInterval(interval time.Duration) {
}

func (self *Zyre) Join(peerOrGroup string) {
}

func (self *Zyre) Leave(peerOrGroup string) {
}

func (self *Zyre) Messages() <-chan *Message {
}

// Send sends a message to peer or group. You can specify the send type.
func (self *Zyre) Send(name string, type_ uint8, message []byte, args ...string) error {
}

// Peers returns list of peers. You can specify group name to filter peers by group
func (self *Zyre) Peers(group string) []string {
}
func (self *Zyre) MyGroup() []string {
}
func (self *Zyre) PeerGroups() []string {
}
func (self *Zyre) Address(peer string) error {
}

func (self *Zyre) Header(peer, header string) string {
}

// Event struct wraps the C zyre_event_t struct
type Event struct {
	zevent *C.struct__zyre_event_t
}
