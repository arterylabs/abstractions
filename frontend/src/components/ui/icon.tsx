// Material Symbols icons are font-based: their visual size is controlled by font-size, NOT width/height.
// Always use text-* Tailwind classes (e.g. text-base, text-[2rem]) to size icons.
// size-* classes (e.g. size-8) only set width/height and have no effect on the rendered icon size.
export default function Icon({
  type,
  className,
}: {
  type: string;
  className?: string;
}) {
  return (
    <span
      className={`material-symbols-rounded icon ${className || ""} cursor-default`}
    >
      {type}
    </span>
  );
}
