import { api } from "../../api/client";

export interface Page<T> {
  items: T[];
  total: number;
}

export interface PowerStats {
  total: number;
  today: number;
  balance: number;
}

export interface InviteStats {
  invite_count: number;
  reward_total: number;
  today_invite: number;
  invite_code: string;
  invite_link: string;
}

export interface PowerLog {
  id: number;
  user_id: number;
  username?: string;
  type: number;
  type_str?: string;
  amount: number;
  mark: number;
  balance: number;
  model?: string;
  remark?: string;
  created_at: number;
}

export async function fetchPowerStats(): Promise<PowerStats> {
  const { data } = await api.get<{ data: PowerStats }>("/powerLog/stats");
  return data.data;
}

export async function fetchInviteStats(): Promise<InviteStats> {
  const { data } = await api.get<{ data: InviteStats }>("/invite/stats");
  return data.data;
}

export async function fetchPowerLogs(params: {
  page: number;
  page_size: number;
  date?: [string, string];
}): Promise<Page<PowerLog>> {
  const { data } = await api.post<{ data: Page<PowerLog> }>("/powerLog/list", params);
  return data.data;
}

export function formatTime(value: number | string | undefined, language: string): string {
  if (value === undefined || value === null || value === "") return "—";
  const date = typeof value === "number" ? new Date(value * 1000) : new Date(value);
  if (Number.isNaN(date.getTime())) return "—";
  return date.toLocaleString(language);
}

export function formatNumber(value: number | undefined, language: string): string {
  return (value ?? 0).toLocaleString(language);
}
