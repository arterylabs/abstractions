import LandingPageNavbar from "@/components/landing/Navbar";

export default function ContentLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col min-h-screen">
      <LandingPageNavbar />
      {children}
    </div>
  );
}
