import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center justify-center rounded-md border font-medium w-fit whitespace-nowrap shrink-0 gap-1 [&>svg]:pointer-events-none ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 focus-visible:ring-4 focus-visible:outline-1 aria-invalid:focus-visible:ring-0 transition-[color,box-shadow]",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
        secondary:
          "border-transparent bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
        destructive:
          "border-transparent bg-destructive text-destructive-foreground [a&]:hover:bg-destructive/90",
        success:
          "border-transparent bg-success text-success-foreground [a&]:hover:bg-success/90",
        warning:
          "border-transparent bg-warning text-warning-foreground [a&]:hover:bg-warning/90",
        info: "border-transparent bg-info text-info-foreground [a&]:hover:bg-info/90",
        outline:
          "text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
      },
      size: {
        sm: "h-6 px-2 text-[10px] leading-none [&>svg]:size-3",
        default: "h-7 px-2.5 text-xs leading-none [&>svg]:size-3",
        lg: "h-8 px-3 text-sm leading-none [&>svg]:size-4",
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
