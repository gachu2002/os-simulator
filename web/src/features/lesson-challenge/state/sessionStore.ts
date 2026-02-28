import { create } from "zustand";

import type { SessionEvent, SnapshotDTO } from "../../../lib/types";

interface SessionStore {
  connected: boolean;
  snapshot: SnapshotDTO | null;
  lastSequence: number;
  error: string;
  reset: () => void;
  onSocketConnected: () => void;
  onSocketDisconnected: () => void;
  onEvent: (event: SessionEvent) => void;
  setError: (message: string) => void;
}

export const useSessionStore = create<SessionStore>((set, get) => ({
  connected: false,
  snapshot: null,
  lastSequence: 0,
  error: "",
  reset: () => set({ connected: false, snapshot: null, lastSequence: 0, error: "" }),
  onSocketConnected: () => set({ connected: true, error: "" }),
  onSocketDisconnected: () => set({ connected: false }),
  onEvent: (event) => {
    const current = get();
    if (event.sequence <= current.lastSequence) {
      return;
    }
    if (event.type === "session.error") {
      set({ lastSequence: event.sequence, error: event.error ?? "unknown session error" });
      return;
    }
    if (!event.snapshot) {
      return;
    }
    set({ lastSequence: event.sequence, error: "", snapshot: event.snapshot });
  },
  setError: (message) => set({ error: message }),
}));
