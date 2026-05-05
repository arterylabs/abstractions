import { format, formatDistanceToNow } from "date-fns";
import Tooltip from "@/components/ui/tooltip";

export type TimestampFormat =
  | "PPpp" // "Apr 29, 2023 at 3:45 PM"
  | "PP" // "Apr 29, 2023"
  | "PPP" // "April 29th, 2023"
  | "p" // "3:45 PM"
  | "Pp" // "Apr 29, 2023, 3:45 PM"
  | "yyyy-MM-dd" // "2023-04-29"
  | "MMM dd, yyyy" // "Apr 29, 2023"
  | "MMMM dd, yyyy" // "April 29, 2023"
  | "MM/dd/yyyy" // "04/29/2023"
  | "dd/MM/yyyy" // "29/04/2023"
  | "HH:mm" // "15:45"
  | "h:mm a" // "3:45 PM"
  | "time-ago"; // "2 hours ago"

interface TimestampProps {
  timestamp: string | Date;
  format?: TimestampFormat;
  showTime?: boolean;
  className?: string;
}

export default function Timestamp({
  timestamp,
  format: formatType = "time-ago",
  showTime = true,
  className,
}: TimestampProps) {
  const date = new Date(timestamp);

  // Check if date is valid
  if (isNaN(date.getTime())) {
    return (
      <span className={className || "text-xs text-muted-foreground"}>
        Invalid date
      </span>
    );
  }

  const formatTimestamp = (date: Date, formatType: TimestampFormat): string => {
    if (formatType === "time-ago") {
      return formatDistanceToNow(date, { addSuffix: true });
    }

    try {
      return format(date, formatType);
    } catch {
      // Fallback to default format if custom format fails
      return format(date, "MM/dd/yyyy");
    }
  };

  const formattedTimestamp = formatTimestamp(date, formatType);

  // Create tooltip content showing full date and time in user's timezone
  const tooltipContent = showTime ? (
    <div className="text-center">
      <div className="font-medium">{format(date, "EEEE, MMMM do, yyyy")}</div>
      <div className="text-xs opacity-80">{format(date, "h:mm:ss a O")}</div>
    </div>
  ) : (
    format(date, "EEEE, MMMM do, yyyy")
  );

  return (
    <Tooltip value={tooltipContent}>
      <span
        className={className || "text-xs text-muted-foreground cursor-help"}
      >
        {formattedTimestamp}
      </span>
    </Tooltip>
  );
}
