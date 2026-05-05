"use client";

import * as React from "react";
import { Input } from "./input";
import { cn } from "@/lib/utils";

interface SearchBarProps {
  placeholder?: string;
  className?: string;
  value?: string;
  onChange?: (value: string) => void;
  onSearch?: (value: string) => void;
  debounceMs?: number;
}

export function SearchBar({
  placeholder = "Search...",
  className,
  value: controlledValue,
  onChange: controlledOnChange,
  onSearch,
  debounceMs = 300,
}: SearchBarProps) {
  const [internalValue, setInternalValue] = React.useState("");
  const [searchTimeout, setSearchTimeout] = React.useState<NodeJS.Timeout>();

  const value = controlledValue ?? internalValue;
  const onChange = controlledOnChange ?? setInternalValue;

  const handleSearch = React.useCallback(
    (searchValue: string) => {
      onChange(searchValue);

      if (onSearch) {
        if (searchTimeout) {
          clearTimeout(searchTimeout);
        }
        setSearchTimeout(
          setTimeout(() => {
            onSearch(searchValue);
          }, debounceMs),
        );
      }
    },
    [onChange, onSearch, searchTimeout, debounceMs],
  );

  return (
    <Input
      id="search"
      type="search"
      placeholder={placeholder}
      value={value}
      onChange={handleSearch}
      className={cn("w-full", className)}
    />
  );
}
