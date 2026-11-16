const tokenKey = (import.meta.env.VITE_KEY_PREFIX || "GeekAI_") + "Authorization";

export function getUserToken(): string {
  try {
    const value: unknown = JSON.parse(localStorage.getItem(tokenKey) || '""');
    return typeof value === "string" ? value : "";
  } catch {
    return "";
  }
}

export function removeUserToken(): void {
  localStorage.removeItem(tokenKey);
  localStorage.removeItem("USER_INFO_CACHE_KEY");
}

export function loginURL(location: Pick<Location, "pathname" | "search" | "hash"> = window.location): string {
  const redirect = location.pathname + location.search + location.hash;
  const params = new URLSearchParams({ redirect });
  const invitation = new URLSearchParams(location.search).get("aff");
  if (invitation) params.set("invite_code", invitation);
  return (invitation ? "/register?" : "/login?") + params;
}

export function openLoginPage(): void {
  window.location.assign(loginURL());
}

export const docsURL = import.meta.env.VITE_DOCS_URL?.trim();
