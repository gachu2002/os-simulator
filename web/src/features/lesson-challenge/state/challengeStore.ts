import { create } from "zustand";

import type { ChallengeGrade, ChallengeStart } from "../../../entities/challenge/model";

type Policy = "fifo" | "rr" | "mlfq";

interface ChallengeStore {
  attempt: ChallengeStart | null;
  runResult: ChallengeGrade | null;
  runError: string;
  policy: Policy;
  quantum: number;
  frames: number;
  tlbEntries: number;
  diskLatency: number;
  terminalLatency: number;
  setAttempt: (attempt: ChallengeStart | null) => void;
  setRunResult: (result: ChallengeGrade | null) => void;
  setRunError: (message: string) => void;
  clearRunError: () => void;
  setPolicy: (value: Policy) => void;
  setQuantum: (value: number) => void;
  setFrames: (value: number) => void;
  setTLBEntries: (value: number) => void;
  setDiskLatency: (value: number) => void;
  setTerminalLatency: (value: number) => void;
}

export const useChallengeStore = create<ChallengeStore>((set) => ({
  attempt: null,
  runResult: null,
  runError: "",
  policy: "rr",
  quantum: 2,
  frames: 8,
  tlbEntries: 4,
  diskLatency: 3,
  terminalLatency: 1,
  setAttempt: (attempt) => set({ attempt }),
  setRunResult: (runResult) => set({ runResult }),
  setRunError: (runError) => set({ runError }),
  clearRunError: () => set({ runError: "" }),
  setPolicy: (policy) => set({ policy }),
  setQuantum: (quantum) => set({ quantum }),
  setFrames: (frames) => set({ frames }),
  setTLBEntries: (tlbEntries) => set({ tlbEntries }),
  setDiskLatency: (diskLatency) => set({ diskLatency }),
  setTerminalLatency: (terminalLatency) => set({ terminalLatency }),
}));
