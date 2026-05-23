"use client";

import Link from "next/link";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Icon from "@/components/ui/icon";

interface NewsletterPromptCardProps {
  repositoryId: string;
  newsletterEnabled: boolean;
}

export function NewsletterPromptCard({
  repositoryId,
  newsletterEnabled,
}: NewsletterPromptCardProps) {
  if (!newsletterEnabled) {
    return null;
  }

  return (
    <Card className="border-primary/20 bg-gradient-to-r from-primary/5 to-transparent">
      <CardContent className="flex items-center justify-between gap-4 pt-6">
        <div className="flex items-center gap-4">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
            <Icon type="mail" className="!text-lg text-primary" />
          </div>
          <div className="space-y-1">
            <p className="font-semibold text-foreground">
              Ready to send your next newsletter?
            </p>
            <p className="text-sm text-muted-foreground">
              Generate a draft from your queued topics and send it to your team.
            </p>
          </div>
        </div>
        <Link href={`/repositories/${repositoryId}/newsletters`}>
          <Button className="shrink-0">
            <Icon type="mail_outline" className="!text-base" />
            Generate Draft
          </Button>
        </Link>
      </CardContent>
    </Card>
  );
}
