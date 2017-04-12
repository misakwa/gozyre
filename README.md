# gozyre - Local area clustering for P2P applications

## Introduction
A golang interface to the [Zyre](https://github.com/zeromq/zyre) library.

The interface is not exactly one-to-one, but exposes the basics needed for
writing applications.

## Usage Examples

// create a new node named 'node1' bound to the 'eth1' interface.
// The node will be using udp beaconing on port 5670 and reporting the
// key-value pairs as headers.

```golang
import "github.com/misakwa/gozyre"

func main() {
    node1 := gozyre.New("node1", "eth1", 5670, map[string]string{"key":"value"}, false)
    defer node1.Destroy()
}
```

Create a new node with autogenerated name and default chosen interface

```golang
import "github.com/misakwa/gozyre"

Create a new node with autogenerated name and default chosen interface
func main() {
    node2 := gozyre.New("", "", 0, map[string]string{}, false)
    defer node2.Destroy()
}
```

Much more complete example inspired by the tests

```golang
import (
    "log"
    "time"

    "github.com/misakwa/gozyre"
)

func main() {
    node1 := gozyre.New("node1", "", 0, map[string]string{}, false)
    defer node1.Destroy()

    node2 := gozyre.New("node2", "", 0, map[string]string{"X-HELLO": "World"}, false)
    defer node2.Destroy()

    // Setup gossip discovery with endpoint
    var err error
    err = node1.Gossip("inproc://zyre-node1", "inproc://gossip-hub")
    if err != nil {
      log.Fatalf("Unable to setup gossip discovery: %s", err)
    }

    err = node2.Gossip("inproc://zyre-node2", "inproc://gossip-hub")
    if err != nil {
      log.Fatalf("Unable to setup gossip discovery: %s", err)
    }

    // Start nodes
    node1.Start()
    node2.Start()

    // Join group
    node1.Join("GLOBAL")
    node2.Join("GLOBAL")

    // Allow time for discovery
    time.Sleep(10 * time.Millisecond)

    mch := make(chan []byte)
    go func() {
        defer close(mch)
        mch <- []byte("Hello, World")
    }()
    // node1 shouts message to global
    node1.Shout("GLOBAL", mch)

    // Call recieve as many times as there are messages
    for {
        evt, err := node2.Recv()
        if err != nil {
            log.Fatalf("Unable to receive node message: %s", err)
        }
        log.Println("event type = %s", evt.Type())
        for frame := range evt.Message() {
            log.Printf("%s", frame)
        }
        evt.Destroy()
    }
}
```

## Disclaimer

I haven't used this in production yet and there are most likely bugs everywhere.

## TODO

- Implement chat example in go
- Test all the example snippets
- Add more tests
