"use client";

import { useEffect, useMemo } from "react";
import api from "@/api";
import useAPI from "@/api/hooks";
import { PageHeader } from "@/components/ui/page-header";
import Spinner from "@/components/ui/spinner";
import { Alert } from "@/components/ui/alert";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";
import { copyToClipboard } from "@/lib/clipboard";
import type { AdminConfigItem, AdminSecretStatus } from "@/types/admin";

function groupByCategory<T extends { category: string }>(
  items: T[],
): Record<string, T[]> {
  return items.reduce<Record<string, T[]>>((acc, item) => {
    const list = acc[item.category] ?? [];
    list.push(item);
    acc[item.category] = list;
    return acc;
  }, {});
}

function ConfigValueRow({ item }: { item: AdminConfigItem }) {
  const empty = !item.value;

  return (
    <div className="flex flex-col gap-1 border-b border-border py-3 last:border-b-0 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
      <code className="shrink-0 text-sm font-medium text-foreground">
        {item.key}
      </code>
      <div className="flex min-w-0 items-center gap-2 sm:max-w-[60%] sm:justify-end">
        {empty ? (
          <span className="text-sm text-muted-foreground italic">Not set</span>
        ) : (
          <>
            <code className="truncate rounded-md bg-muted px-2 py-1 text-sm text-muted-foreground">
              {item.value}
            </code>
            <Button
              variant="ghost"
              size="sm"
              className="h-8 w-8 shrink-0 p-0"
              onClick={() => copyToClipboard(item.value, item.key)}
              aria-label={`Copy ${item.key}`}
            >
              <Icon type="content_copy" className="text-sm" />
            </Button>
          </>
        )}
      </div>
    </div>
  );
}

function SecretStatusRow({ item }: { item: AdminSecretStatus }) {
  return (
    <div className="flex items-center justify-between gap-4 border-b border-border py-3 last:border-b-0">
      <code className="text-sm font-medium text-foreground">{item.key}</code>
      <Badge variant={item.configured ? "success" : "warning"} size="sm">
        {item.configured ? "Configured" : "Not set"}
      </Badge>
    </div>
  );
}

export default function ConfigurationPage() {
  const [configResponse, getConfig] = useAPI(api.admin.getConfig);

  useEffect(() => {
    getConfig();
  }, [getConfig]);

  const configByCategory = useMemo(
    () => groupByCategory(configResponse.result?.config ?? []),
    [configResponse.result?.config],
  );

  const secretsByCategory = useMemo(
    () => groupByCategory(configResponse.result?.secrets ?? []),
    [configResponse.result?.secrets],
  );

  if (configResponse.isFetching && !configResponse.result) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Spinner />
      </div>
    );
  }

  if (configResponse.error) {
    return (
      <Alert variant="destructive">
        <p>Failed to load configuration: {configResponse.error.message}</p>
      </Alert>
    );
  }

  const configCategories = Object.keys(configByCategory);
  const secretCategories = Object.keys(secretsByCategory);

  return (
    <div className="space-y-8">
      <PageHeader
        title="Configuration"
        description="Non-sensitive environment values currently loaded by the API. Secret values are never shown — only whether they are set."
      >
        <Button
          variant="outline"
          size="sm"
          onClick={() => getConfig()}
          disabled={configResponse.isFetching}
        >
          <Icon type="refresh" className="text-sm" />
          Refresh
        </Button>
      </PageHeader>

      <section className="space-y-4">
        <div>
          <h2 className="text-base font-medium">Public configuration</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Values safe to display. Empty means the variable is unset.
          </p>
        </div>
        <div className="grid grid-cols-1 gap-4 xl:grid-cols-2">
          {configCategories.map((category) => (
            <Card key={category}>
              <CardHeader className="pb-2">
                <CardTitle className="text-base font-medium">
                  {category}
                </CardTitle>
                <CardDescription>
                  {configByCategory[category].length} setting
                  {configByCategory[category].length === 1 ? "" : "s"}
                </CardDescription>
              </CardHeader>
              <CardContent className="pt-0">
                {configByCategory[category].map((item) => (
                  <ConfigValueRow key={item.key} item={item} />
                ))}
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="space-y-4">
        <div>
          <h2 className="text-base font-medium">Secrets</h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Sensitive keys are never returned. This only shows whether each one
            is configured on the server.
          </p>
        </div>
        <div className="grid grid-cols-1 gap-4 xl:grid-cols-2">
          {secretCategories.map((category) => (
            <Card key={category}>
              <CardHeader className="pb-2">
                <CardTitle className="text-base font-medium">
                  {category}
                </CardTitle>
                <CardDescription>
                  {
                    secretsByCategory[category].filter((s) => s.configured)
                      .length
                  }{" "}
                  of {secretsByCategory[category].length} configured
                </CardDescription>
              </CardHeader>
              <CardContent className="pt-0">
                {secretsByCategory[category].map((item) => (
                  <SecretStatusRow key={item.key} item={item} />
                ))}
              </CardContent>
            </Card>
          ))}
        </div>
      </section>
    </div>
  );
}
