// Package gozyre provides a wrapping around the zyre c-library. It provides
// clustering for peer-to-peer applications.
// See: https://github.com/zeromq/zyre.
package gozyre

/*
#cgo !windows pkg-config: libczmq libzmq libzyre
#cgo !windows CFLAGS: -DLIBCZMQ_EXPORTS -DZMQ_DEFINED_STDINT -DLIBCZMQ_EXPORTS -DZYRE_EXPORTS
#cgo windows LDFLAGS: -lws2_32 -liphlpapi -lrpcrt4 -lzmq -lczmq -lzyre
#cgo windows CFLAGS: -Wno-pedantic-ms-format -DLIBCZMQ_EXPORTS -DZMQ_DEFINED_STDINT -DLIBCZMQ_EXPORTS -DZYRE_EXPORTS

#include "czmq.h"
#include "zyre.h"
*/
import "C"

import (
	"errors"
)

var (
	//ZMQVersionMajor is the major version of the underlying ZeroMQ library
	ZMQVersionMajor = int(C.ZMQ_VERSION_MAJOR)

	//ZMQVersionMinor is the minor version of the underlying ZeroMQ library
	ZMQVersionMinor = int(C.ZMQ_VERSION_MINOR)

	//CZMQVersionMajor is the major version of the underlying CZMQ library
	CZMQVersionMajor = int(C.CZMQ_VERSION_MAJOR)

	// CZMQVersionMinor is the minor version of the underlying CZMQ library
	CZMQVersionMinor = int(C.CZMQ_VERSION_MINOR)

	// ZyreVersion is the underlying version of the Zyre library
	ZyreVersion = uint64(C.ZYRE_VERSION) // major * 10000 + minor * 100 + patch, as a single integer

	// ZyreVersionMajor is the major version of the underlying Zyre library
	ZyreVersionMajor = int(C.ZYRE_VERSION_MAJOR)

	// ZyreVersionMinor is the minor version of the underlying Zyre library
	ZyreVersionMinor = int(C.ZYRE_VERSION_MINOR)

	// ZyreVersionPatch is the patch version of the underlying Zyre library
	ZyreVersionPatch = int(C.ZYRE_VERSION_PATCH)
)

var (
	// ErrStartFailed is to used to represent startup failures
	ErrStartFailed = errors.New("error starting zyre node")

	// ErrNodeInterrupted is used to indicate an interruped node when receiving messages
	ErrNodeInterrupted = errors.New("error receiving message from node")

	// ErrInvalidEndpoint returned when setting an invalid endpoint
	ErrInvalidEndpoint = errors.New("error setting endpoint")

	// ErrAddingFrame is returned when constructing message frames for sending
	ErrAddingFrame = errors.New("error adding message frame")
)
