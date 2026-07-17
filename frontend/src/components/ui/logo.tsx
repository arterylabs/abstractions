/* eslint-disable @next/next/no-img-element */
export default function Logo({
  height,
  width,
  icon = false,
  className,
}: {
  height?: number | string;
  width?: number | string;
  icon?: boolean;
  className?: string;
}) {
  if (icon) {
    return (
      <img
        src="/assets/icon.svg"
        height={height ?? 200}
        width={width ?? 200}
        alt="Abstractions Logo"
        className={className || ""}
      />
    );
  }

  return (
    // <div className={`flex items-center justify-center ${className || ""}`}>
    <img
      src="/assets/logo.svg"
      height={height ?? 52}
      width={width ?? 408}
      alt="Abstractions Logo"
      className={className || ""}
    />
    // </div>
  );
}
