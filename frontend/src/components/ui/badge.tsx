import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const badgeVariants = cva(
  // Set as a mono folio — uppercase, letterspaced, square — so status reads as
  // a stamp on the page rather than a pill.
  "inline-flex items-center justify-center rounded-sm border font-mono font-medium uppercase tracking-[0.14em] w-fit whitespace-nowrap shrink-0 gap-1 [&>svg]:pointer-events-none ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 focus-visible:ring-4 focus-visible:outline-1 aria-invalid:focus-visible:ring-0 transition-[color,box-shadow]",
  {
    variants: {
      variant: {
        // Stamps, not chips: a hairline box and coloured ink on a faint ground.
        // Solid fills read far too loud against the near-black editorial page.
        default:
          "border-primary/40 bg-primary/10 text-primary [a&]:hover:bg-primary/15",
        secondary:
          "border-border bg-accent text-muted-foreground [a&]:hover:bg-accent-light",
        destructive:
          "border-destructive-dark/40 bg-destructive-dark/10 text-destructive-dark [a&]:hover:bg-destructive-dark/15",
        success:
          "border-success-dark/40 bg-success-dark/10 text-success-dark [a&]:hover:bg-success-dark/15",
        warning:
          "border-warning-dark/40 bg-warning-dark/10 text-warning-dark [a&]:hover:bg-warning-dark/15",
        info: "border-info-dark/40 bg-info-dark/10 text-info-dark [a&]:hover:bg-info-dark/15",
        outline:
          "border-border text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
      },
      size: {
        sm: "h-5 px-1.5 text-[9px] leading-none [&>svg]:size-3",
        default: "h-6 px-2 text-[10px] leading-none [&>svg]:size-3",
        lg: "h-7 px-2.5 text-xs leading-none [&>svg]:size-4",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

type BadgeVariant = NonNullable<VariantProps<typeof badgeVariants>["variant"]>;

function Badge({
  className,
  variant,
  size,
  enableShadow = false,
  asChild = false,
  ...props
}: React.ComponentProps<"span"> &
  VariantProps<typeof badgeVariants> & {
    asChild?: boolean;
    enableShadow?: boolean;
  }) {
  const Comp = asChild ? Slot : "span";
  const effectiveVariant: BadgeVariant = (variant ?? "default") as BadgeVariant;
  const shadowEligibleVariants: Record<BadgeVariant, boolean> = {
    default: true,
    destructive: true,
    success: true,
    info: true,
    secondary: false,
    warning: false,
    outline: false,
  };

  return (
    <Comp
      data-slot="badge"
      // Stable hook for tests and styling overrides, so they don't have to
      // assert on whichever utility classes the variant happens to use.
      data-variant={effectiveVariant}
      className={cn(
        badgeVariants({ variant: effectiveVariant, size }),
        enableShadow && shadowEligibleVariants[effectiveVariant] && "shadow-sm",
        className,
      )}
      {...props}
    />
  );
}

export { Badge, badgeVariants };
