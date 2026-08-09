"use client";

import { NewsletterTopic } from "@/api/newsletters";
import { UserAvatar } from "@/components/ui/user-avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { SimpleDropdownMenu } from "@/components/ui/dropdown-menu";
import Icon from "@/components/ui/icon";

interface TopicCardProps {
  topic: NewsletterTopic;
  isFirst: boolean;
  isLast: boolean;
  isSent: boolean;
  onMoveUp: () => void;
  onMoveDown: () => void;
  onMoveToTop: () => void;
  onEdit: () => void;
  onDelete: () => void;
}

export function TopicCard({
  topic,
  isFirst,
  isLast,
  isSent,
  onMoveUp,
  onMoveDown,
  onMoveToTop,
  onEdit,
  onDelete,
}: TopicCardProps) {
  const moreOptions = [];

  if (!isSent) {
    moreOptions.push({
      label: "Edit",
      icon: <Icon type="edit" />,
      onClick: onEdit,
    });
  }

  if (!isSent && !isFirst) {
    moreOptions.push({
      label: "Move to top",
      icon: <Icon type="vertical_align_top" />,
      onClick: onMoveToTop,
    });
  }

  if (!isSent) {
    moreOptions.push({
      label: "Delete",
      icon: <Icon type="delete" />,
      onClick: onDelete,
      variant: "destructive" as const,
    });
  }

  return (
    <Card className="press-card overflow-hidden">
      <CardContent className="p-5">
        <div className="flex items-start gap-3">
          {!isSent && (
            <div className="flex flex-col gap-0.5 shrink-0 mt-0.5">
              <Button
                variant="ghost"
                size="sm"
                className="h-7 w-7 p-0"
                disabled={isFirst}
                onClick={onMoveUp}
                title="Move up"
              >
                <Icon type="arrow_upward" className="!text-base" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="h-7 w-7 p-0"
                disabled={isLast}
                onClick={onMoveDown}
                title="Move down"
              >
                <Icon type="arrow_downward" className="!text-base" />
              </Button>
            </div>
          )}
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <h3 className="text-xl font-medium leading-tight">
                {topic.title}
              </h3>
              {moreOptions.length > 0 && (
                <SimpleDropdownMenu options={moreOptions} widthClassName="w-44">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-7 w-7 p-0 shrink-0"
                  >
                    <Icon type="more_vert" className="!text-base" />
                  </Button>
                </SimpleDropdownMenu>
              )}
            </div>
            {topic.description && (
              <p className="mt-2 line-clamp-2 max-w-[72ch] text-sm leading-relaxed text-muted-foreground">
                {topic.description}
              </p>
            )}
            <div className="mt-4 flex items-center gap-2 border-t border-border pt-3">
              {topic.requestedByUserId ? (
                <>
                  <UserAvatar
                    name={topic.requestedByName ?? ""}
                    imageUrl={topic.requestedByAvatar}
                    className="h-5 w-5"
                    fallbackClassName="text-[10px]"
                  />
                  <span className="press-label text-muted-foreground">
                    {topic.requestedByName ?? "Unknown user"}
                  </span>
                </>
              ) : (
                <>
                  <Icon
                    type="auto_awesome"
                    className="!text-sm text-muted-foreground"
                  />
                  <span className="press-label text-muted-foreground">
                    System generated
                  </span>
                </>
              )}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
