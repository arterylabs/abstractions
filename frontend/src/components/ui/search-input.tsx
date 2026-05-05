"use client";

import * as React from "react";
import { Input } from "@/components/ui/input";
import Icon from "@/components/ui/icon";
import { cn } from "@/lib/utils";

interface SearchInputProps {
  placeholder?: string;
  className?: string;
  value?: string;
  onChange?: (value: string) => void;
  onSearch?: (value: string) => void;
  debounceMs?: number;
}

export function SearchInput({
  placeholder = "Search...",
  className,
  value: controlledValue,
  onChange: controlledOnChange,
  onSearch,
  debounceMs = 300,
}: SearchInputProps) {
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
    <div className={cn("relative w-full", className)}>
      <Input
        id="search"
        type="text"
        placeholder={placeholder}
        value={value}
        onChange={handleSearch}
        inputClassName="pl-10"
        prepend={
          <Icon
            type="search"
            className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4"
          />
        }
      />
    </div>
  );
}
