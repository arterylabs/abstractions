import * as React from "react";
import Link from "next/link";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";
import Spinner from "@/components/ui/spinner";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium text-foreground no-underline hover:no-underline [&>a]:text-inherit [&>a]:no-underline hover:[&>a]:no-underline focus:[&>a]:no-underline active:[&>a]:no-underline transition-[color,box-shadow] disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 [&_svg]:shrink-0 ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 focus-visible:ring-4 focus-visible:outline-1 aria-invalid:focus-visible:ring-0 cursor-pointer",
  {
    variants: {
      variant: {
        default:
          "bg-primary !text-primary-foreground shadow-sm hover:bg-primary/90",
        destructive:
          "bg-destructive !text-destructive-foreground shadow-xs hover:bg-[var(--destructive-hover)]",
        outline:
          "border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground",
        secondary:
          "bg-secondary !text-secondary-foreground shadow-xs hover:bg-secondary/80",
        ghost: "hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline",
        warning:
          "bg-warning !text-warning-foreground shadow-xs hover:bg-[var(--warning-hover)]",
        info: "bg-info !text-info-foreground shadow-xs hover:bg-[var(--info-hover)]",
      },
      size: {
        default: "h-9 px-4 py-2 has-[>svg]:px-3",
        sm: "h-8 rounded-md px-3 has-[>svg]:px-2.5",
        lg: "h-10 rounded-md px-6 has-[>svg]:px-4",
        icon: "size-9",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

function Button({
  className,
  variant,
  size,
  asChild = false,
  loading,
  href,
  children,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean;
    loading?: boolean;
    href?: string;
    children: React.ReactNode;
  }) {
  const computedClassName = cn(buttonVariants({ variant, size, className }));

  if (href) {
    return (
      <Link
        href={href}
        data-slot="button"
        className={cn(
          computedClassName,
          // Reset link-specific browser/Tailwind styles so button styles take precedence
          "!text-inherit !no-underline !hover:no-underline !focus:no-underline !active:no-underline",
        )}
        data-loading={loading}
      >
        {children}
      </Link>
    );
  }

  if (asChild) {
    return (
      <Slot
        data-slot="button"
        className={computedClassName}
        data-loading={loading}
        {...props}
      >
        {children}
      </Slot>
    );
  }

  return (
    <button
      data-slot="button"
      className={cn(
        computedClassName,
        loading ? "opacity-70 pointer-events-none" : "",
      )}
      data-loading={loading}
      {...props}
    >
      {loading && <Spinner />} {children}
    </button>
  );
}

export { Button, buttonVariants };
