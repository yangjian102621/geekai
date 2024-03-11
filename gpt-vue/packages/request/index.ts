import axios from "axios";
import tokenHandler from "./token";

const { _tokenData, refreshToken, setCurRequest } = tokenHandler();

const createInstance = (baseURL: string) => {

  const instance = axios.create({
    baseURL,
    timeout: 10000,
    withCredentials: true,
  });

  instance.interceptors.request.use((config) => {
    if (config.url !== _tokenData.get("lastRequest")) {
      refreshToken();
    }
    if (config.method === "post") {
      setCurRequest(config.url);
      config.headers["request-id"] = _tokenData.get("__token");
    }
    return config;
  });

  return instance;
}

export default createInstance;