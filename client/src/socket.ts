let appUnloading = false;

export type Message = {
  readonly channel: string;
  readonly cmd: string;
  readonly param: { [key: string]: unknown };
}

function socketUrl(u?: string) {
  if (!u) {
    u = "/connect";
  }
  if (u.indexOf("ws") === 0) {
    return u;
  }
  const l = document.location;
  let protocol = "ws";
  if (l.protocol === "https:") {
    protocol = "wss";
  }
  if (u.indexOf("/") !== 0) {
    u = "/" + u;
  }
  return protocol + `://${l.host}${u}`;
}

export class Socket {
  readonly debug: boolean;
  private readonly open: () => void;
  private readonly recv: (m: Message) => void;
  private readonly err: (svc: string, err: string) => void;
  readonly url?: string;
  connected: boolean;
  pauseSeconds: number;
  pendingMessages: Message[];
  connectTime?: number;
  closed?: boolean;
  sock?: WebSocket;

  constructor(debug: boolean, o: () => void, r: (m: Message) => void, e: (svc: string, err: string) => void, url?: string) {
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
    window.onbeforeunload = () => {
      appUnloading = true;
    };
    this.connectTime = Date.now();
    this.sock = new WebSocket(socketUrl(this.url));
    const s = this; // eslint-disable-line @typescript-eslint/no-this-alias
    this.sock.onopen = () => {
      s.connected = true;
      s.pendingMessages.forEach(s.send);
      s.pendingMessages = [];
      if (s.debug) {
        console.log("WebSocket connected");
      }
      s.open();
    };
    this.sock.onmessage = (event) => {
      const msg = JSON.parse(event.data) as Message;
      if (s.debug) {
        console.debug("[socket]: receive", msg);
      }
      if (msg.cmd === "close-connection") {
        s.disconnect();
      } else {
        s.recv(msg);
      }
    };
    this.sock.onerror = (event) => () => {
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
          console.debug(`socket closed immediately, reconnecting in ${s.pauseSeconds} seconds`);
        }
        setTimeout(() => {
          s.connect();
        }, s.pauseSeconds * 1000);
      } else {
        console.debug("socket closed after [" + elapsed + "ms]");
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

  send(msg: Message) {
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
    }
  }
}

export function socketInit() {
  return Socket;
}
