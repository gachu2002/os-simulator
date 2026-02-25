import type { CreateSessionResponse, SessionConfig } from "./types";
import { fetchJSON } from "./http";

export async function createSession(
  baseUrl: string,
  config: SessionConfig,
): Promise<CreateSessionResponse> {
  return fetchJSON<CreateSessionResponse>(baseUrl, "/sessions", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(config),
  });
}
