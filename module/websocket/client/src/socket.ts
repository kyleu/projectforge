let appUnloading = false;
let unloadHandlerAttached = false;
const maxPendingMessages = 100;

export interface SocketMessage {
  readonly channel: string;
  readonly cmd: string;
  readonly param: Record<string, unknown>;
}

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
    const s = this; // eslint-disable-line @typescript-eslint/no-this-alias
    this.sock.onopen = () => {
      s.connected = true;
      s.pendingMessages.forEach((msg) => {
        s.send(msg);
      });
      s.pendingMessages = [];
      if (s.debug) {
        console.log("WebSocket connected");
      }
      s.open();
    };
    this.sock.onmessage = (event) => {
      let msg: SocketMessage;
      try {
        msg = JSON.parse(String(event.data)) as SocketMessage;
      } catch (err) {
        s.err("socket", `invalid message payload: ${String(err)}`);
        return;
      }
      if (s.debug) {
        console.debug("[socket]: receive", msg);
      }
      if (msg.cmd === "close-connection") {
        s.disconnect();
      } else {
        s.recv(msg);
      }
    };
    this.sock.onerror = (event) => {
      s.err("socket", event.type);
    };
    this.sock.onclose = () => {
      if (appUnloading || s.closed) {
        return;
      }
      s.connected = false;
      const elapsed = s.connectTime ? Date.now() - s.connectTime : 0;
      if (elapsed > 0 && elapsed < 2000) {
        s.pauseSeconds = s.pauseSeconds * 2;
        if (s.debug) {
          console.debug(`socket closed immediately, reconnecting in ${s.pauseSeconds.toString()} seconds`);
        }
        setTimeout(() => {
          s.connect();
        }, s.pauseSeconds * 1000);
      } else {
        console.debug("socket closed after [" + elapsed.toString() + "ms]");
        setTimeout(() => {
          s.connect();
        }, s.pauseSeconds * 500);
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
