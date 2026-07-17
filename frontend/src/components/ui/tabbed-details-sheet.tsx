"use client";

import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";

interface Tab {
  value: string;
  label: string;
  content: React.ReactNode;
  contentClassName?: string;
}

interface TabbedDetailsSheetProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  tabs: Tab[];
  defaultTab?: string;
  className?: string;
}

const tabsGridClass: Record<number, string> = {
  1: "grid-cols-1",
  2: "grid-cols-2",
  3: "grid-cols-3",
  4: "grid-cols-4",
};

export function TabbedDetailsSheet({
  isOpen,
  onClose,
  title,
  tabs,
  defaultTab,
  className,
}: TabbedDetailsSheetProps) {
  const [activeTab, setActiveTab] = useState(defaultTab ?? tabs[0]?.value);
  const gridClass = tabsGridClass[tabs.length] ?? "grid-cols-2";

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent
        className={cn(
          "w-full max-w-[90vw] overflow-y-auto p-0 sm:w-[600px] lg:w-[800px]",
          className,
        )}
        side="right"
        aria-describedby={undefined}
      >
        <SheetHeader className="border-b border-border bg-muted/30 px-6 py-4">
          <SheetTitle className="text-lg font-medium">{title}</SheetTitle>
        </SheetHeader>

        <div className="px-6 py-4">
          <Tabs
            value={activeTab}
            onValueChange={setActiveTab}
            className="w-full"
          >
            <TabsList className={cn("mb-6 grid w-full", gridClass)}>
              {tabs.map((tab) => (
                <TabsTrigger key={tab.value} value={tab.value}>
                  {tab.label}
                </TabsTrigger>
              ))}
            </TabsList>

            {tabs.map((tab) => (
              <TabsContent
                key={tab.value}
                value={tab.value}
                className={cn("mt-0", tab.contentClassName)}
              >
                {tab.content}
              </TabsContent>
            ))}
          </Tabs>
        </div>
      </SheetContent>
    </Sheet>
  );
}
