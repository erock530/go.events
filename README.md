# go.events Package

## Overview

The `go.events` package provides a simple and flexible event emitter implementation in Go. It allows you to define and emit custom events, register listeners, and manage event-driven interactions within your Go applications.

### Features

- Register multiple listeners for various events.
- Emit events with arbitrary payloads.
- Add listeners that fire once and then remove themselves.
- Manage and remove listeners dynamically.
- Control the maximum number of listeners to avoid potential memory leaks.

### Usage Examples

#### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/erock530/go.events"
)

func main() {
    // Create a new event emitter
    e := events.New()

    // Register listeners
    e.On("user_created", func(payload ...interface{}) {
        fmt.Println("A new User just created!")
    })

    e.On("user_joined", func(payload ...interface{}) {
        user := payload[0].(string)
        room := payload[1].(string)
        fmt.Printf("%s joined the room: %s\n", user, room)
    })

    // Emit events
    e.Emit("user_created")
    e.Emit("user_joined", "user1", "room1")
}
```
#### Using the Default Event Emitter

```go
package main

import (
    "fmt"
    "github.com/erock530/go.events"
)

func main() {
    // Register listeners to the default event emitter
    events.On("user_left", func(payload ...interface{}) {
        user := payload[0].(string)
        room := payload[1].(string)
        fmt.Printf("%s left the room: %s\n", user, room)
    })

    // Emit events using the default event emitter
    events.Emit("user_left", "user1", "room1")
}
```

#### Removing Listeners

```go
package main

import (
    "fmt"
    "github.com/erock530/go.events"
)

func main() {
    // Create a new event emitter
    e := events.New()

    var count = 0
    listener := func(payload ...interface{}) {
        fmt.Println("Event triggered")
        count++
    }

    e.On("my_event", listener)

    e.Emit("my_event")
    e.RemoveListener("my_event", listener)
    e.Emit("my_event")

    fmt.Printf("Listener was called %d times\n", count) // Should print 1
}
```

## Installation

```
go get github.com/erock530/go.events
```

## Versioning

Current version is 0.0.1

Current Go version: 1.22

Read more about Semantic Versioning @ [SemVer.org](http://semver.org/)