import type { Command, CommandEnvelope, SessionEvent } from "./types";

export interface SessionSocket {
  sendCommand: (command: Command) => void;
  close: () => void;
}

export function connectSessionSocket(
  baseUrl: string,
  sessionID: string,
  onEvent: (event: SessionEvent) => void,
  onError: (error: Error) => void,
): SessionSocket {
  const wsUrl = toWebSocketURL(baseUrl, sessionID);
  const socket = new WebSocket(wsUrl);

  socket.onmessage = (message) => {
    try {
      const payload = JSON.parse(message.data as string) as SessionEvent;
      onEvent(payload);
    } catch (err) {
      onError(err instanceof Error ? err : new Error("invalid event payload"));
    }
  };

  socket.onerror = () => {
    onError(new Error("websocket transport error"));
  };

  return {
    sendCommand(command: Command) {
      const envelope: CommandEnvelope = {
        type: "command",
        command,
      };
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(envelope));
      }
    },
    close() {
      socket.close();
    },
  };
}

function toWebSocketURL(baseUrl: string, sessionID: string): string {
  const parsed = new URL(baseUrl);
  if (parsed.protocol === "https:") {
    parsed.protocol = "wss:";
  } else if (parsed.protocol === "http:") {
    parsed.protocol = "ws:";
  } else {
    throw new Error("base URL must start with http:// or https://");
  }
  const basePath = parsed.pathname.replace(/\/$/, "");
  parsed.pathname = `${basePath}/ws/${sessionID}`;
  parsed.search = "";
  parsed.hash = "";
  return parsed.toString();
}
