"use client";

import Link from "next/link";
import { type ReactNode, useEffect } from "react";
import { useParams, usePathname } from "next/navigation";
import api from "@/api";
import useAPI from "@/api/hooks";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAppSelector } from "@/store/hooks";

const getRepositoryTabs = (repoId: string) => [
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

export function getRepositorySectionValue(pathname: string, repoId: string) {
  return (
    getRepositoryTabs(repoId).find(
      (tab) => pathname === tab.href || pathname.startsWith(`${tab.href}/`),
    )?.value ?? "newsletters"
  );
}

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

  useEffect(() => {
    if (!currentWorkspace?.id || !repoId) {
      return;
    }

    void loadRepository(currentWorkspace.id, repoId);
  }, [currentWorkspace?.id, loadRepository, repoId]);

  const repository = repositoryResponse.result ?? repositoryFromStore;
  const tabs = getRepositoryTabs(repoId);
  const currentTab = getRepositorySectionValue(pathname, repoId);

  return (
    <div className="space-y-6">
      <PageHeader
        title={repository?.fullName ?? "Repository"}
        description={
          repository?.description ||
          "Manage newsletters, indexing, and repository settings."
        }
      />
      <div className="flex flex-row items-start gap-6">
        <Tabs value={currentTab} orientation="vertical">
          <TabsList className="h-auto w-44 shrink-0 flex-col justify-start bg-transparent p-0">
            {tabs.map((tab) => (
              <TabsTrigger
                key={tab.value}
                value={tab.value}
                variant="sidebar"
                className="justify-start"
                asChild
              >
                <Link href={tab.href}>
                  <Icon type={tab.icon} className="!text-base" />
                  {tab.name}
                </Link>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
        <div className="min-w-0 flex-1">
          {!repository && repositoryResponse.isFetching ? (
            <div className="flex min-h-[240px] flex-col items-center justify-center gap-4">
              <Spinner />
              <p className="text-sm text-muted-foreground">
                Loading repository...
              </p>
            </div>
          ) : (
            children
          )}
        </div>
      </div>
    </div>
  );
}
