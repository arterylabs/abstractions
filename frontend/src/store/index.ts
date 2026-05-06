import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./modules/auth";
import onboardingReducer from "./modules/onboardingGuide";
import workspaceUsersReducer from "./modules/workspaceUsers";
import invitationsReducer from "./modules/invitations";
import integrationsReducer from "./modules/integrations";
import repositoriesReducer from "./modules/repositories";
import { workspaceRolesMiddleware } from "./middleware/workspace-roles";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    onboarding: onboardingReducer,
    workspaceUsers: workspaceUsersReducer,
    invitations: invitationsReducer,
    integrations: integrationsReducer,
    repositories: repositoriesReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().prepend(workspaceRolesMiddleware.middleware),
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch;
