// @vitest-environment jsdom
import { afterEach, describe, expect, it } from "vitest";
import { getUserToken, loginURL, removeUserToken } from "./config";

afterEach(() => localStorage.clear());

describe("GeekAI account integration", () => {
  it("reads the same token format as the existing Vue login", () => {
    localStorage.setItem("GeekAI_Authorization", JSON.stringify("user-token"));
    expect(getUserToken()).toBe("user-token");
    localStorage.setItem("GeekAI_Admin-Authorization", JSON.stringify("admin-token"));
    removeUserToken();
    expect(getUserToken()).toBe("");
    expect(localStorage.getItem("GeekAI_Admin-Authorization")).toBe(JSON.stringify("admin-token"));
  });

  it("ignores a damaged token instead of preventing the website from rendering", () => {
    localStorage.setItem("GeekAI_Authorization", "not-json");
    expect(getUserToken()).toBe("");
    localStorage.setItem("GeekAI_Authorization", JSON.stringify({ token: "wrong-format" }));
    expect(getUserToken()).toBe("");
  });

  it("preserves invitation codes when opening registration", () => {
    const url = new URL(loginURL({ pathname: "/invite", search: "?aff=A%2BB", hash: "#share" }), "https://site.example");
    expect(url.pathname).toBe("/register");
    expect(url.searchParams.get("invite_code")).toBe("A+B");
    expect(url.searchParams.get("redirect")).toBe("/invite?aff=A%2BB#share");
  });

  it("opens local login without a Tokengine configuration", () => {
    expect(loginURL({ pathname: "/", search: "", hash: "" })).toBe("/login?redirect=%2F");
  });
});
