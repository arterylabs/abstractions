"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import React from "react";
import { Card, CardContent } from "../ui/card";
import { cn } from "@/lib/utils";

export interface InfoFlowStep {
  content: React.ReactNode;
  image: string;
}

const InfoFlow = ({
  startText,
  steps,
  finalText,
  loop = false,
}: {
  startText: string;
  steps: InfoFlowStep[];
  finalText: string;
  loop?: boolean;
}) => {
  const [highlightedIndex, setHighlightedIndex] = useState(-1);

  useEffect(() => {
    let timer: NodeJS.Timeout;
    if (highlightedIndex < steps.length) {
      timer = setTimeout(() => {
        setHighlightedIndex((prev) => prev + 1);
      }, 1000);
    }
    if (highlightedIndex === steps.length && loop) {
      timer = setTimeout(() => {
        setHighlightedIndex(0);
      }, 2000);
    }
    return () => clearTimeout(timer);
  }, [highlightedIndex, steps.length, loop]);

  return (
    <div className="text-center">
      <div className="text-sm text-muted-foreground mb-4">{startText}</div>
      <div className="space-y-4">
        {steps.map((step, index) => (
          <React.Fragment key={index}>
            <Card
              className={cn(
                "transition-all duration-300 transform",
                highlightedIndex >= index
                  ? "scale-100 opacity-100 border-primary"
                  : "scale-95 opacity-40",
              )}
            >
              <CardContent className="flex items-center gap-4 p-4">
                <div className="w-6 h-6 relative">
                  <Image
                    src={step.image}
                    alt=""
                    fill
                    className="object-contain"
                  />
                </div>
                <div className="flex-1">{step.content}</div>
              </CardContent>
            </Card>
            {index < steps.length - 1 && (
              <div
                className={cn(
                  "w-1 h-8 mx-auto transition-colors duration-300",
                  highlightedIndex > index ? "bg-primary" : "bg-muted",
                )}
              />
            )}
          </React.Fragment>
        ))}
      </div>
      <div
        className={cn(
          "text-sm font-medium mt-4 transition-opacity duration-300",
          highlightedIndex === steps.length ? "opacity-100" : "opacity-0",
        )}
      >
        {finalText}
      </div>
    </div>
  );
};

export default InfoFlow;
