import auth from "./auth";
import users from "./users";
import admin from "./admin";
import workspaces from "./workspaces";
import invitations from "./invitations";
import integrations from "./integrations";
import github from "./github";
import repositories from "./repositories";
import newsletters from "./newsletters";
import files from "./files";
import { axiosInstance, tokenUtils } from "./client";

export { axiosInstance };

const api = {
  auth,
  users,
  utils: tokenUtils,
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
