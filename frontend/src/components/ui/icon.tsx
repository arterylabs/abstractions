export default function Icon({
  type,
  className,
}: {
  type: string;
  className?: string;
}) {
  return (
    <span
      className={`material-symbols-rounded ${className || ""} cursor-default`}
    >
      {type}
    </span>
  );
}
