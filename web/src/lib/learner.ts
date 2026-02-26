const LEARNER_ID_KEY = "os-sim.learner-id";

export function getOrCreateLearnerID(): string {
  if (typeof window === "undefined") {
    return "anonymous";
  }

  const stored = window.localStorage.getItem(LEARNER_ID_KEY);
  if (stored && stored.trim() !== "") {
    return stored;
  }

  const generated =
    typeof crypto !== "undefined" && typeof crypto.randomUUID === "function"
      ? crypto.randomUUID()
      : `${Date.now()}-${Math.random().toString(16).slice(2)}`;
  const id = `learner-${generated}`;
  window.localStorage.setItem(LEARNER_ID_KEY, id);
  return id;
}
