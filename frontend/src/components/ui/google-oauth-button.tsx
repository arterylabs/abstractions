"use client";

import { Button } from "@/components/ui/button";

interface GoogleOAuthButtonProps {
  onClick: () => void;
  label?: string;
}

export function GoogleOAuthButton({
  onClick,
  label = "Continue with Google",
}: GoogleOAuthButtonProps) {
  return (
    <Button
      variant="outline"
      className="w-full hover:text-[#3c4043] bg-white text-[#3c4043] hover:bg-gray-200 border-[#dadce0] hover:border-[#dadce0]"
      onClick={onClick}
      style={{ fontFamily: "'Roboto', sans-serif" }}
    >
      <svg
        className="h-5 w-5"
        viewBox="0 0 100 100"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <g clipPath="url(#clip0_101_21)">
          <rect width="100" height="100" fill="white" />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M98 51.1364C98 47.5909 97.6818 44.1818 97.0909 40.9091H50V60.25H76.9091C75.75 66.5 72.2273 71.7955 66.9318 75.3409V87.8864H83.0909C92.5455 79.1818 98 66.3637 98 51.1364Z"
            fill="#4285F4"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M50 100C63.5 100 74.8182 95.5227 83.0909 87.8864L66.9318 75.3409C62.4546 78.3409 56.7273 80.1136 50 80.1136C36.9773 80.1136 25.9546 71.3182 22.0228 59.5H5.31821V72.4545C13.5455 88.7955 30.4546 100 50 100Z"
            fill="#34A853"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M22.0227 59.5C21.0227 56.5 20.4545 53.2955 20.4545 50C20.4545 46.7046 21.0227 43.5 22.0227 40.5V27.5455H5.31818C1.93182 34.2955 0 41.9318 0 50C0 58.0682 1.93182 65.7046 5.31818 72.4546L22.0227 59.5Z"
            fill="#FBBC05"
          />
          <path
            fillRule="evenodd"
            clipRule="evenodd"
            d="M50 19.8864C57.3409 19.8864 63.9318 22.4091 69.1137 27.3636L83.4546 13.0227C74.7955 4.95455 63.4773 0 50 0C30.4546 0 13.5455 11.2045 5.31821 27.5455L22.0228 40.5C25.9546 28.6818 36.9773 19.8864 50 19.8864Z"
            fill="#EA4335"
          />
        </g>
        <defs>
          <clipPath id="clip0_101_21">
            <rect width="100" height="100" fill="white" />
          </clipPath>
        </defs>
      </svg>
      {label}
    </Button>
  );
}
