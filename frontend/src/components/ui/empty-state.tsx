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
      className="flex flex-col items-center justify-center border border-dashed border-border px-6 py-16 text-center"
      {...rest}
    >
      <Icon type={icon} className="!text-2xl text-primary" />
      <h3 className="mt-5 max-w-[22ch] text-2xl font-medium">{title}</h3>
      <p className="mt-3 mb-7 max-w-[46ch] text-[0.9rem] leading-[1.75] text-muted-foreground">
        {description}
      </p>
      {cta}
    </div>
  );
}
