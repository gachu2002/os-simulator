interface APIErrorEnvelope {
  message?: string;
}

function buildURL(baseURL: string, path: string): string {
  const trimmedBase = baseURL.trim().replace(/\/$/, "");
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  return `${trimmedBase}${normalizedPath}`;
}

export async function fetchJSON<T>(
  baseURL: string,
  path: string,
  init?: Parameters<typeof fetch>[1],
): Promise<T> {
  const response = await fetch(buildURL(baseURL, path), init);
  if (!response.ok) {
    throw new Error(await decodeErrorMessage(response));
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
}

async function decodeErrorMessage(response: Response): Promise<string> {
  try {
    const payload = (await response.json()) as APIErrorEnvelope;
    if (typeof payload.message === "string" && payload.message.trim() !== "") {
      return payload.message;
    }
  } catch {
    // ignore decode errors and use fallback
  }
  return `request failed with status ${response.status}`;
}
