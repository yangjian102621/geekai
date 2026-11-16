import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { api } from "../../api/client";
import { getUserToken, removeUserToken } from "../../config";

export interface AuthUser {
  id: number;
  username: string;
  power?: number;
}

export function useAuth() {
  return useQuery<AuthUser | null>({
    queryKey: ["auth-user"],
    queryFn: async () => {
      if (!getUserToken()) return null;
      try {
        const response = await api.get<{ data: AuthUser }>("/user/session");
        return response.data.data;
      } catch (error) {
        if (axios.isAxiosError(error) && error.response?.status === 401) {
          removeUserToken();
          return null;
        }
        throw error;
      }
    },
    staleTime: 60_000,
  });
}
