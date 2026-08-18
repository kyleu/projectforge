# WebSockets

Project Forge provides a comprehensive WebSocket system that enables real-time, bidirectional communication between the server and client. The WebSocket implementation is designed for reliability, scalability, and ease of use, with automatic reconnection, message queuing, and robust error handling.

## Overview

The WebSocket system offers:
- **Real-time Communication**: Instant bidirectional messaging between server and client
- **Automatic Reconnection**: Intelligent reconnection with exponential backoff
- **Message Queuing**: Automatic queuing of messages when disconnected
- **Channel-based Routing**: Organize communication using logical channels
- **Type-safe Messaging**: Structured message format with TypeScript support
- **Debug Support**: Built-in logging and debugging capabilities

## Server-Side Implementation

### Basic WebSocket Controller

Create a controller action to handle WebSocket upgrades:

```go
package controller

import (
    "net/http"
    "myproject/app"
    "myproject/app/controller/cutil"
)

func MySocketHandler(w http.ResponseWriter, r *http.Request) {
    controller.Act("my.socket", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
        // Define the channel for this connection
        channel := "user-notifications"

        // Upgrade the connection to WebSocket
        id, err := as.Services.Socket.Upgrade(w, r, nil, channel, ps.Profile, ps.Logger)
        if err != nil {
            ps.Logger.Error("WebSocket upgrade failed", "error", err)
            return "", err
        }

        // Connection upgraded successfully
        return "", nil
    })
}
```

## Security Considerations

### Authentication and Authorization

Ensure WebSocket connections are properly authenticated:

```go
func SecureSocketHandler(w http.ResponseWriter, r *http.Request) {
    controller.Act("secure.socket", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
        if ps.Profile == nil {
            return "", errors.New("authentication required")
        }
        channel := fmt.Sprintf("secure-user-%d", ps.Profile.ID)
        id, err := as.Services.Socket.Upgrade(ps.Context, w, r, nil, channel, ps.Profile, ps.Logger)
        return "", err
    })
}
```
