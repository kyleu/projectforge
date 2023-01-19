// Content managed by Project Forge, see [projectforge.md] for details.
export function socketInit() {
  (window as any).projectforge.Socket = Socket;
}

let appUnloading = false;

export interface Message {
  readonly channel: string;
  readonly cmd: string;
  readonly param: any;
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
  connectTime?: number
  sock?: WebSocket;

  constructor(debug: boolean, open: () => void, recv: (m: Message) => void, err: (svc: string, err: string) => void, url?: string) {
    this.debug = debug
    this.open = open;
    this.recv = recv;
    this.err = err;
    this.url = socketUrl(url);
    this.connected = false;
    this.pauseSeconds = 1;
    this.pendingMessages = [];

    this.connect();
  }

  connect() {
    window.onbeforeunload = function(){
      appUnloading = true;
    };
    this.connectTime = Date.now();
    this.sock = new WebSocket(socketUrl(this.url));
    const s = this;
    this.sock.onopen = () => {
      s.connected = true;
      s.pendingMessages.forEach(s.send);
      s.pendingMessages = [];
      if (s.debug) {
        console.log("WebSocket connected")
      }
      s.open();
    };
    this.sock.onmessage = (event) => {
      const msg = JSON.parse(event.data) as Message;
      if (s.debug) {
        console.debug("[socket]: receive", msg);
      }
      s.recv(msg);
    };
    this.sock.onerror = (event) => () => {
      s.err("socket", event.type);
    }
    this.sock.onclose = () => {
      if (appUnloading) {
        return;
      }
      s.connected = false;
      const elapsed = s.connectTime ? Date.now() - s.connectTime : 0;
      if (0 < elapsed && elapsed < 2000) {
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

  }

  send(msg: Message) {
    if (this.debug) {
      console.debug("out", msg);
    }
    if (!this.sock) {
      throw "not initialized";
    }
    if (this.connected) {
      const m = JSON.stringify(msg, null, 2);
      this.sock.send(m);
    } else {
      this.pendingMessages.push(msg);
    }
  }
}

function socketUrl(u?: string) {
  if (!u) {
    u = "/connect";
  }
  if (u.indexOf("ws") == 0) {
    return u
  }
  const l = document.location;
  let protocol = "ws";
  if (l.protocol === "https:") {
    protocol = "wss";
  }
  if (u.indexOf("/") != 0) {
    u = "/" + u;
  }
  return protocol + `://${l.host}${u}`;
}
