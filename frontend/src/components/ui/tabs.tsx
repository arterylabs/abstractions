"use client";

import * as React from "react";
import * as TabsPrimitive from "@radix-ui/react-tabs";

import { cn } from "@/lib/utils";

function Tabs({
  className,
  ...props
}: React.ComponentProps<typeof TabsPrimitive.Root>) {
  return (
    <TabsPrimitive.Root
      data-slot="tabs"
      className={cn("flex flex-col gap-2", className)}
      {...props}
    />
  );
}

function TabsList({
  className,
  ...props
}: React.ComponentProps<typeof TabsPrimitive.List>) {
  return (
    <TabsPrimitive.List
      data-slot="tabs-list"
      className={cn(
        // Tabs sit on a single rule, like running heads across a spread
        "inline-flex w-full items-center justify-start gap-8 overflow-x-auto border-b border-border scrollbar-none",
        className,
      )}
      {...props}
    />
  );
}

type TabsTriggerProps = React.ComponentProps<typeof TabsPrimitive.Trigger> & {
  /** "default" runs the tab along a rule; "sidebar" stacks it in a nav column. */
  variant?: "default" | "sidebar";
};

const triggerBase =
  "ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 inline-flex cursor-pointer items-center gap-2 whitespace-nowrap transition-colors duration-200 focus-visible:ring-4 focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 aria-invalid:focus-visible:ring-0 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4";

function TabsTrigger({
  className,
  variant = "default",
  ...props
}: TabsTriggerProps) {
  return (
    <TabsPrimitive.Trigger
      data-slot="tabs-trigger"
      className={cn(
        triggerBase,
        variant === "default" &&
          // A running head that inks in its rule when selected. -1px pulls the
          // rule down onto the TabsList border so the two read as one line.
          "press-label relative -mb-px justify-center border-b-2 border-transparent px-0 py-3.5 text-muted-foreground hover:text-foreground data-[state=active]:border-primary data-[state=active]:text-foreground",
        variant === "sidebar" &&
          "w-full justify-start border-l border-transparent py-1.5 pl-3 text-sm text-muted-foreground hover:border-border hover:text-foreground data-[state=active]:border-primary data-[state=active]:text-foreground",
        className,
      )}
      {...props}
    />
  );
}

function TabsContent({
  className,
  ...props
}: React.ComponentProps<typeof TabsPrimitive.Content>) {
  return (
    <TabsPrimitive.Content
      data-slot="tabs-content"
      className={cn(
        "ring-ring/10 dark:ring-ring/20 dark:outline-ring/40 outline-ring/50 flex-1 transition-[box-shadow] focus-visible:ring-4 focus-visible:outline-1 aria-invalid:focus-visible:ring-0",
        className,
      )}
      {...props}
    />
  );
}

export { Tabs, TabsList, TabsTrigger, TabsContent };
