import { SectionMark } from "@/components/ui/section-mark";
import { Switch } from "@/components/ui/switch";

/* Shared building blocks for the repository settings surfaces: repository
   Settings, and Newsletters' Schedule and Recipients sub-tabs. */
export function Section({
  label,
  description,
  action,
  children,
}: {
  label: string;
  description?: string;
  action?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <section className="flex flex-col gap-4">
      <SectionMark label={label} action={action} />
      <div className="flex flex-col gap-5 border border-border bg-card p-5">
        {description && (
          <p className="max-w-[68ch] text-sm leading-relaxed text-muted-foreground">
            {description}
          </p>
        )}
        {children}
      </div>
    </section>
  );
}

/** A switch with its own explanation, boxed off from the rest of a section. */
export function SwitchRow({
  title,
  description,
  checked,
  onCheckedChange,
}: {
  title: string;
  description: string;
  checked: boolean;
  onCheckedChange: (checked: boolean) => void;
}) {
  return (
    <div className="flex items-center justify-between gap-4 border border-border p-4">
      <div className="flex flex-col gap-1.5">
        <p className="press-label text-foreground">{title}</p>
        <p className="max-w-[60ch] text-sm leading-relaxed text-muted-foreground">
          {description}
        </p>
      </div>
      <Switch checked={checked} onCheckedChange={onCheckedChange} />
    </div>
  );
}
