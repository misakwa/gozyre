package gozyre

/*
#include "zyre.h"
#include "czmq.h"

int wrap_set_endpoint(zyre_t *self, const char *ep) {
	return zyre_set_endpoint(self, "%s", ep);
}

void wrap_gossip_bind(zyre_t *self, const char *bind) {
	zyre_gossip_bind(self, "%s", bind);
}

void wrap_gossip_connect(zyre_t *self, const char *conn) {
	zyre_gossip_connect(self, "%s", conn);
}

void wrap_set_header(zyre_t *self, const char *name, const char *val) {
	zyre_set_header(self, name, "%s", val);
}
*/
import "C"

import (
	"time"
	"unsafe"
)

// Zyre struct wraps the C zyre_t struct
type Zyre struct {
	czyre *C.struct__zyre_t
}

// Force ZMQ init 
// Ensure ZMQ won't interact with GO signal handling mechanism.
func init() {
	C.zsys_init()
	C.zsys_handler_set(nil)
}

// Instruct ZMQ to exit and free all its resources.
func Exit() { C.zsys_shutdown() }

// New constructs a new node for peer-to-peer discovery
// Constructor, creates a new Zyre node. Note that until you start the
// node, it is silent and invisible to other nodes on the network.
// The node name is provided to other nodes during discovery. If you
// specify an empty string, Zyre generates a randomized node name from the UUID.
// Specify an empty string for the interface to have the system choose the default.
// Specify 0 for the port to use the default port specification from the c library
func New(name, iface string, port uint16, headers map[string]string, verbose bool) *Zyre {
	znode := &Zyre{}
	if name == "" {
		znode.czyre = C.zyre_new(nil)
	} else {
		cname := C.CString(name)
		defer C.free(unsafe.Pointer(cname))
		znode.czyre = C.zyre_new(cname)
	}
	for k, v := range headers {
		ckey := C.CString(k)
		cval := C.CString(v)
		C.wrap_set_header(znode.czyre, ckey, cval)
		C.free(unsafe.Pointer(ckey))
		C.free(unsafe.Pointer(cval))
	}
	if iface != "" {
		ciface := C.CString(iface)
		C.zyre_set_interface(znode.czyre, ciface)
		C.free(unsafe.Pointer(ciface))
	}
	if port != 0 {
		C.zyre_set_port(znode.czyre, C.int(port))
	}
	if verbose {
		C.zyre_set_verbose(znode.czyre)
	}
	return znode
}

// Start begins discovering and receiving messages from other nodes
func (n *Zyre) Start() error {
	ret := int(C.zyre_start(n.czyre))
	if ret != 0 {
		return ErrStartFailed
	}
	return nil
}

// Destroy shuts down the node. Any messages that are being sent or received
// will be discarded.
func (n *Zyre) Destroy() {
	C.zyre_destroy(&n.czyre)
	n.czyre = nil
}

// Name returns our node name, after successful initialization. By default
// is taken from the UUID and shortened.
func (n *Zyre) Name() string { return C.GoString(C.zyre_name(n.czyre)) }

// UUID returns our node UUID string, after successful initialization
func (n *Zyre) UUID() string { return C.GoString(C.zyre_uuid(n.czyre)) }

// Set beacon TCP ephemeral port to a well known value.
func (n *Zyre) SetBeaconPeerPort(port uint16) {
	C.zyre_set_beacon_peer_port(n.czyre, C.int(port))
}

// Old name of the above, deprecated.
func (n *Zyre) SetEphemeralPort(port uint16) {
	C.zyre_set_beacon_peer_port(n.czyre, C.int(port))
}

// SetEvasive sets the node evasiveness timeout. Default is 5 * time.Millisecond.
func (n *Zyre) SetEvasive(timeout time.Duration) {
	ctimeout := C.int(float64(timeout.Seconds()) * 1000)
	C.zyre_set_evasive_timeout(n.czyre, ctimeout)
}

// SetExpire sets the node expiration timeout. Default is 30 * time.Second
func (n *Zyre) SetExpire(expired time.Duration) {
	cexpired := C.int(expired.Seconds() * 1000)
	C.zyre_set_expired_timeout(n.czyre, cexpired)
}

// SetInterval sets the UDP beacon discovery interval. Default is instant
// beacon exploration followed by pinging every 1 * time.Second
func (n *Zyre) SetInterval(interval time.Duration) {
	cinterval := C.size_t(interval.Seconds() * 1000)
	C.zyre_set_interval(n.czyre, cinterval)
}

// Join joins the named group. After joining a group you can send messages to
// the group and all nodes in that group will receive them.
func (n *Zyre) Join(group string) {
	cgroup := C.CString(group)
	defer C.free(unsafe.Pointer(cgroup))
	C.zyre_join(n.czyre, cgroup)
}

// Leave leaves the named group.
func (n *Zyre) Leave(group string) {
	cgroup := C.CString(group)
	defer C.free(unsafe.Pointer(cgroup))
	C.zyre_leave(n.czyre, cgroup)
}

// Recv receives next event peers
func (n *Zyre) Recv() (*Event, error) { return NewEvent(n) }

// Whisper sends messages to a peers. Each of the messages is grouped into frames
func (n *Zyre) Whisper(name string, frames <-chan []byte) error {
	zmsg, err := bytesToZmsg(frames)
	if err == nil {
		cname := C.CString(name)
		defer C.free(unsafe.Pointer(cname))
		C.zyre_whisper(n.czyre, cname, &zmsg)
	}
	return err
}

// Shout sends messages to a group. Each of the messages is grouped into frames
func (n *Zyre) Shout(name string, frames <-chan []byte) error {
	zmsg, err := bytesToZmsg(frames)
	if err == nil {
		cname := C.CString(name)
		defer C.free(unsafe.Pointer(cname))
		C.zyre_shout(n.czyre, cname, &zmsg)
	}
	return err
}

// Groups returns list groups that this node belongs to
func (n *Zyre) Groups() []string {
	return zListToSlice(C.zyre_own_groups(n.czyre), true)
}

// PeerGroups returns list of groups known through peers
func (n *Zyre) PeerGroups() []string {
	return zListToSlice(C.zyre_peer_groups(n.czyre), true)
}

// Peers returns list of peers in a specified group.
// It will return a list of all current peers if an empty string is specified
func (n *Zyre) Peers(group string) []string {
	if group == "" {
		return zListToSlice(C.zyre_peers(n.czyre), true)
	}
	cgroup := C.CString(group)
	defer C.free(unsafe.Pointer(cgroup))
	return zListToSlice(C.zyre_peers_by_group(n.czyre, cgroup), true)
}

// PeerHeader returns a peer header value.
// It will return an empty string if peer or header doesn't exist.
func (n *Zyre) PeerHeader(uuid, name string) string {
	cpeer := C.CString(uuid)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cpeer))
	defer C.free(unsafe.Pointer(cname))
	header := C.zyre_peer_header_value(n.czyre, cpeer, cname)
	defer C.free(unsafe.Pointer(header))
	return C.GoString(header)
}

// PeerAddress returns the endpoint of a connected peer.
func (n *Zyre) PeerAddress(uuid string) string {
	cpeer := C.CString(uuid)
	defer C.free(unsafe.Pointer(cpeer))
	caddr := C.zyre_peer_address(n.czyre, cpeer)
	defer C.free(unsafe.Pointer(caddr))
	return C.GoString(caddr)
}

// Gossip sets up gossip discovery and switches from using the default UDP beaconing
// It will bind and connect to the gossip endpoint
func (n *Zyre) Gossip(endpoint, hub string) error {
	cep := C.CString(endpoint)
	defer C.free(unsafe.Pointer(cep))
	ret := int(C.wrap_set_endpoint(n.czyre, cep))
	if ret != 0 {
		return ErrInvalidEndpoint
	}
	// Bind and connect to gossip hub
	chub := C.CString(hub)
	defer C.free(unsafe.Pointer(chub))
	C.wrap_gossip_bind(n.czyre, chub)
	C.wrap_gossip_connect(n.czyre, chub)
	return nil
}

// Stop signals to other nodes that this node will go away
func (n *Zyre) Stop() { C.zyre_stop(n.czyre) }

// Socket returns the socket, used by ZYRE.
func (n *Zyre) Socket() *C.struct__zsock_t { return C.zyre_socket(n.czyre) }

