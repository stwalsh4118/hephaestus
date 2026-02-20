import type { StatusMessage } from "@/types/deploy";

const WS_BASE =
  process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080";
const WS_STATUS_PATH = "/ws/status";
const RECONNECT_BASE_MS = 1000;
const RECONNECT_MAX_MS = 30000;

let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let reconnectAttempts = 0;

function reconnectDelay(): number {
  const delay = Math.min(
    RECONNECT_BASE_MS * Math.pow(2, reconnectAttempts),
    RECONNECT_MAX_MS
  );
  return delay;
}

/**
 * Connect to the WebSocket status endpoint.
 * Calls onMessage for each received StatusMessage.
 * Automatically reconnects with exponential backoff on disconnection.
 */
export function connectStatusWs(
  onMessage: (msg: StatusMessage) => void
): WebSocket {
  if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
    return ws;
  }

  const socket = new WebSocket(`${WS_BASE}${WS_STATUS_PATH}`);

  socket.onopen = () => {
    reconnectAttempts = 0;
  };

  socket.onmessage = (event: MessageEvent) => {
    try {
      const msg = JSON.parse(event.data as string) as StatusMessage;
      if (msg.type === "status_update") {
        onMessage(msg);
      }
    } catch {
      // Ignore non-JSON or unknown messages.
    }
  };

  socket.onclose = () => {
    ws = null;
    const delay = reconnectDelay();
    reconnectAttempts += 1;
    reconnectTimer = setTimeout(() => {
      connectStatusWs(onMessage);
    }, delay);
  };

  socket.onerror = () => {
    // The close event will fire after error, triggering reconnect.
  };

  ws = socket;
  return socket;
}

/**
 * Disconnect the WebSocket and stop reconnection attempts.
 */
export function disconnectStatusWs(): void {
  if (reconnectTimer !== null) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  reconnectAttempts = 0;

  if (ws) {
    ws.onclose = null; // Prevent reconnection on intentional close.
    ws.close();
    ws = null;
  }
}
