// Package gozyre provides ...
package gozyre

/*
#include "czmq.h"
*/
import "C"

import (
	"unsafe"
)

func zListToSlice(zlist *C.struct__zlist_t, free bool) []string {
	if free {
		defer C.zlist_destroy(&zlist)
	}
	members := make([]string, 0)
	for item := C.zlist_first(zlist); item != nil; item = C.zlist_next(zlist) {
		members = append(members, C.GoString((*C.char)(item)))
	}
	return members
}

func zHashToMap(zhash *C.struct__zhash_t, free bool) map[string]string {
	if free {
		defer C.zhash_destroy(&zhash)
	}

	keyVal := make(map[string]string)
	for item := C.zhash_first(zhash); item != nil; item = C.zhash_next(zhash) {
		ckey := C.zhash_cursor(zhash)
		cval := C.GoString((*C.char)(C.zhash_lookup(zhash, ckey)))
		keyVal[C.GoString(ckey)] = cval
	}
	return keyVal
}

func zMsgToBytes(zmsg *C.struct__zmsg_t) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		defer close(ch)
		if zmsg == nil {
			return
		}
		for frame := C.zmsg_first(zmsg); frame != nil; frame = C.zmsg_next(zmsg) {
			data := C.GoBytes(unsafe.Pointer(C.zframe_data(frame)), C.int(C.zframe_size(frame)))
			ch <- data
		}
	}()
	return ch
}

func bytesToZmsg(ch <-chan []byte) (*C.struct__zmsg_t, error) {
	zmsg := C.zmsg_new()
	var err error
	for frame := range ch {
		// XXX: Need to check if zmq allocates space for the frame in case the
		// gc collects it too early
		ret := C.zmsg_addmem(zmsg, unsafe.Pointer(&frame[0]), C.size_t(len(frame)))
		if ret != 0 {
			C.zmsg_destroy(unsafe.Pointer(zmsg))
			err = ErrAddingFrame
			break
		}
	}
	// Consume all channel data so the sending goroutine does't keep running
	if err != nil {
		go func() {
			for _ = range ch {
			}
		}()
	}
	return zmsg, err
}
