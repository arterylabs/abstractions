import { ReactNode } from "react";
import Icon from "./icon";

interface EmptyStateProps {
  icon: string;
  title: string;
  description: string;
  cta: ReactNode;
  "data-testid"?: string;
}

export function EmptyState({
  icon,
  title,
  description,
  cta,
  ...rest
}: EmptyStateProps) {
  return (
    <div
      className="flex flex-col items-center justify-center py-12 px-4 text-center"
      {...rest}
    >
      <div className="mb-6">
        <Icon type={icon} className="!text-[3rem] text-muted-foreground" />
      </div>
      <h3 className="text-lg font-semibold mb-2">{title}</h3>
      <p className="text-muted-foreground mb-6 max-w-md">{description}</p>
      {cta}
    </div>
  );
}
