import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";

type HeadingSize = "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

// Page-title sizes. h1 is the default masthead; smaller sizes step down for
// nested page titles. Titles use the body (sans) stack by default — opt into
// the serif display face with `display` for marketing / editorial pages.
const sizeClasses: Record<HeadingSize, string> = {
  h1: "text-[2rem] leading-[1.05]",
  h2: "text-2xl",
  h3: "text-xl",
  h4: "text-lg",
  h5: "text-base",
  h6: "text-sm",
};

interface PageHeaderProps extends Omit<React.ComponentProps<"div">, "title"> {
  title: React.ReactNode;
  description?: React.ReactNode;
  /** Uppercase mono running head set above the title, e.g. "Repository". */
  eyebrow?: React.ReactNode;
  showBackButton?: boolean;
  handleBack?: () => void;
  size?: HeadingSize;
  /** Draw the hairline dateline rule beneath the header. On by default for h1. */
  rule?: boolean;
  /**
   * Use the serif display face for the title. Defaults to false so app tab and
   * page titles stay on the body font. Set true for marketing / editorial pages.
   */
  display?: boolean;
}

function PageHeader({
  title,
  description,
  eyebrow,
  showBackButton = false,
  handleBack,
  size = "h1",
  rule,
  display = false,
  className,
  children,
  ...props
}: PageHeaderProps) {
  const Heading = size;
  const showRule = rule ?? size === "h1";

  return (
    <div
      className={cn(
        "flex flex-col gap-1",
        showRule && "border-b border-border pb-5",
        className,
      )}
      {...props}
    >
      {showBackButton && (
        <Button
          variant="ghost"
          size="sm"
          className="press-label mb-4 w-fit gap-2 px-0 text-muted-foreground hover:bg-transparent hover:text-foreground"
          onClick={handleBack}
        >
          <Icon type="arrow_back" className="!text-sm" />
          Back
        </Button>
      )}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div className="flex min-w-0 flex-col gap-2">
          {eyebrow && (
            <span className="press-label text-primary">{eyebrow}</span>
          )}
          <Heading
            className={cn(
              display ? "press-display" : "font-medium",
              sizeClasses[size],
            )}
          >
            {title}
          </Heading>
          {description && (
            <p className="max-w-[68ch] text-[0.95rem] leading-[1.7] text-muted-foreground">
              {description}
            </p>
          )}
        </div>
        {children && <div className="shrink-0">{children}</div>}
      </div>
    </div>
  );
}

export { PageHeader };
