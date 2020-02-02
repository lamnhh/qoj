function range(left: number, right: number): number[] {
  let ans = [];
  for (let i = left; i <= right; ++i) {
    ans.push(i);
  }
  return ans;
}

function parsePage(page: string): number {
  let ans = parseInt(page);
  if (isNaN(ans)) {
    return 1;
  }
  return ans;
}

function buildURL(url: string, params: Array<[string, string]>) {
  let paramString = params.map(([key, value]) => `${key}=${value}`).join("&");
  return url + "?" + paramString;
}

export { range, parsePage, buildURL };
