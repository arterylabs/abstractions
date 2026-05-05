import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";

interface PageHeaderProps extends React.ComponentProps<"div"> {
  title: string;
  description?: string;
  showBackButton?: boolean;
  handleBack?: () => void;
}

function PageHeader({
  title,
  description,
  showBackButton = false,
  handleBack,
  className,
  children,
  ...props
}: PageHeaderProps) {
  return (
    <div
      className={
        cn("", className) + ` flex items-center justify-between relative`
      }
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
          <h1 className="text-2xl font-medium">{title}</h1>
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
