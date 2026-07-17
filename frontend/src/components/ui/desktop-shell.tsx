import { ReactNode } from "react";

interface DesktopShellProps {
  children: ReactNode;
  /** URL shown in the address bar, e.g. "harmony.app/e/my-event" */
  url?: string;
  /** Background applied to the content area (e.g. the event theme panel background).
   *  Accepts any CSS background value: hex color, gradient, etc. Defaults to white. */
  background?: string;
}

/**
 * Google Chrome (light mode) browser shell.
 * Row 1: traffic lights + tab strip (active tab in white).
 * Row 2: navigation buttons + address bar.
 * All in CSS/Tailwind, no images.
 *
 * Usage:
 *   <DesktopShell url="harmony.app/e/my-event">
 *     <iframe src="..." style={{ width: 1280, height: 900 }} />
 *   </DesktopShell>
 */
export default function DesktopShell({
  children,
  url,
  background,
}: DesktopShellProps) {
  const displayUrl = url ?? "harmony.app";

  return (
    <div
      className="overflow-hidden rounded-xl shadow-lg select-none"
      style={{ border: "1px solid #c8cacd" }}
    >
      {/* Row 1: traffic lights + tabs */}
      <div
        className="flex items-end px-3 gap-2"
        style={{ height: 40, background: "#dee1e6" }}
      >
        {/* Traffic lights */}
        <div className="flex items-center gap-[6px] shrink-0 pb-2">
          <div
            className="rounded-full"
            style={{
              width: 12,
              height: 12,
              background: "#ff5f57",
              boxShadow: "inset 0 0 0 0.5px rgba(0,0,0,0.15)",
            }}
          />
          <div
            className="rounded-full"
            style={{
              width: 12,
              height: 12,
              background: "#febc2e",
              boxShadow: "inset 0 0 0 0.5px rgba(0,0,0,0.15)",
            }}
          />
          <div
            className="rounded-full"
            style={{
              width: 12,
              height: 12,
              background: "#28c840",
              boxShadow: "inset 0 0 0 0.5px rgba(0,0,0,0.15)",
            }}
          />
        </div>

        {/* Active tab */}
        <div
          className="flex items-center gap-2 px-3 shrink-0"
          style={{
            height: 32,
            maxWidth: 220,
            background: "#ffffff",
            borderRadius: "8px 8px 0 0",
            boxShadow: "0 1px 0 #ffffff",
          }}
        >
          <div
            className="rounded-sm shrink-0"
            style={{ width: 14, height: 14, background: "#e8eaed" }}
          />
          <span
            className="text-[12px] truncate"
            style={{ color: "#202124", fontWeight: 500, maxWidth: 160 }}
          >
            {displayUrl}
          </span>
        </div>

        {/* Inactive tab */}
        <div
          className="flex items-center gap-2 px-3 shrink-0 pb-1"
          style={{ height: 30, maxWidth: 160 }}
        >
          <div
            className="rounded-sm shrink-0"
            style={{ width: 14, height: 14, background: "#bdc1c6" }}
          />
          <span className="text-[12px] truncate" style={{ color: "#5f6368" }}>
            New Tab
          </span>
        </div>
      </div>

      {/* Row 2: toolbar with address bar */}
      <div
        className="flex items-center gap-2 px-3"
        style={{
          height: 44,
          background: "#ffffff",
          borderBottom: "1px solid #e0e0e0",
        }}
      >
        {/* Back / forward / refresh buttons */}
        <div className="flex items-center gap-1 shrink-0">
          {/* Back */}
          <div
            className="flex items-center justify-center rounded-full"
            style={{ width: 28, height: 28 }}
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path
                d="M10 3L5 8l5 5"
                stroke="#5f6368"
                strokeWidth="1.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
          {/* Forward */}
          <div
            className="flex items-center justify-center rounded-full"
            style={{ width: 28, height: 28 }}
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path
                d="M6 3l5 5-5 5"
                stroke="#bdc1c6"
                strokeWidth="1.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
          {/* Refresh */}
          <div
            className="flex items-center justify-center rounded-full"
            style={{ width: 28, height: 28 }}
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <path
                d="M13 8A5 5 0 1 1 8 3h3M11 1v4H7"
                stroke="#5f6368"
                strokeWidth="1.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </div>
        </div>

        {/* Address bar */}
        <div
          className="flex-1 flex items-center gap-2 px-3"
          style={{
            height: 32,
            background: "#f1f3f4",
            borderRadius: 16,
            border: "1px solid #e0e0e0",
          }}
        >
          {/* Lock icon */}
          <svg
            width="10"
            height="12"
            viewBox="0 0 10 12"
            fill="none"
            style={{ flexShrink: 0 }}
          >
            <rect x="1" y="5" width="8" height="7" rx="1.5" fill="#5f6368" />
            <path
              d="M3 5V3.5a2 2 0 0 1 4 0V5"
              stroke="#5f6368"
              strokeWidth="1.5"
              strokeLinecap="round"
            />
          </svg>
          <span className="text-[12px] truncate" style={{ color: "#202124" }}>
            {displayUrl}
          </span>
        </div>

        {/* Right icons placeholder */}
        <div className="flex items-center gap-1 shrink-0">
          <div
            className="rounded-full"
            style={{ width: 28, height: 28, background: "#f1f3f4" }}
          />
          <div
            className="rounded-full"
            style={{ width: 28, height: 28, background: "#f1f3f4" }}
          />
        </div>
      </div>

      {/* Content */}
      <div
        className="overflow-hidden"
        style={{ background: background ?? "#ffffff" }}
      >
        {children}
      </div>
    </div>
  );
}
