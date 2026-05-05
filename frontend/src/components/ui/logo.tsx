/* eslint-disable @next/next/no-img-element */
export default function Logo({
  height = "auto",
  width = "auto",
  icon = false,
  className,
}: {
  height?: string;
  width?: string;
  icon?: boolean;
  className?: string;
}) {
  if (icon) {
    return (
      <img
        src="/assets/icon.svg"
        height={height}
        width={width}
        alt="Abstractions Logo"
        className={className || ""}
      />
    );
  }

  return (
    // <div className={`flex items-center justify-center ${className || ""}`}>
    <img
      src="/assets/logo.svg"
      height={height}
      width={width}
      alt="Abstractions Logo"
      className={className || ""}
    />
    // </div>
  );
}
