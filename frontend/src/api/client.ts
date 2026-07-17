import axios from "axios";
import Cookies from "js-cookie";

/**
 * Shared axios instance for all domain API modules.
 *
 * Lives outside `index.ts` so domain modules can import it without creating a
 * circular dependency (index → domain → index), which surfaces as a TDZ
 * ReferenceError under some browsers/bundler chunk orders (notably WebKit).
 */
export const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

const token = Cookies.get("token");
if (token) {
  axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
}

export const tokenUtils = {
  setToken: (token: string, expiresAt: string) => {
    axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
    Cookies.set("token", token, { expires: new Date(expiresAt) });
  },
  getToken: () => {
    return Cookies.get("token");
  },
  removeToken: () => {
    axiosInstance.defaults.headers.common["Authorization"] = "";
    Cookies.remove("token");
  },
};
