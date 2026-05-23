"use client";

import { useEffect } from "react";
import StatsCard from "@/components/admin/StatsCard";
import AdminAnalyticsChart from "@/components/admin/AdminAnalyticsChart";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import api from "@/api";
import useAPI from "@/api/hooks";
import Spinner from "@/components/ui/spinner";
import { Alert } from "@/components/ui/alert";

export default function GodModeDashboard() {
  const [statsResponse, getStats] = useAPI(api.admin.getStats);

  useEffect(() => {
    getStats();
  }, [getStats]);

  if (statsResponse.isFetching) {
    return (
      <div className="flex items-center justify-center h-64">
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
    <div>
      <div className="mb-8">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold">Dashboard</h1>
            <p className="mt-2 text-muted-foreground">
              Overview of platform statistics and key metrics
            </p>
          </div>
        </div>
      </div>

      {stats && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
          <StatsCard title="Total Users" value={stats.usersCount} />
          <StatsCard title="Active Users" value={stats.activeUsersCount} />
          <StatsCard title="Total Workspaces" value={stats.workspacesCount} />
          <StatsCard
            title="Total Repositories"
            value={stats.repositoriesCount}
          />
          <StatsCard title="Total Newsletters" value={stats.newslettersCount} />
        </div>
      )}

      <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
        <AdminAnalyticsChart
          metric="users_created"
          title="Users Created (Last 30 Days)"
          color="#3b82f6"
        />
        <AdminAnalyticsChart
          metric="workspaces_created"
          title="Workspaces Created (Last 30 Days)"
          color="#10b981"
        />
      </div>

      <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
        <AdminAnalyticsChart
          metric="repositories_created"
          title="Repositories Created (Last 30 Days)"
          color="#8b5cf6"
        />
        <AdminAnalyticsChart
          metric="newsletters_created"
          title="Newsletters Created (Last 30 Days)"
          color="#f59e0b"
        />
      </div>

      <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <a
                href="/god-mode/users"
                className="block p-3 rounded-md bg-primary/10 hover:bg-primary/20 transition-colors"
              >
                <div className="flex items-center">
                  <span className="text-2xl mr-3">👥</span>
                  <div>
                    <p className="font-medium text-foreground">Manage Users</p>
                    <p className="text-sm text-primary">
                      View, edit, and manage user accounts
                    </p>
                  </div>
                </div>
              </a>

              <a
                href="/god-mode/workspaces"
                className="block p-3 rounded-md bg-success/20 hover:bg-success/30 transition-colors"
              >
                <div className="flex items-center">
                  <span className="text-2xl mr-3">🏢</span>
                  <div>
                    <p className="font-medium text-foreground">
                      Manage Workspaces
                    </p>
                    <p className="text-sm text-success-dark">
                      View, edit, and manage workspaces
                    </p>
                  </div>
                </div>
              </a>
            </div>
          </CardContent>
        </Card>

        {/* Recent Activity placeholder */}
        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Database</span>
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-success text-success-foreground">
                  ✓ Online
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">API</span>
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-success text-success-foreground">
                  ✓ Online
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">
                  Email Service
                </span>
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-success text-success-foreground">
                  ✓ Online
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
