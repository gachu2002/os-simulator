import type { CreateSessionResponse, SessionConfig } from "./types";

export async function createSession(
  baseUrl: string,
  config: SessionConfig,
): Promise<CreateSessionResponse> {
  const response = await fetch(`${baseUrl}/sessions`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(config),
  });

  if (!response.ok) {
    throw new Error(`create session failed with status ${response.status}`);
  }

  return (await response.json()) as CreateSessionResponse;
}
