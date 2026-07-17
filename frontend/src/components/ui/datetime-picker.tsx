"use client";

import * as React from "react";
import { format, set } from "date-fns";

import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";
import Icon from "@/components/ui/icon";

interface DateTimePickerProps {
  value?: Date;
  onChange: (date: Date | undefined) => void;
  placeholder?: string;
  className?: string;
  id?: string;
}

export function DateTimePicker({
  value,
  onChange,
  placeholder = "Select date & time",
  className,
  id,
}: DateTimePickerProps) {
  const [open, setOpen] = React.useState(false);

  // Derive HH:MM string from value, defaulting to "00:00"
  const timeValue = value
    ? `${String(value.getHours()).padStart(2, "0")}:${String(value.getMinutes()).padStart(2, "0")}`
    : "00:00";

  const handleDaySelect = (day: Date | undefined) => {
    if (!day) {
      onChange(undefined);
      return;
    }
    // Preserve the current time when changing the day
    const [hours, minutes] = timeValue.split(":").map(Number);
    onChange(set(day, { hours, minutes, seconds: 0, milliseconds: 0 }));
  };

  const handleTimeChange = (timeString: string) => {
    const [hours, minutes] = timeString.split(":").map(Number);
    const base = value ?? new Date();
    onChange(set(base, { hours, minutes, seconds: 0, milliseconds: 0 }));
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          id={id}
          data-testid={id ? `${id}-trigger` : undefined}
          className={cn(
            "justify-between font-normal text-left",
            !value && "text-muted-foreground",
            className,
          )}
        >
          {value ? format(value, "PPP, h:mm a") : placeholder}
          <Icon type="calendar_month" className="size-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto overflow-hidden p-0" align="start">
        <Calendar
          mode="single"
          selected={value}
          captionLayout="dropdown"
          defaultMonth={value}
          onSelect={handleDaySelect}
        />
        {/* Time input below the calendar, inside the same popover */}
        <div className="border-t border-border px-3 py-3 flex flex-col gap-2">
          <p className="text-xs text-muted-foreground font-medium">Time</p>
          <Input
            id={`time-input-${id}`}
            type="time"
            value={timeValue}
            onChange={handleTimeChange}
            className="appearance-none [&::-webkit-calendar-picker-indicator]:hidden"
          />
        </div>
        <div className="px-3 pb-3 flex justify-end">
          <Button size="sm" onClick={() => setOpen(false)}>
            Done
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}
