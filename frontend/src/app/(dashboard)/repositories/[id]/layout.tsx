"use client";

import Link from "next/link";
import { type ReactNode, useEffect, useState } from "react";
import { useParams, usePathname } from "next/navigation";
import api from "@/api";
import useAPI from "@/api/hooks";
import Icon from "@/components/ui/icon";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAppSelector } from "@/store/hooks";
import { UpdateDescriptionModal } from "@/components/repositories/UpdateDescriptionModal";
import { RepositoryMasthead } from "@/components/repositories/RepositoryMasthead";
import { RepositoryMetaColumn } from "@/components/repositories/RepositoryMetaColumn";
import { RepositoryIndexingStatusBanner } from "@/components/repositories/RepositoryIndexingStatusBanner";
import { getRepositorySectionValue } from "./utils";
import Spinner from "@/components/ui/spinner";

/* The sections of the repository, in the order the work moves through them:
   read the state of things, run the newsletter, then change how it runs. */
const getRepositoryTabs = (repoId: string) => [
  {
    name: "Overview",
    href: `/repositories/${repoId}/overview`,
    value: "overview",
    icon: "space_dashboard",
  },
  {
    name: "Newsletters",
    href: `/repositories/${repoId}/newsletters`,
    value: "newsletters",
    icon: "mail",
  },
  {
    name: "Settings",
    href: `/repositories/${repoId}/settings`,
    value: "settings",
    icon: "settings",
  },
];

export default function RepositoryLayout({
  children,
}: {
  children: ReactNode;
}) {
  const pathname = usePathname();
  const params = useParams<{ id: string }>();
  const repoId = params.id;
  const currentWorkspace = useAppSelector(
    (state) => state.auth.currentWorkspace,
  );
  const repositoryFromStore = useAppSelector((state) =>
    state.repositories.items.find((repository) => repository.id === repoId),
  );
  const [repositoryResponse, loadRepository] = useAPI(api.repositories.get);
  const [isDescriptionModalOpen, setIsDescriptionModalOpen] = useState(false);

  useEffect(() => {
    if (!currentWorkspace?.id || !repoId) {
      return;
    }

    void loadRepository(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadRepository, repoId]);

  const repository = repositoryFromStore ?? repositoryResponse.result;
  const tabs = getRepositoryTabs(repoId);
  const currentTab = getRepositorySectionValue(pathname, repoId);

  // Everything below the masthead reads from the repository, so the whole
  // section waits on it rather than each tab handling the empty case itself.
  if (!repository) {
    return (
      <div className="flex min-h-[240px] flex-col items-center justify-center gap-4 py-16">
        {repositoryResponse.isFetching ? (
          <Spinner />
        ) : (
          <p className="text-sm text-muted-foreground">
            This repository could not be loaded.
          </p>
        )}
      </div>
    );
  }

  const workspaceId = currentWorkspace?.id ?? repository.workspaceId;

  return (
    <div className="flex min-h-full flex-col gap-8 md:flex-row md:items-stretch md:gap-0">
      <UpdateDescriptionModal
        repository={repository}
        open={isDescriptionModalOpen}
        onOpenChange={setIsDescriptionModalOpen}
      />

      <div className="flex min-w-0 flex-1 flex-col gap-8 md:pr-10">
        <RepositoryMasthead
          repository={repository}
          workspaceId={workspaceId}
          onEditDescription={() => setIsDescriptionModalOpen(true)}
        />

        {/* On small screens the meta rail sits under the header rather than
            after the full page body. Desktop uses the sticky aside below. */}
        <div className="border-t border-border pt-6 md:hidden">
          <RepositoryMetaColumn repository={repository} />
        </div>

        <Tabs value={currentTab}>
          <TabsList>
            {tabs.map((tab) => (
              <TabsTrigger key={tab.value} value={tab.value} asChild>
                <Link href={tab.href}>
                  <Icon type={tab.icon} className="!text-sm" />
                  {tab.name}
                </Link>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>

        <RepositoryIndexingStatusBanner
          workspaceId={workspaceId}
          repository={repository}
        />

        {children}
      </div>

      {/* Full-height right rail: status + standing facts. Sticky so it stays
          in view while the main column scrolls. Hidden on small screens. */}
      <aside
        className="hidden shrink-0 flex-col md:sticky md:top-0 md:flex md:min-h-[calc(100svh-5rem)] md:w-56 md:self-start md:border-l md:border-border md:pl-8 lg:w-64"
        aria-label="Repository details"
      >
        <RepositoryMetaColumn repository={repository} />
      </aside>
    </div>
  );
}
