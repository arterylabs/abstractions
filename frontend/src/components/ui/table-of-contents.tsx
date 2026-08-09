"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";

export interface TocEntry {
  id: string;
  label: string;
  level: 2 | 3;
}

export default function TableOfContents({ toc }: { toc: TocEntry[] }) {
  const [activeId, setActiveId] = useState<string>("");
  const observerRef = useRef<IntersectionObserver | null>(null);

  useEffect(() => {
    if (!toc.length) return;

    const handleIntersect: IntersectionObserverCallback = (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          setActiveId(entry.target.id);
          break;
        }
      }
    };

    observerRef.current = new IntersectionObserver(handleIntersect, {
      rootMargin: "-10% 0px -70% 0px",
      threshold: 0,
    });

    toc.forEach(({ id }) => {
      const el = document.getElementById(id);
      if (el) observerRef.current?.observe(el);
    });

    return () => observerRef.current?.disconnect();
  }, [toc]);

  if (!toc.length) return null;

  return (
    <nav className="text-sm">
      <p className="press-label text-muted-foreground/70 mb-3">On this page</p>
      <ul className="flex flex-col gap-1.5">
        {toc.map(({ id, label, level }) => (
          <li key={id} className={level === 3 ? "pl-3" : ""}>
            <Link
              href={`#${id}`}
              className={`block py-0.5 no-underline transition-colors text-foreground/55 hover:text-foreground leading-snug ${
                activeId === id ? "text-primary font-medium" : ""
              }`}
            >
              {label}
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}
