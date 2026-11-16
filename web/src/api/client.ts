import axios from "axios";
import { getUserToken } from "../config";

export const api = axios.create({
  baseURL: "/api",
  timeout: 10_000,
  withCredentials: true,
  headers: {
    Accept: "application/json",
  },
});

api.interceptors.request.use((config) => {
  const token = getUserToken();
  if (token) config.headers.Authorization = token;
  return config;
});

export interface InviteSummary {
  invite_code: string;
  invite_url: string;
  invited_count: number;
  total_invite_points: number;
  today_invites: number;
}

export interface ApiError {
  error?: string;
  status?: number;
}

export function isApiError(err: unknown): err is ApiError {
  return typeof err === "object" && err !== null && "response" in err;
}
