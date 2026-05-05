import Spinner from "@/components/ui/spinner";

export default function Loading() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-slate-50">
      <Spinner className="w-8 h-8" />
    </div>
  );
}
