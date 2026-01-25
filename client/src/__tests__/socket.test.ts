/* @vitest-environment happy-dom */
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { Socket } from "../socket";

class MockWebSocket {
  static instances: MockWebSocket[] = [];
  url: string;
  sent: string[] = [];
  closeCalled = false;
  onopen?: () => void;
  onmessage?: (event: { data: string }) => void;
  onerror?: (event: { type: string }) => void;
  onclose?: () => void;

  constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
  }

  send(data: string) {
    this.sent.push(data);
  }

  close() {
    this.closeCalled = true;
    this.onclose?.();
  }

  triggerOpen() {
    this.onopen?.();
  }

  triggerMessage(data: unknown) {
    this.onmessage?.({ data: JSON.stringify(data) });
  }

  triggerRawMessage(data: string) {
    this.onmessage?.({ data });
  }

  triggerError(type = "error") {
    this.onerror?.({ type });
  }

  triggerClose() {
    this.onclose?.();
  }
}

beforeEach(() => {
  MockWebSocket.instances = [];
  globalThis.WebSocket = MockWebSocket as unknown as typeof WebSocket;
});

afterEach(() => {
  vi.useRealTimers();
  vi.clearAllTimers();
  document.body.innerHTML = "";
});

describe("socket", () => {
  it("queues messages until connected", () => {
    const sock = new Socket(
      false,
      () => undefined,
      () => undefined,
      () => undefined,
      "ws://example"
    );
    const ws = MockWebSocket.instances[0];

    sock.send({ channel: "c", cmd: "ping", param: { ok: true } });
    expect(sock.pendingMessages).toHaveLength(1);
    expect(ws.sent).toHaveLength(0);

    ws.triggerOpen();
    expect(sock.pendingMessages).toHaveLength(0);
    expect(ws.sent).toHaveLength(1);

    const payload = JSON.parse(ws.sent[0] ?? "{}") as {
      cmd?: string;
      param?: { ok?: boolean };
    };
    expect(payload.cmd).toBe("ping");
    expect(payload.param?.ok).toBe(true);
  });

  it("handles close-connection and disconnects", () => {
    vi.useFakeTimers();
    const sock = new Socket(
      false,
      () => undefined,
      () => undefined,
      () => undefined,
      "ws://example"
    );
    const ws = MockWebSocket.instances[0];
    ws.triggerOpen();

    ws.triggerMessage({ channel: "c", cmd: "close-connection", param: {} });
    expect(sock.closed).toBe(true);

    vi.advanceTimersByTime(500);
    expect(ws.closeCalled).toBe(true);
  });

  it("reports invalid message payloads", () => {
    const err = vi.fn();
    const sock = new Socket(
      false,
      () => undefined,
      () => undefined,
      err,
      "ws://example"
    );
    const ws = MockWebSocket.instances[0];
    ws.triggerOpen();

    ws.triggerRawMessage("not-json");
    expect(err).toHaveBeenCalledTimes(1);
    expect(err.mock.calls[0]?.[0]).toBe("socket");
    expect(String(err.mock.calls[0]?.[1])).toContain("invalid message payload");
    expect(sock.connected).toBe(true);
  });

  it("reconnects with backoff after quick close", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2024-01-01T00:00:00Z"));

    const sock = new Socket(
      false,
      () => undefined,
      () => undefined,
      () => undefined,
      "ws://example"
    );
    const ws = MockWebSocket.instances[0];

    vi.advanceTimersByTime(1);
    ws.triggerClose();
    expect(sock.pauseSeconds).toBe(2);
    expect(MockWebSocket.instances).toHaveLength(1);

    vi.advanceTimersByTime(1999);
    expect(MockWebSocket.instances).toHaveLength(1);

    vi.advanceTimersByTime(1);
    expect(MockWebSocket.instances).toHaveLength(2);
  });
});
