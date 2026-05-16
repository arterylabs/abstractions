import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";

type HeadingSize = "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

const sizeClasses: Record<HeadingSize, string> = {
  h1: "text-2xl",
  h2: "text-xl",
  h3: "text-lg",
  h4: "text-base",
  h5: "text-sm",
  h6: "text-xs",
};

interface PageHeaderProps extends Omit<React.ComponentProps<"div">, "title"> {
  title: React.ReactNode;
  description?: React.ReactNode;
  showBackButton?: boolean;
  handleBack?: () => void;
  size?: HeadingSize;
}

function PageHeader({
  title,
  description,
  showBackButton = false,
  handleBack,
  size = "h1",
  className,
  children,
  ...props
}: PageHeaderProps) {
  const Heading = size;
  return (
    <div
      className={cn("flex items-center justify-between relative", className)}
      {...props}
    >
      <div className="flex items-start gap-4">
        {showBackButton && (
          <Button
            variant="outline"
            className="absolute -top-12"
            size="sm"
            onClick={handleBack}
          >
            <Icon type="arrow_back" className="size-4" />
          </Button>
        )}
        <div className="flex flex-col m-0">
          <Heading className={cn("font-medium", sizeClasses[size])}>
            {title}
          </Heading>
          {description && (
            <p className="mt-2 text-muted-foreground">{description}</p>
          )}
        </div>
      </div>
      {children}
    </div>
  );
}

export { PageHeader };
