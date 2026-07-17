import type { Metadata } from "next";
import AuthLayoutInner from "./AuthLayoutInner";

export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <AuthLayoutInner>{children}</AuthLayoutInner>;
}
