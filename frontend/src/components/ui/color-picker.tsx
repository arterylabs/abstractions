"use client";

import * as React from "react";
import { HexColorPicker } from "react-colorful";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const HEX_REGEX = /^#[0-9a-fA-F]{6}$/;

interface ColorPickerProps {
  id?: string;
  value: string;
  onChange: (value: string) => void;
  className?: string;
}

function ColorPicker({ id, value, onChange, className }: ColorPickerProps) {
  // Local state lets the text input show in-progress typing without
  // pushing invalid values upstream. We sync back when the user blurs
  // or types a full valid hex color.
  const [inputValue, setInputValue] = React.useState(value);

  // Keep the text input in sync when the parent changes the value
  // (e.g. when the color swatch picker fires onChange directly).
  React.useEffect(() => {
    setInputValue(value);
  }, [value]);

  function handleInputChange(raw: string) {
    setInputValue(raw);
    // Prefix a '#' if the user forgot it
    const normalized = raw.startsWith("#") ? raw : `#${raw}`;
    if (HEX_REGEX.test(normalized)) {
      onChange(normalized);
    }
  }

  function handleInputBlur() {
    // On blur, reset to the last valid value so the field is never left in an
    // invalid state visible to the user.
    if (!HEX_REGEX.test(inputValue)) {
      setInputValue(value);
    }
  }

  return (
    <div className={cn("flex items-center gap-2", className)}>
      <Popover>
        <PopoverTrigger asChild>
          <button
            type="button"
            aria-label={`Pick color, current value ${value}`}
            className={cn(
              "h-9 w-9 shrink-0 rounded-md border border-input shadow-xs transition-opacity",
              "focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-ring/20 focus-visible:outline-1 focus-visible:outline-ring/50",
              "disabled:pointer-events-none disabled:opacity-50",
            )}
            style={{
              backgroundColor: HEX_REGEX.test(value) ? value : "#000000",
            }}
          />
        </PopoverTrigger>
        <PopoverContent className="w-auto p-3">
          <HexColorPicker color={value} onChange={onChange} />
        </PopoverContent>
      </Popover>

      <Input
        id={id ?? "color-picker-input"}
        value={inputValue}
        onChange={handleInputChange}
        onBlur={handleInputBlur}
        placeholder="#000000"
        inputClassName="font-mono"
        className="flex-1"
      />
    </div>
  );
}

export { ColorPicker };
