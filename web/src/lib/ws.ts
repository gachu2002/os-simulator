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
  if (baseUrl.startsWith("https://")) {
    return `${baseUrl.replace("https://", "wss://")}/ws/${sessionID}`;
  }
  if (baseUrl.startsWith("http://")) {
    return `${baseUrl.replace("http://", "ws://")}/ws/${sessionID}`;
  }
  throw new Error("base URL must start with http:// or https://");
}
