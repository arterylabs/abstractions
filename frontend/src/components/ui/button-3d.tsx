import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";
import Spinner from "@/components/ui/spinner";

// 3D effect technique:
// - An inset bottom box-shadow in a darker shade creates the "raised edge" without increasing visual height
// - On hover: background darkens slightly — no size/position change
// - On active: button translates down + shadow shrinks → satisfying tactile press feedback
// - ghost and link remain flat intentionally (no visible base makes physical sense for them)
const button3dVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium text-foreground no-underline hover:no-underline [&>a]:text-inherit [&>a]:no-underline hover:[&>a]:no-underline focus:[&>a]:no-underline active:[&>a]:no-underline disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 [&_svg]:shrink-0 ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 focus-visible:ring-4 focus-visible:outline-1 aria-invalid:focus-visible:ring-0 transition-all duration-100 ease-out cursor-pointer select-none",
  {
    variants: {
      variant: {
        // Primary — purple raised button with darker purple base shadow
        default:
          "bg-primary !text-primary-foreground shadow-[inset_0_-4px_0_0_var(--primary-dark)] hover:bg-[var(--primary-hover)] active:translate-y-0.5 active:shadow-[inset_0_-2px_0_0_var(--primary-dark)]",

        // Destructive — red raised button with darker red base shadow
        destructive:
          "bg-destructive !text-destructive-foreground shadow-[inset_0_-4px_0_0_var(--destructive-dark)] hover:bg-[var(--destructive-hover)] active:translate-y-0.5 active:shadow-[inset_0_-2px_0_0_var(--destructive-dark)]",

        // Outline — bordered button with a subtle border-color bottom shadow
        outline:
          "border border-input bg-background shadow-[inset_0_-4px_0_0_var(--accent-dark)] hover:bg-[var(--accent-hover)] hover:text-accent-foreground active:translate-y-0.5 active:shadow-[inset_0_-2px_0_0_var(--accent-dark)]",

        // Secondary — muted raised button with accent-dark base shadow
        secondary:
          "bg-secondary !text-secondary-foreground shadow-[inset_0_-4px_0_0_var(--accent-dark)] hover:bg-[var(--accent-hover)] active:translate-y-0.5 active:shadow-[inset_0_-2px_0_0_var(--accent-dark)]",

        // Ghost — intentionally flat, minimal hover bg, no shadow depth
        ghost: "hover:bg-accent hover:text-accent-foreground active:scale-95",

        // Link — text-only, no 3D effect makes sense here
        link: "text-primary underline-offset-4 hover:underline active:opacity-70",
      },
      size: {
        default: "h-9 px-4 py-2 has-[>svg]:px-3",
        sm: "h-8 rounded-md px-3 text-xs has-[>svg]:px-2.5",
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

function Button3D({
  className,
  variant,
  size,
  asChild = false,
  loading,
  children,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof button3dVariants> & {
    asChild?: boolean;
    loading?: boolean;
    children: React.ReactNode;
  }) {
  if (asChild) {
    return (
      <Slot
        data-slot="button"
        className={cn(button3dVariants({ variant, size, className }))}
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
        button3dVariants({ variant, size, className }),
        loading ? "opacity-70 pointer-events-none" : "",
      )}
      data-loading={loading}
      {...props}
    >
      {loading && <Spinner />} {children}
    </button>
  );
}

export { Button3D, button3dVariants };
