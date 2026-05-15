import { redirect } from "next/navigation";

export default function RepositoryPage({ params }: { params: { id: string } }) {
  redirect(`/repositories/${params.id}/overview`);
}
