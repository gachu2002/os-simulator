import { create } from "zustand";

import type { ChallengeGrade, ChallengeStart } from "../../../entities/challenge/model";

interface ChallengeStore {
  attempt: ChallengeStart | null;
  runResult: ChallengeGrade | null;
  runError: string;
  setAttempt: (attempt: ChallengeStart | null) => void;
  setRunResult: (result: ChallengeGrade | null) => void;
  setRunError: (message: string) => void;
  clearRunError: () => void;
}

export const useChallengeStore = create<ChallengeStore>((set) => ({
  attempt: null,
  runResult: null,
  runError: "",
  setAttempt: (attempt) => set({ attempt }),
  setRunResult: (runResult) => set({ runResult }),
  setRunError: (runError) => set({ runError }),
  clearRunError: () => set({ runError: "" }),
}));
