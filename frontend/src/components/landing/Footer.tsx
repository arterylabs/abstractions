import Link from "next/link";
import { Separator } from "../ui/separator";

export default function Footer() {
  return (
    <footer className="flex flex-col gap-2 sm:flex-row py-6 w-full shrink-0 items-center px-4 md:px-6 border-t">
      <p className="text-xs text-gray-500 dark:text-gray-400">
        © 2023 Abstractions. All rights reserved.
      </p>
      <Separator className="sm:hidden h-[1px] w-full bg-border" />
      <nav className="sm:ml-auto flex gap-4 sm:gap-6">
        <Link
          className="text-xs hover:underline underline-offset-4 hover:text-primary"
          href="/terms"
        >
          Terms of Service
        </Link>
        <Link
          className="text-xs hover:underline underline-offset-4 hover:text-primary"
          href="/privacy"
        >
          Privacy
        </Link>
      </nav>
    </footer>
  );
}
