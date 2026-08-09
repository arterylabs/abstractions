"use client";

import { useEffect } from "react";
import Link from "next/link";
import StatsCard from "@/components/admin/StatsCard";
import AdminAnalyticsChart from "@/components/admin/AdminAnalyticsChart";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardDescription,
} from "@/components/ui/card";
import { PageHeader } from "@/components/ui/page-header";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { Alert } from "@/components/ui/alert";
import Icon from "@/components/ui/icon";

const quickActions = [
  {
    href: "/god-mode/users",
    title: "Manage users",
    description: "View, verify, impersonate, and manage accounts",
    icon: "group",
  },
  {
    href: "/god-mode/workspaces",
    title: "Manage workspaces",
    description: "View, edit, and recover workspaces",
    icon: "workspaces",
  },
  {
    href: "/god-mode/configuration",
    title: "Configuration",
    description: "Review non-sensitive environment settings",
    icon: "tune",
  },
] as const;

export default function GodModeDashboard() {
  const [statsResponse, getStats] = useAPI(api.admin.getStats);

  useEffect(() => {
    getStats();
  }, [getStats]);

  if (statsResponse.isFetching) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Spinner />
      </div>
    );
  }

  if (statsResponse.error) {
    return (
      <Alert variant="destructive">
        <p>Failed to load dashboard stats: {statsResponse.error.message}</p>
      </Alert>
    );
  }

  const stats = statsResponse.result;

  return (
    <div className="space-y-8">
      <PageHeader
        title="Dashboard"
        description="Overview of platform activity and key metrics"
      />

      {stats && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
          <StatsCard
            title="Total users"
            value={stats.usersCount}
            icon="group"
          />
          <StatsCard
            title="Active users"
            value={stats.activeUsersCount}
            icon="person"
          />
          <StatsCard
            title="Workspaces"
            value={stats.workspacesCount}
            icon="workspaces"
          />
          <StatsCard
            title="Repositories"
            value={stats.repositoriesCount}
            icon="source"
          />
          <StatsCard
            title="Newsletters"
            value={stats.newslettersCount}
            icon="mail"
          />
        </div>
      )}

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <AdminAnalyticsChart
          metric="users_created"
          title="Users created (last 30 days)"
          color="var(--chart-1)"
        />
        <AdminAnalyticsChart
          metric="workspaces_created"
          title="Workspaces created (last 30 days)"
          color="var(--chart-2)"
        />
        <AdminAnalyticsChart
          metric="repositories_created"
          title="Repositories created (last 30 days)"
          color="var(--chart-3)"
        />
        <AdminAnalyticsChart
          metric="newsletters_created"
          title="Newsletters created (last 30 days)"
          color="var(--chart-4)"
        />
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Quick actions</CardTitle>
          <CardDescription>Jump into common admin tasks</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            {quickActions.map((action) => (
              <Link
                key={action.href}
                href={action.href}
                className="press-card group flex items-start gap-3 overflow-hidden border border-border bg-card p-4 hover:no-underline"
              >
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-sm border border-violet/30 text-violet-light">
                  <Icon type={action.icon} className="text-lg" />
                </div>
                <div className="min-w-0">
                  <p className="font-medium text-foreground group-hover:text-violet-light">
                    {action.title}
                  </p>
                  <p className="mt-0.5 text-sm text-muted-foreground">
                    {action.description}
                  </p>
                </div>
                <Icon
                  type="chevron_right"
                  className="ml-auto mt-1 shrink-0 text-base text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100"
                />
              </Link>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
