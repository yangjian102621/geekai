// @vitest-environment jsdom
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, expect, it, vi } from "vitest";
import { api } from "../../api/client";
import { ModelTokenPanel } from "./ModelTokenPanel";

vi.mock("../../api/client", () => ({ api: { get: vi.fn(), post: vi.fn(), delete: vi.fn() } }));
vi.mock("react-i18next", () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

const secret = "sk-test-only";
const copy = vi.fn().mockResolvedValue(undefined);

beforeEach(() => {
  vi.mocked(api.get).mockResolvedValue({ data: { code: 0, data: { enabled: false } } });
  vi.mocked(api.post).mockResolvedValue({ data: { code: 0, data: { token: secret, created_at: "2026-11-17" } } });
  vi.mocked(api.delete).mockResolvedValue({ data: { code: 0 } });
  vi.spyOn(window, "confirm").mockReturnValue(true);
  Object.defineProperty(navigator, "clipboard", { configurable: true, value: { writeText: copy } });
});
afterEach(() => { cleanup(); vi.restoreAllMocks(); vi.clearAllMocks(); });

function showPanel() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } });
  const view = render(<QueryClientProvider client={client}><ModelTokenPanel userId={7} /></QueryClientProvider>);
  return { client, ...view };
}

it("creates and copies a key without persisting the secret, and clears it on revocation", async () => {
  const { client } = showPanel();
  fireEvent.click(await screen.findByRole("button", { name: "model_api.create" }));
  expect((await screen.findByLabelText("model_api.token") as HTMLInputElement).value).toBe(secret);
  expect(JSON.stringify(client.getQueryData(["model-token", 7]))).not.toContain(secret);
  expect(Object.values(localStorage)).not.toContain(secret);
  fireEvent.click(screen.getByRole("button", { name: "model_api.copy" }));
  await waitFor(() => expect(copy).toHaveBeenCalledWith(secret));
  fireEvent.click(screen.getByRole("button", { name: "model_api.revoke" }));
  await screen.findByRole("button", { name: "model_api.create" });
  expect(api.delete).toHaveBeenCalledWith("/user/model-token");
  expect(screen.queryByLabelText("model_api.token")).toBeNull();
});

it("requires confirmation before invalidating an existing key", async () => {
  vi.mocked(api.get).mockResolvedValue({ data: { code: 0, data: { enabled: true } } });
  vi.mocked(window.confirm).mockReturnValue(false);
  showPanel();
  fireEvent.click(await screen.findByRole("button", { name: "model_api.reset" }));
  fireEvent.click(screen.getByRole("button", { name: "model_api.revoke" }));
  expect(api.post).not.toHaveBeenCalled();
  expect(api.delete).not.toHaveBeenCalled();
});

it("reports business errors without displaying a key", async () => {
  vi.mocked(api.post).mockResolvedValue({ data: { code: 1, message: "unavailable" } });
  showPanel();
  fireEvent.click(await screen.findByRole("button", { name: "model_api.create" }));
  expect((await screen.findByRole("alert")).textContent).toBe("model_api.action_error");
  expect(screen.queryByLabelText("model_api.token")).toBeNull();
});
