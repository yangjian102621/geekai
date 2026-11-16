import { useQuery } from "@tanstack/react-query";
import { api, InviteSummary } from "../../api/client";

export function useInviteSummary(enabled: boolean) {
  return useQuery<InviteSummary>({
    queryKey: ["invite-summary"],
    queryFn: async () => {
      const code = await api.get<{ data: { code: string } }>("/invite/code");
      const { data: response } = await api.get<{
        data: { invite_count: number; reward_total: number; today_invite: number };
      }>("/invite/stats");
      return {
        invite_code: code.data.data.code,
        invite_url: `${window.location.origin}/invite?aff=${encodeURIComponent(code.data.data.code)}`,
        invited_count: response.data.invite_count,
        total_invite_points: response.data.reward_total,
        today_invites: response.data.today_invite,
      };
    },
    enabled,
    retry: false,
  });
}

export function useInviteReward() {
  return useQuery({
    queryKey: ["invite-reward"],
    queryFn: async () => {
      const { data } = await api.get<{ data: { invite_power?: number } | null }>("/config/get?key=system");
      return data.data?.invite_power ?? null;
    },
  });
}
