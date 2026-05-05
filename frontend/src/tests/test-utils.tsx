import React from "react";
import { render, RenderOptions } from "@testing-library/react";
import { Provider } from "react-redux";
import { configureStore } from "@reduxjs/toolkit";

// Re-export everything from RTL for convenience
export * from "@testing-library/react";

// Minimal store for tests — add slices here if your component tests need state
const testStore = configureStore({ reducer: {} });

interface RenderWithProvidersOptions extends Omit<RenderOptions, "wrapper"> {
  store?: typeof testStore;
}

/**
 * Custom render function that wraps the component under test with global
 * providers (Redux store). Use this instead of RTL's `render` when the
 * component tree accesses the Redux store.
 *
 * Usage:
 *   import { renderWithProviders } from "@/tests/test-utils";
 *   renderWithProviders(<MyComponent />);
 */
export function renderWithProviders(
  ui: React.ReactElement,
  { store = testStore, ...renderOptions }: RenderWithProvidersOptions = {},
) {
  function Wrapper({ children }: { children: React.ReactNode }) {
    return <Provider store={store}>{children}</Provider>;
  }
  return render(ui, { wrapper: Wrapper, ...renderOptions });
}
