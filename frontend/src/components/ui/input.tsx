"use client";

import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";
import { Label } from "@/components/ui/label";
import Icon from "@/components/ui/icon";
import Tooltip from "@/components/ui/tooltip";

const inputVariants = cva(
  "border-input file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground aria-invalid:outline-destructive/60 aria-invalid:ring-destructive/20 dark:aria-invalid:outline-destructive dark:aria-invalid:ring-destructive/50 ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 aria-invalid:outline-destructive/60 dark:aria-invalid:outline-destructive dark:aria-invalid:ring-destructive/40 aria-invalid:ring-destructive/20 aria-invalid:border-destructive/60 dark:aria-invalid:border-destructive flex w-full min-w-0 rounded-md border bg-accent shadow-xs transition-[color,box-shadow] file:inline-flex file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-4 focus-visible:outline-1 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:focus-visible:ring-[3px] aria-invalid:focus-visible:outline-none dark:aria-invalid:focus-visible:ring-4",
  {
    variants: {
      size: {
        xs: "h-7 px-2 py-1 text-xs file:h-5",
        sm: "h-8 px-3 py-1 text-sm file:h-6",
        default: "h-9 px-3 py-1 text-base file:h-7 md:text-sm",
        lg: "h-10 px-4 py-2 text-base file:h-8",
      },
    },
    defaultVariants: {
      size: "default",
    },
  },
);

function Input({
  id,
  placeholder,
  className,
  inputClassName,
  labelClassName,
  value,
  onChange,
  label,
  disabled,
  required = false,
  type = "text",
  labelAction,
  step,
  toggleVisibility,
  prepend,
  size,
  ...rest
}: {
  id: string;
  placeholder?: string;
  value: string;
  onChange: (value: string) => void;
  inputClassName?: string;
  labelClassName?: string;
  label?: string;
  required?: boolean;
  type?: string;
  disabled?: boolean;
  className?: string;
  labelAction?: React.ReactNode;
  toggleVisibility?: boolean;
  step?: string;
  prepend?: React.ReactNode;
  size?: VariantProps<typeof inputVariants>["size"];
} & Omit<
  React.ComponentPropsWithoutRef<"input">,
  | "id"
  | "placeholder"
  | "className"
  | "value"
  | "onChange"
  | "disabled"
  | "required"
  | "type"
  | "step"
  | "size"
>) {
  const [inputVisible, setInputVisible] = React.useState(false);
  const toggle = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setInputVisible(!inputVisible);
  };

  let inputType = type;
  if (toggleVisibility && type === "password") {
    inputType = inputVisible ? "text" : "password";
  }

  return (
    <div className={cn("grid gap-2", className)}>
      {label && (
        <div className="flex items-center justify-between text-sm leading-none select-none">
          <Label className={labelClassName} htmlFor={id}>
            {label}
          </Label>
          {labelAction}
        </div>
      )}
      <div className="relative">
        {prepend}
        <input
          id={id}
          type={inputType}
          placeholder={placeholder}
          data-slot="input"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={required}
          disabled={disabled}
          step={step}
          {...rest}
          className={cn(inputVariants({ size }), inputClassName)}
        />
        {toggleVisibility && (
          <Tooltip value={inputVisible ? "Hide Password" : "Show Password"}>
            <div
              className="absolute top-[5%] bottom-[5%] right-[0.15rem] z-10 text-sm hover:bg-accent hover:text-accent-foreground flex items-center justify-center rounded-md px-2 has-[>svg]:px-1 scale-90 py-1 h-[95%] select-none"
              onClick={toggle}
              tabIndex={-1}
            >
              <Icon type={inputVisible ? "visibility_off" : "visibility"} />
            </div>
          </Tooltip>
        )}
      </div>
    </div>
  );
}

export { Input };
