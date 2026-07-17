import type { Metadata } from "next";
import GodModeLayoutInner from "./GodModeLayoutInner";

export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default function GodModeLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <GodModeLayoutInner>{children}</GodModeLayoutInner>;
}
