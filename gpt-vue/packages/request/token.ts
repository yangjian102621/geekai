import { getUUID } from "../utils";

const _tokenData = new Map();
export default function tokenHandler() {
  const refreshToken = () => {
    _tokenData.set("__token", getUUID());
    _tokenData.set("lastRequest", null);
  };
  const setCurRequest = (curRequest?: string) => {
    _tokenData.set("lastRequest", curRequest);
  };
  return { _tokenData, refreshToken, setCurRequest };
}
