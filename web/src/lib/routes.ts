export type AppRoute =
  | { kind: "overview" }
  | { kind: "learn"; lessonID: string; stageIndex?: number }
  | { kind: "challenge"; lessonID: string; stageIndex?: number };

export function parseRoute(pathname: string): AppRoute {
  const [pathOnly, search = ""] = pathname.split("?");
  const normalized = pathOnly.length > 1 ? pathOnly.replace(/\/+$/, "") : pathOnly;
  const params = new URLSearchParams(search);
  const stageParam = Number(params.get("stage") ?? 0);
  const stageIndex = Number.isFinite(stageParam) ? stageParam : 0;

  if (normalized === "/") {
    return { kind: "overview" };
  }

  const lessonLearnPrefix = "/lesson/";
  if (normalized.startsWith(lessonLearnPrefix) && normalized.endsWith("/learn")) {
    const lessonID = decodeURIComponent(
      normalized.slice(lessonLearnPrefix.length, normalized.length - "/learn".length),
    ).trim();
    if (lessonID !== "") {
      return { kind: "learn", lessonID, stageIndex };
    }
  }

  if (normalized.startsWith(lessonLearnPrefix) && normalized.endsWith("/challenge")) {
    const lessonID = decodeURIComponent(
      normalized.slice(lessonLearnPrefix.length, normalized.length - "/challenge".length),
    ).trim();
    if (lessonID !== "") {
      return { kind: "challenge", lessonID, stageIndex };
    }
  }

  return { kind: "overview" };
}
