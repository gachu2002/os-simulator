import type { SessionEvent, SnapshotDTO } from "../lib/types";

export interface LogEntry {
  id: number;
  sequence: number;
  type: string;
  tick: number;
  traceHash: string;
  detail: string;
}

export interface SessionState {
  connected: boolean;
  sessionID: string;
  snapshot: SnapshotDTO | null;
  lastSequence: number;
  nextLogID: number;
  error: string;
  logs: LogEntry[];
}

export type SessionAction =
  | { type: "session.created"; sessionID: string; snapshot: SnapshotDTO }
  | { type: "socket.connected" }
  | { type: "socket.disconnected" }
  | { type: "event.received"; event: SessionEvent }
  | { type: "error"; message: string };

export const initialSessionState: SessionState = {
  connected: false,
  sessionID: "",
  snapshot: null,
  lastSequence: 0,
  nextLogID: 1,
  error: "",
  logs: [],
};

export function sessionReducer(
  state: SessionState,
  action: SessionAction,
): SessionState {
  switch (action.type) {
    case "session.created":
      return {
        ...state,
        sessionID: action.sessionID,
        snapshot: action.snapshot,
        nextLogID: state.nextLogID + 1,
        logs: appendLog(state.logs, makeLogEntry(state.nextLogID, {
          sequence: 1,
          type: "session.created",
          tick: action.snapshot.tick,
          traceHash: action.snapshot.trace_hash,
          detail: `session=${action.sessionID}`,
        })),
      };
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
          nextLogID: state.nextLogID + 1,
          logs: appendLog(state.logs, makeLogEntry(state.nextLogID, {
            sequence: event.sequence,
            type: event.type,
            tick: state.snapshot?.tick ?? 0,
            traceHash: state.snapshot?.trace_hash ?? "",
            detail: event.error ?? "unknown session error",
          })),
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
        nextLogID: state.nextLogID + 1,
        logs: appendLog(state.logs, makeLogEntry(state.nextLogID, {
          sequence: event.sequence,
          type: event.type,
          tick: event.snapshot.tick,
          traceHash: event.snapshot.trace_hash,
          detail: event.snapshot.last_command ?? "snapshot",
        })),
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

function appendLog(logs: LogEntry[], entry: LogEntry): LogEntry[] {
  const next = [...logs, entry];
  if (next.length <= 200) {
    return next;
  }
  return next.slice(next.length - 200);
}

function makeLogEntry(
  id: number,
  entry: Omit<LogEntry, "id">,
): LogEntry {
  return { id, ...entry };
}
