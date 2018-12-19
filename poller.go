// Package gozyre provides ...
package gozyre

/*
#include "czmq.h"
#include "czmq_library.h"

zpoller_t *wrap_zpoller_new() {
	return  zpoller_new(NULL);
}
*/
import "C"

import (
	"unsafe"
)

type ZyrePoller struct {
	czpoller *C.struct__zpoller_t
	czsocket *C.struct__zsock_t
}

// New Zyre-socket POLLER.
func NewPoller(node *Zyre) *ZyrePoller {
	sk := node.Socket()
	if (sk == nil) {
		return nil
	}

	zp := C.wrap_zpoller_new()
	if (zp == nil) {
		return nil
	}
	if (C.zpoller_add(zp, unsafe.Pointer(sk)) < 0) {
		C.zpoller_destroy(&zp)
		return nil
	}
	
	pol := &ZyrePoller{}
	pol.czpoller = zp
	pol.czsocket = sk

	return pol
}

// Zyre-Poller destruction
func (poller *ZyrePoller) Destroy() {
	if (poller == nil) {
		return
	}
	if (poller.czpoller == nil) {
		return
	}
	if (poller.czsocket != nil) {
		C.zpoller_remove(poller.czpoller, unsafe.Pointer(poller.czsocket))
		poller.czsocket = nil
	}
	C.zpoller_destroy(&poller.czpoller)
}

// Poll ZYRE socket, until timeout or message is received.
// <timeout> is given in MILISECOND.
// -  0 : no wait.
// - -1 : wait indefinitely.
// 
// Note:
// Polling is not stopped by any SIGNAL (thx to GO signal handler).
// Hence, timeout has to be rather small if response time must be short.
// Good compromise should be around 200 (ms).
func (poller *ZyrePoller) Poll(timeout int) (bool, error) {
	socket := C.zpoller_wait(poller.czpoller, C.int(timeout))
	if (socket == nil) {
		if C.zpoller_terminated(poller.czpoller) {
			// We've been interrupted by SIGINT.
			return false, ErrPollInterrupted
		}
		if C.zpoller_expired(poller.czpoller) {
			return false, ErrPollExpired
		}

		// Wall... not EXPIRED & not INTERRUPTED...
		// What is this strange state ?
		return false, ErrPollUnknown
	}

	return true, nil
}

