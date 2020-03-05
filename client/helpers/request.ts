import { getAccessToken, hasToken, setAccessToken } from "./auth";
import jwtDecode from "jwt-decode";

function refreshToken() {
  return fetch("/api/refresh").then(function(res) {
    if (res.ok) {
      return res.json().then(function({ accessToken }) {
        return accessToken;
      });
    }
    throw res.json();
  });
}

function isExpired(token: string): boolean {
  try {
    let { exp } = jwtDecode(token);
    return Date.now() > exp * 1000;
  } catch (err) {
    return true;
  }
}

async function parseInit(init: RequestInit) {
  init.headers = new Headers(init.headers);
  if (hasToken()) {
    let token = getAccessToken();
    try {
      if (isExpired(token)) {
        token = await refreshToken();
        setAccessToken(token);
      }
      init.headers.set("Authorization", "Bearer " + token);
    } catch (err) {}
  }
}

async function request(url: RequestInfo, init: RequestInit = {}) {
  parseInit(init);
  return await fetch(url, init).then(function(res) {
    return res.json().then(function(data) {
      if (res.ok) {
        return data;
      }
      throw data;
    });
  });
}

// eslint-disable
async function requestWithHeaders(
  url: RequestInfo,
  init: RequestInit = {}
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
): Promise<[any, Headers]> {
  parseInit(init);
  return await fetch(url, init).then(function(res) {
    return res.json().then(function(data) {
      if (res.ok) {
        return [data, res.headers];
      }
      throw data;
    });
  });
}

export default request;
export { requestWithHeaders };
