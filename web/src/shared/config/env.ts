export function apiBaseURL(): string {
  const envURL = import.meta.env.VITE_API_BASE_URL;
  if (typeof envURL === "string" && envURL.trim() !== "") {
    return envURL.trim();
  }
  if (typeof window === "undefined") {
    return "http://localhost:8080";
  }
  return window.location.origin;
}
