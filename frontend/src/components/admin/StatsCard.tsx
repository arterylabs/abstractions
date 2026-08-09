import { Card, CardContent } from "@/components/ui/card";
import Icon from "@/components/ui/icon";
import { cn } from "@/lib/utils";

interface StatsCardProps {
  title: string;
  value: number;
  icon?: string;
  className?: string;
}

export default function StatsCard({
  title,
  value,
  icon,
  className,
}: StatsCardProps) {
  return (
    <Card className={cn("overflow-hidden", className)}>
      <CardContent className="p-5">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="press-label text-muted-foreground">{title}</p>
            <p className="press-display mt-2 text-4xl tabular-nums">
              {value.toLocaleString()}
            </p>
          </div>
          {icon && (
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-sm border border-violet/30 text-violet-light">
              <Icon type={icon} className="text-lg" />
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
