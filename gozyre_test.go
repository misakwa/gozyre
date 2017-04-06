package gozyre_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/misakwa/gozyre"
)

var verbose bool

func init() {
	_, ok := os.LookupEnv("GOZYRE_TEST_VERBOSE")
	if ok {
		verbose = true
	}
}

func TestEmptyNode(t *testing.T) {
	node := gozyre.New("", "", 0, map[string]string{}, verbose)
	defer node.Destroy()
	if node.Name() == "" {
		t.Error("Name should not be empty")
	}
	if node.UUID() == "" {
		t.Error("UUID should not be empty")
	}
}

// @see https://github.com/zeromq/zyre/blob/v2.0.0/src/zyre.c#L589
func TestGoZyre(t *testing.T) {
	node2 := gozyre.New("node2", "", 0, map[string]string{}, verbose)
	defer node2.Destroy()
	node1 := gozyre.New("node1", "", 0, map[string]string{"X-HELLO": "World"}, verbose)
	defer node1.Destroy()
	if node1.Name() != "node1" {
		t.Error("Expected node1")
	}

	if node2.Name() != "node2" {
		t.Error("Expected node2")
	}

	var (
		event *gozyre.Event
		err   error
	)
	err = node1.Gossip("inproc://zyre-node1", "inproc://gossip-hub")
	if err != nil {
		t.Error("Failed to setup node gossip")
	}
	node1.Start()
	err = node2.Gossip("inproc://zyre-node1", "inproc://gossip-hub")
	if err == nil {
		t.Error("Can't use same node endpoint more than once")
	}
	err = node2.Gossip("inproc://zyre-node2", "inproc://gossip-hub")
	if err != nil {
		t.Error("Able to use free endpoint for node gossip")
	}
	node2.Start()

	if node1.UUID() == node2.UUID() {
		t.Error("Nodes should have different UUIDS")
	}
	node1.Join("GLOBAL")
	node2.Join("GLOBAL")

	// Give nodes time to fine each other
	time.Sleep(10 * time.Millisecond)

	peers := node1.Peers("")
	if len(peers) != 1 {
		t.Error("Got incorrect peer count")
	}
	node1.Join("node1 group of one")
	node2.Join("node2 group of one")

	time.Sleep(10 * time.Millisecond)

	if len(node1.Groups()) != 2 {
		t.Error("Expected exactly 2 groups")
	}

	if len(node1.PeerGroups()) != 2 {
		t.Error("Expected exactly 2 peer groups")
	}

	if node2.PeerHeader(node1.UUID(), "X-HELLO") != "World" {
		t.Error("Expected X-HELLO World header")
	}

	mch := make(chan []byte)
	go func() {
		defer close(mch)
		mch <- []byte("Hello, World")
	}()
	// node1 shouts message to global
	node1.Shout("GLOBAL", mch)

	// node2 will receive ENTER, JOIN, SHOUT
	event, err = node2.Recv()
	if err == nil {
		defer event.Destroy()
	} else {
		t.Errorf("Failed to recv: %s", err)
	}
	et := event.Type()
	if et != gozyre.EnterEvent {
		t.Errorf("Expected ENTER event, got '%s' instead", et)
	}
	peer := event.Peer()
	if event.Peer() != "node1" {
		t.Errorf("Expected 'node1', got '%s' instead", peer)
	}
	if event.UUID() != node1.UUID() {
		t.Error("Expected event uuid to be node's uuid")
	}

	headers := map[string]string{
		"X-HELLO": "World",
	}
	if !reflect.DeepEqual(event.Headers(), headers) {
		t.Error("Expected headers to be equal")
	}
	evAddr := event.Address()
	peerAddr := node2.PeerAddress(node1.UUID())
	if evAddr != peerAddr {
		t.Errorf("Address should be endpoint of sending node: %s != %s", evAddr, peerAddr)
	}
	event.Destroy()

	event, err = node2.Recv()
	if event.Type() != gozyre.JoinEvent {
		t.Error("Expected JOIN event")
	}
	event.Destroy()

	event, err = node2.Recv()
	if event.Type() != gozyre.JoinEvent {
		t.Error("Expected JOIN event")
	}
	event.Destroy()

	event, err = node2.Recv()
	if event.Type() != gozyre.ShoutEvent {
		t.Error("Expected SHOUT event")
	}
	data := make([]byte, 0)
	for frame := range event.Message() {
		data = append(frame)
	}
	event.Destroy()

	if !bytes.Equal(data, []byte("Hello, World")) {
		t.Error("Expected event data")
	}

	node2.Stop()
	event, err = node2.Recv()
	if event.Type() != gozyre.StopEvent {
		t.Error("Expected STOP event")
	}
	event.Destroy()
	node1.Stop()
}
