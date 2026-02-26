import type { SessionEvent, SnapshotDTO } from "../lib/types";

export interface SessionState {
  connected: boolean;
  snapshot: SnapshotDTO | null;
  lastSequence: number;
  error: string;
}

export type SessionAction =
  | { type: "session.reset" }
  | { type: "socket.connected" }
  | { type: "socket.disconnected" }
  | { type: "event.received"; event: SessionEvent }
  | { type: "error"; message: string };

export const initialSessionState: SessionState = {
  connected: false,
  snapshot: null,
  lastSequence: 0,
  error: "",
};

export function sessionReducer(
  state: SessionState,
  action: SessionAction,
): SessionState {
  switch (action.type) {
    case "session.reset":
      return initialSessionState;
    case "socket.connected":
      return {
        ...state,
        connected: true,
        error: "",
      };
    case "socket.disconnected":
      return {
        ...state,
        connected: false,
      };
    case "event.received": {
      const event = action.event;
      if (event.sequence <= state.lastSequence) {
        return state;
      }
      if (event.type === "session.error") {
        return {
          ...state,
          lastSequence: event.sequence,
          error: event.error ?? "unknown session error",
        };
      }
      if (!event.snapshot) {
        return state;
      }
      return {
        ...state,
        lastSequence: event.sequence,
        error: "",
        snapshot: event.snapshot,
      };
    }
    case "error":
      return {
        ...state,
        error: action.message,
      };
    default:
      return state;
  }
}
