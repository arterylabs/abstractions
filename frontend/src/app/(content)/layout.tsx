import LandingPageNavbar from "@/components/landing/Navbar";
import Footer from "@/components/landing/Footer";

export default function ContentLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col min-h-screen">
      <LandingPageNavbar />
      {children}
      <Footer />
    </div>
  );
}
