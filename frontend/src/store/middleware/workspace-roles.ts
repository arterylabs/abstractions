import { createListenerMiddleware } from "@reduxjs/toolkit";
import {
  fetchRoles,
  syncWorkspaces,
  setCurrentWorkspace,
} from "../modules/auth";
import { updateWorkspace } from "../modules/workspace";
import type { RootState } from "../index";

export const workspaceRolesMiddleware = createListenerMiddleware();

// Listen for changes to currentWorkspace in state
workspaceRolesMiddleware.startListening({
  predicate: (action, currentState, previousState) => {
    const current = (currentState as RootState).auth.currentWorkspace;
    const previous = (previousState as RootState).auth.currentWorkspace;
    return current?.id !== previous?.id;
  },
  effect: async (action, listenerApi) => {
    const state = listenerApi.getState() as RootState;
    const currentWorkspace = state.auth.currentWorkspace;
    if (currentWorkspace?.id) {
      await listenerApi.dispatch(fetchRoles(currentWorkspace.id));
    }
  },
});

// Fetch roles when workspace sync is complete and there's a current workspace
workspaceRolesMiddleware.startListening({
  actionCreator: syncWorkspaces.fulfilled,
  effect: async (action, listenerApi) => {
    const state = listenerApi.getState() as RootState;
    const currentWorkspace = state.auth.currentWorkspace;
    if (currentWorkspace) {
      await listenerApi.dispatch(fetchRoles(currentWorkspace.id));
    }
  },
});

// Fetch roles when workspace is updated
workspaceRolesMiddleware.startListening({
  actionCreator: updateWorkspace.fulfilled,
  effect: async (action, listenerApi) => {
    const updatedWorkspace = action.payload;
    const state = listenerApi.getState() as RootState;
    const currentWorkspace = state.auth.currentWorkspace;
    if (currentWorkspace?.id === updatedWorkspace.id) {
      await listenerApi.dispatch(fetchRoles(updatedWorkspace.id));
    }
  },
});

// Add listener for manual workspace changes
workspaceRolesMiddleware.startListening({
  actionCreator: setCurrentWorkspace,
  effect: async (action, listenerApi) => {
    const workspace = action.payload;
    if (workspace?.id) {
      await listenerApi.dispatch(fetchRoles(workspace.id));
    }
  },
});
