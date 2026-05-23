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

export function TabbedDetailsSheet({
  isOpen,
  onClose,
  title,
  tabs,
  defaultTab,
  className,
}: TabbedDetailsSheetProps) {
  const [activeTab, setActiveTab] = useState(defaultTab ?? tabs[0]?.value);

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent
        className={cn(
          "w-full sm:w-[600px] lg:w-[800px] max-w-[90vw] overflow-y-auto p-0",
          className,
        )}
        side="right"
        aria-describedby={undefined}
      >
        <SheetHeader className="px-6 py-4 border-b bg-gray-50/50">
          <SheetTitle className="text-xl">{title}</SheetTitle>
        </SheetHeader>

        <div className="px-6 py-4">
          <Tabs
            value={activeTab}
            onValueChange={setActiveTab}
            className="w-full"
          >
            <TabsList className={`grid w-full grid-cols-${tabs.length} mb-6`}>
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
