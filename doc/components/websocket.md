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
        err := as.Services.Socket.Upgrade(w, r, channel, ps.Profile, ps.Logger)
        if err != nil {
            ps.Logger.Error("WebSocket upgrade failed", "error", err)
            return "", err
        }

        // Connection upgraded successfully
        return "", nil
    })
}
```

### Advanced WebSocket Setup

For more complex scenarios with custom logic:

```go
func ChatRoomSocket(w http.ResponseWriter, r *http.Request) {
    controller.Act("chat.room", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
        // Get room ID from URL parameters
        roomID := ps.RequestString("room", true)
        if roomID == "" {
            return "", errors.New("room ID is required")
        }

        // Verify user has access to this room
        hasAccess, err := models.UserHasRoomAccess(as.DB, ps.Profile.ID, roomID)
        if err != nil {
            return "", err
        }
        if !hasAccess {
            return "", errors.New("access denied to chat room")
        }

        // Create channel name for this room
        channel := fmt.Sprintf("chat-room-%s", roomID)

        // Set up connection metadata
        metadata := map[string]interface{}{
            "room_id": roomID,
            "user_id": ps.Profile.ID,
            "username": ps.Profile.Name,
        }

        // Upgrade with metadata
        err = as.Services.Socket.UpgradeWithMetadata(w, r, channel, ps.Profile, metadata, ps.Logger)
        if err != nil {
            return "", err
        }

        // Notify room that user joined
        joinMessage := &SocketMessage{
            Channel: channel,
            Cmd:     "user-joined",
            Param: map[string]interface{}{
                "user_id": ps.Profile.ID,
                "username": ps.Profile.Name,
                "timestamp": time.Now(),
            },
        }

        as.Services.Socket.Broadcast(channel, joinMessage)

        return "", nil
    })
}
```

### Message Broadcasting

Send messages to connected clients:

```go
// Broadcast to all clients in a channel
func BroadcastNotification(as *app.State, channel string, notification *Notification) {
    message := &SocketMessage{
        Channel: channel,
        Cmd:     "notification",
        Param: map[string]interface{}{
            "id":       notification.ID,
            "type":     notification.Type,
            "title":    notification.Title,
            "message":  notification.Message,
            "timestamp": notification.CreatedAt,
        },
    }

    as.Services.Socket.Broadcast(channel, message)
}

// Send message to specific user
func SendUserMessage(as *app.State, userID int, messageType string, data interface{}) {
    channel := fmt.Sprintf("user-%d", userID)

    message := &SocketMessage{
        Channel: channel,
        Cmd:     messageType,
        Param:   data,
    }

    as.Services.Socket.SendToChannel(channel, message)
}

// Send message to multiple channels
func BroadcastSystemUpdate(as *app.State, channels []string, updateData interface{}) {
    message := &SocketMessage{
        Channel: "system",
        Cmd:     "system-update",
        Param:   updateData,
    }

    for _, channel := range channels {
        as.Services.Socket.SendToChannel(channel, message)
    }
}
```

## Client-Side Implementation

Set up a WebSocket connection on the client side:

```typescript
import {Socket, Message} from "./socket";

// Define callback functions
function onSocketOpen() {
    console.log("WebSocket connection established");
    // Update UI to show connected state
    updateConnectionStatus(true);
}

function onMessageReceived(message: Message) {
    console.log("Message received:", message);

    // Handle different message types
    switch (message.cmd) {
        case "notification":
            showNotification(message.param);
            break;
        case "user-joined":
            handleUserJoined(message.param);
            break;
        case "chat-message":
            displayChatMessage(message.param);
            break;
        case "system-update":
            handleSystemUpdate(message.param);
            break;
        default:
            console.warn("Unknown message type:", message.cmd);
    }
}

function onSocketError(service: string, error: string) {
    console.error("WebSocket error:", service, error);
    // Update UI to show error state
    updateConnectionStatus(false, error);
}

// Initialize WebSocket connection
document.addEventListener("DOMContentLoaded", function() {
    // Enable debug mode in development
    const debug = window.location.hostname === "localhost";

    // Create socket connection
    const socket = new Socket(debug, onSocketOpen, onMessageReceived, onSocketError, "/ws/notifications");

    // Store socket reference for later use
    (window as any).appSocket = socket;
});
```

## Message Format

```typescript
interface Message {
    readonly channel: string;  // Channel identifier
    readonly cmd: string;      // Command/message type
    readonly param: Record<string, unknown>; // Message payload
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
        err := as.Services.Socket.Upgrade(w, r, channel, ps.Profile, ps.Logger)
        return "", err
    })
}
```
