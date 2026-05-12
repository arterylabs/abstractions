import axios from "axios";
import Cookies from "js-cookie";
import auth from "./auth";
import users from "./users";
import admin from "./admin";
import subscriptions from "./subscriptions";
import workspaces from "./workspaces";
import invitations from "./invitations";
import integrations from "./integrations";
import github from "./github";
import repositories from "./repositories";
import newsletters from "./newsletters";
import files from "./files";

export const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL, // Replace with your API base URL
  headers: {
    "Content-Type": "application/json",
  },
});

const token = Cookies.get("token");
if (token) {
  axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
}

const utils = {
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

const api = {
  auth,
  users,
  subscriptions,
  utils,
  admin,
  workspaces,
  invitations,
  integrations,
  github,
  repositories,
  newsletters,
  files,
};

export default api;
