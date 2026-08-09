/* eslint-disable @typescript-eslint/no-require-imports */
const nextJest = require("next/jest.js");

const createJestConfig = nextJest({
  // Points to the Next.js app root so next/jest can read next.config.ts and .env files
  dir: "./",
});

/** @type {import('jest').Config} */
const config = {
  coverageProvider: "babel",
  testEnvironment: "jsdom",
  setupFilesAfterEnv: ["<rootDir>/src/setupTests.ts"],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/src/$1",
  },
  testMatch: ["**/*.test.ts", "**/*.test.tsx"],
  collectCoverageFrom: [
    "src/components/ui/**/*.{ts,tsx}",
    "!src/components/ui/__tests__/**",
  ],
};

module.exports = createJestConfig(config);
