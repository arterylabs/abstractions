"use client";

import { Toaster as Sonner, ToasterProps, toast } from "sonner";
import Icon from "@/components/ui/icon";
import { Button } from "@/components/ui/button";

type NotificationAction = {
  label: string;
  onClick: () => void;
};

export type NotificationParams = {
  title: string;
  message?: string;
  actions?: NotificationAction[];
};

const Toaster = ({ ...props }: ToasterProps) => {
  return (
    <Sonner
      className="toaster group"
      toastOptions={{
        classNames: {
          toast:
            "group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:shadow-lg",
          description: "group-[.toast]:text-muted-foreground",
          actionButton:
            "group-[.toast]:bg-primary group-[.toast]:text-primary-foreground",
          cancelButton:
            "group-[.toast]:bg-muted group-[.toast]:text-muted-foreground",
        },
      }}
      {...props}
    />
  );
};

const getParams = (params: NotificationParams | string) => {
  if (typeof params === "string") {
    params = { title: params };
  }
  return params;
};

const Notifier = {
  success: (params: NotificationParams | string) => {
    toast(<Notification type="success" params={params} />);
  },
  error: (params: NotificationParams | string) => {
    toast(<Notification type="error" params={params} />);
  },
  info: (params: NotificationParams | string) => {
    toast(<Notification type="info" params={params} />);
  },
  warning: (params: NotificationParams | string) => {
    toast(<Notification type="warning" params={params} />);
  },
};

export function Notification({
  type,
  params,
}: {
  type: "success" | "error" | "info" | "warning";
  params: NotificationParams | string;
}) {
  params = getParams(params);

  return (
    <div className="flex items-center gap-4">
      <div className="mt-1 !text-lg">
        {type === "success" && (
          <Icon type="check_circle" className="text-success" />
        )}

        {type === "error" && (
          <Icon type="error" className=" text-destructive" />
        )}

        {type === "info" && <Icon type="info" className=" text-primary" />}

        {type === "warning" && (
          <Icon type="warning" className=" text-amber-500" />
        )}
      </div>
      <div className="flex flex-col h-full">
        <div className="font-medium">{params.title}</div>
        {params.message && (
          <div className="text-xs mt-1 opacity-60">{params.message}</div>
        )}
        {params.actions && (
          <div className="flex gap-2 mt-4">
            {params.actions.map((action) => (
              <Button
                key={action.label}
                variant="ghost"
                size="sm"
                onClick={action.onClick}
              >
                {action.label}
              </Button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export { Toaster, toast, Notifier };
