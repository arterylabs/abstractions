"use client";

import {
  useEffect,
  useRef,
  useState,
  type ComponentPropsWithoutRef,
} from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";
import { Label } from "@/components/ui/label";

const COUNTRY_CODES = [
  { code: "+91", flag: "🇮🇳", country: "India" },
  { code: "+1", flag: "🇺🇸", country: "USA / Canada" },
  { code: "+44", flag: "🇬🇧", country: "UK" },
  { code: "+61", flag: "🇦🇺", country: "Australia" },
  { code: "+64", flag: "🇳🇿", country: "New Zealand" },
  { code: "+971", flag: "🇦🇪", country: "UAE" },
  { code: "+966", flag: "🇸🇦", country: "Saudi Arabia" },
  { code: "+65", flag: "🇸🇬", country: "Singapore" },
  { code: "+60", flag: "🇲🇾", country: "Malaysia" },
  { code: "+880", flag: "🇧🇩", country: "Bangladesh" },
  { code: "+92", flag: "🇵🇰", country: "Pakistan" },
  { code: "+94", flag: "🇱🇰", country: "Sri Lanka" },
  { code: "+977", flag: "🇳🇵", country: "Nepal" },
  { code: "+49", flag: "🇩🇪", country: "Germany" },
  { code: "+33", flag: "🇫🇷", country: "France" },
  { code: "+39", flag: "🇮🇹", country: "Italy" },
  { code: "+34", flag: "🇪🇸", country: "Spain" },
  { code: "+31", flag: "🇳🇱", country: "Netherlands" },
  { code: "+7", flag: "🇷🇺", country: "Russia" },
  { code: "+86", flag: "🇨🇳", country: "China" },
  { code: "+81", flag: "🇯🇵", country: "Japan" },
  { code: "+82", flag: "🇰🇷", country: "South Korea" },
  { code: "+55", flag: "🇧🇷", country: "Brazil" },
  { code: "+52", flag: "🇲🇽", country: "Mexico" },
  { code: "+27", flag: "🇿🇦", country: "South Africa" },
  { code: "+234", flag: "🇳🇬", country: "Nigeria" },
  { code: "+254", flag: "🇰🇪", country: "Kenya" },
  { code: "+20", flag: "🇪🇬", country: "Egypt" },
];

const LOCALE_TO_DIAL_CODE: Record<string, string> = {
  IN: "+91",
  US: "+1",
  CA: "+1",
  GB: "+44",
  AU: "+61",
  NZ: "+64",
  AE: "+971",
  SA: "+966",
  SG: "+65",
  MY: "+60",
  BD: "+880",
  PK: "+92",
  LK: "+94",
  NP: "+977",
  DE: "+49",
  FR: "+33",
  IT: "+39",
  ES: "+34",
  NL: "+31",
  RU: "+7",
  CN: "+86",
  JP: "+81",
  KR: "+82",
  BR: "+55",
  MX: "+52",
  ZA: "+27",
  NG: "+234",
  KE: "+254",
  EG: "+20",
};

function detectCountryCodeFromLocale(): string {
  try {
    const locale = navigator.language || "";
    const parts = locale.split("-");
    if (parts.length >= 2) {
      const country = parts[parts.length - 1].toUpperCase();
      if (LOCALE_TO_DIAL_CODE[country]) return LOCALE_TO_DIAL_CODE[country];
    }
  } catch {
    // ignore
  }
  return "+91";
}

async function detectCountryCodeFromIP(): Promise<string> {
  try {
    const resp = await fetch("https://www.cloudflare.com/cdn-cgi/trace", {
      signal: AbortSignal.timeout(3000),
    });
    const text = await resp.text();
    const loc = text
      .split("\n")
      .find((l) => l.startsWith("loc="))
      ?.split("=")[1]
      ?.trim();
    if (loc && LOCALE_TO_DIAL_CODE[loc]) return LOCALE_TO_DIAL_CODE[loc];
  } catch {
    // fall back to locale
  }
  return detectCountryCodeFromLocale();
}

/** Returns true for a valid local phone number segment (digits only, 4-15 chars). */
function isValidLocalNumber(value: string): boolean {
  const cleaned = value.replace(/[\s\-()\./]/g, "");
  if (!cleaned) return true; // empty is fine (field is optional)
  return /^\d{4,15}$/.test(cleaned);
}

interface PhoneInputProps {
  id: string;
  label?: string;
  /** Receives the full E.164-formatted phone string, e.g. "+919876543210", or empty string. */
  onChange: (value: string) => void;
  className?: string;
  value?: string;
  size?: VariantProps<typeof phoneInputVariants>["size"];
}

type NativeInputProps = Omit<
  ComponentPropsWithoutRef<"input">,
  "id" | "value" | "onChange" | "className" | "size"
>;

const phoneInputVariants = cva(
  "flex w-full overflow-hidden rounded-md border bg-transparent shadow-xs transition-[color,box-shadow]",
  {
    variants: {
      size: {
        xs: "h-7",
        sm: "h-8",
        default: "h-9",
        lg: "h-10",
      },
    },
    defaultVariants: {
      size: "default",
    },
  },
);

const phoneInputDisplayVariants = cva(
  "pointer-events-none flex items-center gap-1",
  {
    variants: {
      size: {
        xs: "px-2 py-1",
        sm: "px-2 py-1",
        default: "px-2 py-1",
        lg: "px-3 py-2",
      },
    },
    defaultVariants: {
      size: "default",
    },
  },
);

const phoneInputFieldVariants = cva(
  "flex-1 min-w-0 bg-transparent placeholder:text-muted-foreground focus:outline-none disabled:cursor-not-allowed disabled:opacity-50",
  {
    variants: {
      size: {
        xs: "px-2 py-1 text-xs",
        sm: "px-3 py-1 text-sm",
        default: "px-3 py-1 text-sm",
        lg: "px-4 py-2 text-base",
      },
    },
    defaultVariants: {
      size: "default",
    },
  },
);

const COUNTRY_CODES_BY_LENGTH = [...COUNTRY_CODES].sort(
  (a, b) => b.code.length - a.code.length,
);

function parsePhoneValue(value: string): {
  dialCode: string;
  localNumber: string;
} {
  const cleaned = value.replace(/[\s\-()\./]/g, "");
  if (!cleaned) {
    return { dialCode: "", localNumber: "" };
  }

  if (!cleaned.startsWith("+")) {
    return { dialCode: "", localNumber: cleaned };
  }

  const matchedCountryCode = COUNTRY_CODES_BY_LENGTH.find(
    ({ code }) => cleaned.startsWith(code) && cleaned.length > code.length,
  );
  if (matchedCountryCode) {
    return {
      dialCode: matchedCountryCode.code,
      localNumber: cleaned.slice(matchedCountryCode.code.length),
    };
  }

  for (
    let prefixLength = Math.min(4, cleaned.length - 4);
    prefixLength >= 2;
    prefixLength -= 1
  ) {
    const dialCode = cleaned.slice(0, prefixLength);
    const localNumber = cleaned.slice(prefixLength);
    if (/^\+\d+$/.test(dialCode) && /^\d+$/.test(localNumber)) {
      return { dialCode, localNumber };
    }
  }

  return { dialCode: "", localNumber: cleaned.slice(1) };
}

/**
 * A phone number input with an integrated country-code selector.
 * Auto-detects the user's country via IP geolocation (Cloudflare) with a
 * browser-locale fallback. Emits the full E.164 number on every change.
 */
export function PhoneInput({
  id,
  label,
  onChange,
  className,
  value,
  onBlur,
  onKeyDown,
  autoFocus,
  disabled,
  size,
  ...rest
}: PhoneInputProps & NativeInputProps) {
  const [dialCode, setDialCode] = useState("+91");
  const [localNumber, setLocalNumber] = useState("");
  const [touched, setTouched] = useState(false);
  const detected = useRef(false);

  // Detect country once on mount
  useEffect(() => {
    if (detected.current) return;
    if (value !== undefined && value !== "") return;
    detected.current = true;
    detectCountryCodeFromIP().then((code) => {
      setDialCode(code);
    });
  }, [value]);

  useEffect(() => {
    if (value === undefined) return;
    if (value === "") {
      setLocalNumber("");
      return;
    }

    const parsed = parsePhoneValue(value);
    if (parsed.dialCode) {
      setDialCode(parsed.dialCode);
    }
    setLocalNumber(parsed.localNumber);
  }, [value]);

  const isInvalid =
    touched && localNumber !== "" && !isValidLocalNumber(localNumber);

  const emit = (code: string, number: string) => {
    const cleaned = number.replace(/[\s\-()\./]/g, "");
    if (!cleaned) {
      onChange("");
      return;
    }
    onChange(`${code}${cleaned}`);
  };

  const handleCodeChange = (code: string) => {
    setDialCode(code);
    emit(code, localNumber);
  };

  const handleNumberChange = (value: string) => {
    setLocalNumber(value);
    emit(dialCode, value);
  };

  const countryCodeOptions = COUNTRY_CODES.some((c) => c.code === dialCode)
    ? COUNTRY_CODES
    : [{ code: dialCode, flag: "🌐", country: "Custom" }, ...COUNTRY_CODES];

  const selectedEntry = countryCodeOptions.find((c) => c.code === dialCode);

  return (
    <div className={cn("grid gap-2", className)}>
      {label && (
        <Label htmlFor={id} className="text-sm leading-none select-none">
          {label}
        </Label>
      )}
      <div
        data-phone-input-wrapper
        className={cn(
          phoneInputVariants({ size }),
          isInvalid
            ? "border-destructive/60 ring-destructive/20 focus-within:ring-[3px] focus-within:outline-none"
            : "border-input focus-within:ring-4 focus-within:outline-1 focus-within:ring-ring/10 focus-within:outline-ring/50",
        )}
      >
        {/* Country code selector */}
        <div className="relative flex shrink-0 items-center border-r border-input bg-muted/40">
          {/* Invisible native select covers the whole button area for interaction */}
          <select
            value={dialCode}
            onChange={(e) => handleCodeChange(e.target.value)}
            aria-label="Country dial code"
            disabled={disabled}
            className="absolute inset-0 h-full w-full cursor-pointer opacity-0 disabled:cursor-not-allowed"
          >
            {countryCodeOptions.map((c) => (
              <option key={c.code + c.country} value={c.code}>
                {c.flag} {c.code} — {c.country}
              </option>
            ))}
          </select>
          {/* Visual display of the selected country */}
          <div className={cn(phoneInputDisplayVariants({ size }))}>
            <span className="text-base leading-none">
              {selectedEntry?.flag}
            </span>
            <span className="text-sm font-medium">{dialCode}</span>
            <svg
              className="ml-0.5 size-3 shrink-0 text-muted-foreground"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fillRule="evenodd"
                d="M5.22 8.22a.75.75 0 0 1 1.06 0L10 11.94l3.72-3.72a.75.75 0 1 1 1.06 1.06l-4.25 4.25a.75.75 0 0 1-1.06 0L5.22 9.28a.75.75 0 0 1 0-1.06Z"
                clipRule="evenodd"
              />
            </svg>
          </div>
        </div>

        {/* Local number input */}
        <input
          id={id}
          type="tel"
          inputMode="numeric"
          placeholder="Phone number"
          value={localNumber}
          onChange={(e) => handleNumberChange(e.target.value)}
          onBlur={(e) => {
            setTouched(true);
            onBlur?.(e);
          }}
          onKeyDown={onKeyDown}
          autoFocus={autoFocus}
          disabled={disabled}
          {...rest}
          className={cn(phoneInputFieldVariants({ size }))}
        />
      </div>
      {isInvalid && (
        <p className="text-xs text-destructive">
          Enter a valid phone number (digits only).
        </p>
      )}
    </div>
  );
}
