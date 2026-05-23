import { type ReactNode } from "react";
import Icon from "@/components/ui/icon";

type CalloutType = "note" | "tip" | "warning" | "danger";

interface CalloutProps {
  type?: CalloutType;
  title?: string;
  children: ReactNode;
}

const config: Record<
  CalloutType,
  { icon: string; className: string; titleDefault: string }
> = {
  note: {
    icon: "info",
    className: "border-info bg-info-light text-info-foreground",
    titleDefault: "Note",
  },
  tip: {
    icon: "lightbulb",
    className: "border-success bg-success-light text-success-foreground",
    titleDefault: "Tip",
  },
  warning: {
    icon: "warning",
    className: "border-warning bg-warning-light text-warning-foreground",
    titleDefault: "Warning",
  },
  danger: {
    icon: "error",
    className:
      "border-destructive bg-destructive-light text-destructive-foreground",
    titleDefault: "Danger",
  },
};

export default function Callout({
  type = "note",
  title,
  children,
}: CalloutProps) {
  const { icon, className, titleDefault } = config[type];
  const label = title ?? titleDefault;

  return (
    <div
      className={`my-4 flex gap-3 rounded-lg border px-4 py-3 [&_a]:text-current [&_a]:underline ${className} callout callout-${type}`}
    >
      <Icon type={icon} className="text-[1.1rem] mt-0.5 shrink-0 opacity-80" />
      <div className="flex-1 min-w-0">
        <div className="!m-0 font-semibold text-sm mb-4">{label}</div>
        <div className="text-sm leading-relaxed [&>p]:!m-0 [&>p]:leading-relaxed">
          {children}
        </div>
      </div>
    </div>
  );
}

export type { CalloutType, CalloutProps };
