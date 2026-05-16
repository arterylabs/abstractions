"use client";

import Link from "next/link";
import { type ReactNode, useEffect, useState } from "react";
import { useParams, usePathname } from "next/navigation";
import api from "@/api";
import useAPI from "@/api/hooks";
import { Badge } from "@/components/ui/badge";
import Icon from "@/components/ui/icon";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import Tooltip from "@/components/ui/tooltip";
import { formatDateTime } from "@/lib/format";
import { useAppSelector } from "@/store/hooks";
import { Repository } from "@/types";
import { UpdateDescriptionModal } from "@/components/repositories/UpdateDescriptionModal";

const STATUS_CONFIG: Record<
  Repository["status"],
  {
    label: string;
    variant: "secondary" | "warning" | "success" | "destructive";
    className?: string;
  }
> = {
  pending: { label: "Pending", variant: "secondary" },
  indexing: {
    label: "Indexing…",
    variant: "warning",
    className: "animate-pulse",
  },
  indexed: { label: "Indexed", variant: "success" },
  error: { label: "Error", variant: "destructive" },
};

function RepositoryStatusBadge({
  status,
  indexedAt,
}: {
  status: Repository["status"];
  indexedAt?: string | null;
}) {
  const config = STATUS_CONFIG[status];
  const badge = (
    <Badge variant={config.variant} className={config.className}>
      <Icon type="database" className="text-[10px]" />
      {config.label}
    </Badge>
  );
  if (status === "indexed" && indexedAt) {
    return (
      <Tooltip value={`Last indexed ${formatDateTime(indexedAt)}`}>
        {badge}
      </Tooltip>
    );
  }
  return badge;
}

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
    name: "Topics",
    href: `/repositories/${repoId}/topics`,
    value: "topics",
    icon: "subject",
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
    )?.value ?? "overview"
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
  const [isDescriptionModalOpen, setIsDescriptionModalOpen] = useState(false);

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
      {repository && (
        <UpdateDescriptionModal
          repository={repository}
          open={isDescriptionModalOpen}
          onOpenChange={setIsDescriptionModalOpen}
        />
      )}
      <PageHeader
        title={repository?.fullName ?? "Repository"}
        description={
          repository ? (
            repository.description ? (
              repository.description
            ) : (
              <button
                className="text-muted-foreground/50 hover:text-muted-foreground transition-colors cursor-pointer text-sm"
                onClick={() => setIsDescriptionModalOpen(true)}
              >
                Add description
              </button>
            )
          ) : undefined
        }
      >
        {repository && (
          <div className="flex shrink-0 items-center gap-2">
            <Badge
              variant={repository.newsletterEnabled ? "success" : "secondary"}
            >
              <Icon type="mail" className="text-[10px]" />
              {repository.newsletterEnabled
                ? "Newsletters On"
                : "Newsletters Off"}
            </Badge>
            <RepositoryStatusBadge
              status={repository.status}
              indexedAt={repository.indexedAt}
            />
          </div>
        )}
      </PageHeader>
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
