import { redirect } from "next/navigation";

export default async function RepositoryPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  redirect(`/repositories/${id}/overview`);
}
