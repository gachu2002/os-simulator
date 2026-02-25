import { useMutation } from "@tanstack/react-query";
import { useCallback, useEffect, useMemo, useReducer, useRef, useState } from "react";

import { createSession } from "../lib/api";
import type { Command } from "../lib/types";
import { connectSessionSocket, type SessionSocket } from "../lib/ws";
import { initialSessionState, sessionReducer } from "../state/sessionReducer";

export function useSession() {
  const [state, dispatch] = useReducer(sessionReducer, initialSessionState);
  const [baseURL, setBaseURL] = useState(defaultBaseURL());
  const [seed, setSeed] = useState(1);
  const [policy, setPolicy] = useState<"fifo" | "rr" | "mlfq">("rr");
  const [quantum, setQuantum] = useState(2);
  const socketRef = useRef<SessionSocket | null>(null);

  const createSessionMutation = useMutation({
    mutationFn: () =>
      createSession(baseURL, {
        seed,
        policy,
        quantum,
      }),
  });

  const canSend = useMemo(
    () => Boolean(state.sessionID && state.connected),
    [state.connected, state.sessionID],
  );

  useEffect(() => {
    return () => {
      socketRef.current?.close();
      socketRef.current = null;
    };
  }, []);

  const handleCreateSession = useCallback(async () => {
    dispatch({ type: "error", message: "" });
    try {
      socketRef.current?.close();
      socketRef.current = null;

      const created = await createSessionMutation.mutateAsync();

      dispatch({
        type: "session.created",
        sessionID: created.session_id,
        snapshot: created.snapshot,
      });

      const socket = connectSessionSocket(
        baseURL,
        created.session_id,
        (event) => {
          dispatch({ type: "event.received", event });
          dispatch({ type: "socket.connected" });
        },
        (error) => {
          dispatch({ type: "socket.disconnected" });
          dispatch({ type: "error", message: error.message });
        },
      );

      socketRef.current = socket;
    } catch (error) {
      dispatch({
        type: "error",
        message:
          error instanceof Error ? error.message : "failed to create session",
      });
    }
  }, [baseURL, createSessionMutation]);

  const handleCommand = useCallback(
    (command: Command) => {
      if (!canSend) {
        return;
      }
      socketRef.current?.sendCommand(command);
    },
    [canSend],
  );

  return {
    state,
    baseURL,
    seed,
    policy,
    quantum,
    canSend,
    isCreatingSession: createSessionMutation.isPending,
    setBaseURL,
    setSeed,
    setPolicy,
    setQuantum,
    handleCreateSession,
    handleCommand,
  };
}

function defaultBaseURL(): string {
  const envURL = import.meta.env.VITE_API_BASE_URL;
  if (typeof envURL === "string" && envURL.trim() !== "") {
    return envURL.trim();
  }
  if (typeof window === "undefined") {
    return "http://localhost:8080";
  }
  const host = window.location.hostname || "localhost";
  return `http://${host}:8080`;
}
