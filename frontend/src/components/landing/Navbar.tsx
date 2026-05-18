import Link from "next/link";
import Logo from "@/components/ui/logo";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import Icon from "@/components/ui/icon";

const NavLinks = ({ mobile = false }: { mobile?: boolean }) => {
  const links = [
    { href: "/#features", label: "Features" },
    { href: "/#how-it-works", label: "How It Works" },
  ];

  return (
    <div
      className={`${
        mobile ? "flex flex-col space-y-4" : "hidden lg:flex items-center gap-6"
      }`}
    >
      {links.map((link) => (
        <Link key={link.href} href={link.href}>
          <Button variant="ghost" className="text-gray-600 hover:text-gray-900">
            {link.label}
          </Button>
        </Link>
      ))}
    </div>
  );
};

export default function Navbar() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container h-14 max-w-[1000px] mx-auto flex items-center justify-between">
        <Link className="flex items-center space-x-2" href="/#">
          <Logo className="h-8 w-auto" />
        </Link>

        <NavLinks />

        <div className="flex items-center gap-4">
          <Link href="/login" className="hidden lg:block">
            <Button variant="ghost">Login</Button>
          </Link>
          <Link href="/signup" className="hidden lg:block">
            <Button>Sign up</Button>
          </Link>

          <Sheet>
            <SheetTrigger asChild className="lg:hidden">
              <Button variant="ghost" size="icon">
                <Icon type="menu" className="h-5 w-5" />
              </Button>
            </SheetTrigger>
            <SheetContent>
              <div className="flex flex-col space-y-4 mt-8">
                <NavLinks mobile />
                <Link href="/login">
                  <Button variant="ghost" className="w-full">
                    Login
                  </Button>
                </Link>
                <Link href="/signup">
                  <Button className="w-full">Sign up</Button>
                </Link>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  );
}
