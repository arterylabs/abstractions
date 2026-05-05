import { ReactNode } from "react";

interface EmptyStateProps {
  image?: ReactNode;
  title: string;
  description: string;
  cta: ReactNode;
  "data-testid"?: string;
}

export function EmptyState({
  image,
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
      {image && <div className="mb-6 text-[48px]">{image}</div>}
      <h3 className="text-lg font-semibold mb-2">{title}</h3>
      <p className="text-muted-foreground mb-6 max-w-md">{description}</p>
      {cta}
    </div>
  );
}
