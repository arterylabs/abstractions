import "@testing-library/jest-dom";

// window.matchMedia is not implemented in jsdom
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: jest.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(),
    removeListener: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// ResizeObserver is not available in jsdom
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// IntersectionObserver is not available in jsdom
global.IntersectionObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
  root: null,
  rootMargin: "",
  thresholds: [],
}));

// scrollIntoView is not implemented in jsdom — needed by cmdk (Command component)
Element.prototype.scrollIntoView = jest.fn();

// PointerEvent is missing in jsdom — needed by Radix UI primitives
if (typeof window.PointerEvent === "undefined") {
  class PointerEvent extends MouseEvent {
    constructor(type: string, init?: PointerEventInit) {
      super(type, init);
    }
  }
  window.PointerEvent = PointerEvent as typeof window.PointerEvent;
}
