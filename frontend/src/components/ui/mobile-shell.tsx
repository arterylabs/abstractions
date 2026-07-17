import { ReactNode } from "react";

interface MobileShellProps {
  children: ReactNode;
  /** Background applied to the phone screen (e.g. the event theme panel background).
   *  Accepts any CSS background value: hex color, gradient, etc. Defaults to white. */
  background?: string;
  /** Width of the phone shell in px. Defaults to 390. */
  width?: number;
}

/**
 * iPhone-style mobile device shell.
 * Wraps content in a realistic phone chrome: Dynamic Island,
 * status bar, and home indicator — all in CSS/Tailwind, no images.
 *
 * Usage:
 *   <MobileShell>
 *     <iframe src="..." style={{ width: 390, height: 844 }} />
 *   </MobileShell>
 */
export default function MobileShell({
  children,
  background,
  width = 390,
}: MobileShellProps) {
  return (
    <div className="relative flex justify-center py-2 select-none">
      {/* Phone body */}
      <div
        className="relative overflow-hidden bg-[#1c1c1e]"
        style={{
          width,
          borderRadius: 50,
          boxShadow:
            "0 0 0 1px #3a3a3c, 0 0 0 3px #1c1c1e, 0 0 0 4px #3a3a3c, 0 20px 60px rgba(0,0,0,0.5)",
        }}
      >
        {/* Screen */}
        <div
          className="relative overflow-hidden"
          style={{ borderRadius: 44, background: background ?? "#ffffff" }}
        >
          {/* Status bar — absolutely positioned so it overlays content without adding height */}
          <div className="absolute top-0 left-0 right-0 flex items-center justify-between px-8 pt-3 pb-1 bg-transparent z-10 pointer-events-none">
            {/* Time */}
            <span className="text-[13px] font-semibold text-black leading-none">
              9:41
            </span>

            {/* Dynamic Island */}
            <div
              className="absolute left-1/2 -translate-x-1/2 top-2 bg-black"
              style={{ width: 120, height: 34, borderRadius: 20 }}
            />

            {/* Right icons: signal + wifi + battery */}
            <div className="flex items-center gap-[5px]">
              {/* Cellular signal */}
              <div className="flex items-end gap-[2px] h-[11px]">
                {[3, 5, 7, 9].map((h, i) => (
                  <div
                    key={i}
                    className="w-[3px] bg-black rounded-[1px]"
                    style={{ height: h, opacity: i < 3 ? 1 : 0.3 }}
                  />
                ))}
              </div>
              {/* Wifi */}
              <svg width="15" height="11" viewBox="0 0 15 11" fill="none">
                <path d="M7.5 8.5a1 1 0 1 1 0 2 1 1 0 0 1 0-2z" fill="black" />
                <path
                  d="M4.5 6.2a4.5 4.5 0 0 1 6 0"
                  stroke="black"
                  strokeWidth="1.3"
                  strokeLinecap="round"
                />
                <path
                  d="M2 3.8a7.5 7.5 0 0 1 11 0"
                  stroke="black"
                  strokeWidth="1.3"
                  strokeLinecap="round"
                />
              </svg>
              {/* Battery */}
              <div className="flex items-center">
                <div
                  className="relative border border-black rounded-[2px] flex items-center"
                  style={{ width: 22, height: 11 }}
                >
                  <div
                    className="bg-black rounded-[1px] m-[2px]"
                    style={{ width: 16, height: 7 }}
                  />
                </div>
                <div
                  className="bg-black rounded-r-[1px] ml-[1px]"
                  style={{ width: 2, height: 5 }}
                />
              </div>
            </div>
          </div>

          {/* Content area */}
          <div className="overflow-hidden">{children}</div>

          {/* Home indicator — absolutely positioned so it overlays content without adding height */}
          <div className="absolute bottom-0 left-0 right-0 flex justify-center pb-2 bg-transparent pointer-events-none z-10">
            <div
              className="bg-black/20 rounded-full"
              style={{ width: 130, height: 5 }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
