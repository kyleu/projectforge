# WebSocket

This module provides a comprehensive WebSocket service for [Project Forge](https://projectforge.dev) applications, enabling real-time bidirectional communication between clients and the server.

## Features

- **Connection Management**: Automatic connection registration, heartbeat, and cleanup
- **Channel System**: Multi-channel communication with join/leave functionality
- **Message Routing**: Type-safe message handling with custom commands
- **User Integration**: Seamless integration with user authentication and profiles
- **Admin Interface**: Built-in WebSocket status and debugging tools
- **TypeScript Client**: Full-featured client with automatic reconnection
- **Broadcasting**: Send messages to channels, users, or all connections

## Architecture

### Core Components

- **Service**: Central WebSocket service managing connections and channels
- **Connection**: Individual WebSocket connection with user context
- **Channel**: Named communication channels for grouping connections
- **Message**: Structured message format with commands and parameters
- **Handler**: Custom message processing logic

## Usage

### 1. Backend Setup

#### Create Controller Handlers

```go
// Page handler - serves the WebSocket interface
func ExamplePage(w http.ResponseWriter, r *http.Request) {
    controller.Act("example", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
        ps.SetTitleAndData("WebSocket Example", nil)
        return controller.Render(r, as, &views.ExamplePage{}, ps)
    })
}

// WebSocket upgrade handler
func ExampleSocket(w http.ResponseWriter, r *http.Request) {
    controller.Act("example.socket", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
        // Get or generate channel
        channel := r.URL.Query().Get("ch")
        if channel == "" {
            channel = "example-" + util.RandomString(8)
        }

        // Create custom message handler
        handler := &ExampleHandler{}

        // Upgrade connection
        connID, err := as.Services.Socket.Upgrade(
            ps.Context, ps.W, ps.R, channel,
            ps.Profile, handler, ps.Logger,
        )
        if err != nil {
            return "", err
        }

        // Start read loop (blocks until connection closes)
        return "", as.Services.Socket.ReadLoop(ps.Context, connID, ps.Logger)
    })
}
```

#### Implement Message Handler

```go
type ExampleHandler struct{}

func (h *ExampleHandler) On(s *websocket.Service, c *websocket.Connection, cmd string, param []byte, logger util.Logger) error {
    switch cmd {
    case "ping":
        // Echo back a pong
        return s.WriteChannel("pong", util.ValueMap{"timestamp": util.TimeCurrentMillis()}, c.Channel, logger)

    case "chat":
        var msg util.ValueMap
        if err := util.FromJSON(param, &msg); err != nil {
            return err
        }
        // Broadcast to all users in channel
        return s.WriteChannel("chat", util.ValueMap{
            "user": c.Profile.Name,
            "text": msg["text"],
            "time": util.TimeCurrentMillis(),
        }, c.Channel, logger)

    case "join-room":
        var data util.ValueMap
        if err := util.FromJSON(param, &data); err != nil {
            return err
        }
        room := fmt.Sprint(data["room"])
        _, err := s.Join(c.ID, room, logger)
        return err

    default:
        logger.Warnf("unhandled websocket command [%s]", cmd)
        return nil
    }
}
```

#### Register Routes

```go
// In your routes setup
makeRoute(r, http.MethodGet, "/example", controller.ExamplePage)
makeRoute(r, http.MethodGet, "/example/socket", controller.ExampleSocket)
```

### 2. Frontend Setup

#### TypeScript Client

```typescript
import { Socket, Message } from "./socket";

class ExampleClient {
    private socket: Socket;
    private channel: string;

    constructor(channel: string = "example") {
        this.channel = channel;

        this.socket = new Socket(
            true,                    // debug mode
            this.onOpen.bind(this),  // connection opened
            this.onMessage.bind(this), // message received
            this.onError.bind(this), // error occurred
            `/example/socket?ch=${channel}` // WebSocket URL
        );
    }

    private onOpen(): void {
        console.log("Connected to WebSocket");
        this.sendPing();
    }

    private onMessage(msg: Message): void {
        switch (msg.cmd) {
            case "pong":
                console.log("Received pong:", msg.param);
                break;

            case "chat":
                this.displayChatMessage(msg.param);
                break;

            case "user-joined":
                this.showUserJoined(msg.param);
                break;

            default:
                console.log("Unknown message:", msg);
        }
    }

    private onError(service: string, error: string): void {
        console.error(`WebSocket error in ${service}:`, error);
    }

    // Public methods
    sendPing(): void {
        this.socket.send({ channel: this.channel, cmd: "ping", param: {} });
    }

    sendChatMessage(text: string): void {
        this.socket.send({
            channel: this.channel,
            cmd: "chat",
            param: { text }
        });
    }

    joinRoom(room: string): void {
        this.socket.send({
            channel: this.channel,
            cmd: "join-room",
            param: { room }
        });
    }
}

// Initialize when page loads
document.addEventListener("DOMContentLoaded", () => {
    const client = new ExampleClient("my-channel");

    // Wire up UI events
    const chatForm = document.getElementById("chat-form") as HTMLFormElement;
    const chatInput = document.getElementById("chat-input") as HTMLInputElement;

    chatForm?.addEventListener("submit", (e) => {
        e.preventDefault();
        const text = chatInput.value.trim();
        if (text) {
            client.sendChatMessage(text);
            chatInput.value = "";
        }
    });
});
```

#### HTML Template Integration

```html
<!-- In your quicktemplate view -->
<div id="websocket-example">
    <div id="messages"></div>
    <form id="chat-form">
        <input type="text" id="chat-input" placeholder="Type a message..." />
        <button type="submit">Send</button>
    </form>
</div>

<script>
    // Simple JavaScript version
    const sock = new YourProjectName.Socket(
        true,
        () => console.log("Connected"),
        (msg) => {
            const messages = document.getElementById("messages");
            const div = document.createElement("div");
            div.textContent = `${msg.cmd}: ${JSON.stringify(msg.param)}`;
            messages.appendChild(div);
        },
        (svc, err) => console.error(`${svc}: ${err}`),
        "/example/socket"
    );

    function sendMessage(cmd, param) {
        sock.send({ channel: "example", cmd, param });
    }
</script>
```

## Service API

### Core Methods

```go
// Upgrade HTTP connection to WebSocket
func (s *Service) Upgrade(ctx context.Context, w http.ResponseWriter, r *http.Request,
    channel string, profile *user.Profile, handler Handler, logger util.Logger) (uuid.UUID, error)

// Join a channel
func (s *Service) Join(connID uuid.UUID, channel string, logger util.Logger) (bool, error)

// Leave a channel
func (s *Service) Leave(connID uuid.UUID, channel string, logger util.Logger) (bool, error)

// Send message to specific connection
func (s *Service) WriteConnection(cmd string, param any, connID uuid.UUID, logger util.Logger) error

// Send message to all connections in channel
func (s *Service) WriteChannel(cmd string, param any, channel string, logger util.Logger) error

// Broadcast to all connections
func (s *Service) WriteAll(cmd string, param any, logger util.Logger) error
```

### Status and Monitoring

```go
// Get service status
func (s *Service) Status() *Status

// Get all connections
func (s *Service) GetConnections() []*Connection

// Get connections by channel
func (s *Service) GetConnectionsByChannel(channel string) []*Connection
```

## Admin Interface

The module includes built-in admin pages accessible at:

- `/admin/websocket` - WebSocket service status and active connections
- Real-time connection monitoring
- Channel membership display
- Message broadcasting tools

## Best Practices

1. **Error Handling**: Always handle WebSocket errors gracefully with reconnection logic
2. **Message Validation**: Validate all incoming message parameters
3. **Rate Limiting**: Implement rate limiting for message-heavy applications
4. **Channel Management**: Use meaningful channel names and clean up unused channels
5. **Security**: Validate user permissions before joining channels or processing commands
6. **Performance**: Avoid blocking operations in message handlers

## Common Patterns

### Chat Application
- Use channels for chat rooms
- Broadcast messages to channel members
- Handle user join/leave events

### Real-time Updates
- Subscribe clients to data channels
- Push updates when server data changes
- Use different message types for different data

### Gaming/Collaboration
- Use per-game/session channels
- Implement custom game logic in handlers
- Synchronize state between clients

### Notifications
- Global notification channel for all users
- User-specific channels for private notifications
- Different message types for different notification levels

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/websocket
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)
