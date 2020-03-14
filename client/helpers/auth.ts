let accessToken: string | null = null;

export function hasToken(): boolean {
  return accessToken !== null;
}

export function getAccessToken(): string {
  if (!hasToken() || !accessToken) {
    throw new Error("No token available");
  }
  return accessToken;
}

export function setAccessToken(tkn: string): void {
  accessToken = tkn;
}

export function clearToken(): void {
  accessToken = null;
  fetch("/api/logout");
}
