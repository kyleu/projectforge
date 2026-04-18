let appUnloading = false;
let unloadHandlerAttached = false;
const maxPendingMessages = 100;

export interface SocketMessage {
  readonly channel: string;
  readonly cmd: string;
  readonly param: unknown;
}

export type Message = SocketMessage;

function socketUrl(u?: string) {
  u ??= "/connect";
  if (u.startsWith("ws")) {
    return u;
  }
  const l = document.location;
  let protocol = "ws";
  if (l.protocol === "https:") {
    protocol = "wss";
  }
  if (!u.startsWith("/")) {
    u = "/" + u;
  }
  return protocol + `://${l.host}${u}`;
}

export class Socket {
  readonly debug: boolean;
  private readonly open: () => void;
  private readonly recv: (m: SocketMessage) => void;
  private readonly err: (svc: string, err: string) => void;
  readonly url?: string;
  connected: boolean;
  pauseSeconds: number;
  pendingMessages: SocketMessage[];
  connectTime?: number;
  closed?: boolean;
  sock?: WebSocket;

  constructor(
    debug: boolean,
    o: () => void,
    r: (m: SocketMessage) => void,
    e: (svc: string, err: string) => void,
    url?: string
  ) {
    this.debug = debug;
    this.open = o;
    this.recv = r;
    this.err = e;
    this.url = socketUrl(url);
    this.connected = false;
    this.pauseSeconds = 1;
    this.pendingMessages = [];

    this.connect();
  }

  connect() {
    if (!unloadHandlerAttached) {
      window.addEventListener("beforeunload", () => {
        appUnloading = true;
      });
      unloadHandlerAttached = true;
    }
    this.connectTime = Date.now();
    this.sock = new WebSocket(socketUrl(this.url));
    this.sock.onopen = () => {
      this.connected = true;
      this.pendingMessages.forEach((msg) => {
        this.send(msg);
      });
      this.pendingMessages = [];
      if (this.debug) {
        console.log("WebSocket connected");
      }
      this.open();
    };
    this.sock.onmessage = (event) => {
      let msg: SocketMessage;
      try {
        msg = JSON.parse(String(event.data)) as SocketMessage;
      } catch (err) {
        this.err("socket", `invalid message payload: ${String(err)}`);
        return;
      }
      if (this.debug) {
        console.debug("[socket]: receive", msg);
      }
      if (msg.cmd === "close-connection") {
        this.disconnect();
      } else {
        this.recv(msg);
      }
    };
    this.sock.onerror = (event) => {
      this.err("socket", event.type);
    };
    this.sock.onclose = () => {
      if (appUnloading || this.closed) {
        return;
      }
      this.connected = false;
      const elapsed = this.connectTime ? Date.now() - this.connectTime : 0;
      if (elapsed > 0 && elapsed < 2000) {
        this.pauseSeconds = this.pauseSeconds * 2;
        if (this.debug) {
          console.debug(`socket closed immediately, reconnecting in ${this.pauseSeconds.toString()} seconds`);
        }
        setTimeout(() => {
          this.connect();
        }, this.pauseSeconds * 1000);
      } else {
        console.debug("socket closed after [" + elapsed.toString() + "ms]");
        setTimeout(() => {
          this.connect();
        }, this.pauseSeconds * 500);
      }
    };
  }

  disconnect() {
    this.closed = true;
    setTimeout(() => {
      this.sock?.close();
      console.debug("[socket] closed");
    }, 500);
  }

  send(msg: SocketMessage) {
    if (this.debug) {
      console.debug("out", msg);
    }
    if (!this.sock) {
      throw new Error("socket not initialized");
    }
    if (this.connected) {
      const m = JSON.stringify(msg, null, 2);
      this.sock.send(m);
    } else {
      this.pendingMessages.push(msg);
      if (this.pendingMessages.length > maxPendingMessages) {
        this.pendingMessages.shift();
      }
    }
  }
}

export function socketInit() {
  return Socket;
}
