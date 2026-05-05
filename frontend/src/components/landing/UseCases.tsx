"use client";

import { useEffect, useState } from "react";
import InfoFlow from "./flow";
import { ChevronRightIcon } from "lucide-react";

export default function UseCases() {
  const useCases = [
    {
      title: "Find Total Funding Raised",
      flow: {
        startText: "Column = Company Funding",
        finalText: "Company Funding = $104 Million",
        steps: [
          {
            image: "/assets/img/openai.svg",
            content: (
              <span>
                Search for{" "}
                <span className="border-primary border text-primary p-1 px-2 rounded text-sm">
                  Company Name
                </span>{" "}
                funding
              </span>
            ),
          },
          {
            image: "/assets/img/chrome.svg",
            content: <span>Visit Top Search Results</span>,
          },
          {
            image: "/assets/img/openai.svg",
            content: <span>Get Total Funding Amount</span>,
          },
        ],
      },
    },
    {
      title: "Find Contact Information",
      flow: {
        startText: "Column = Company Contact",
        finalText: "Company Contact = John Doe",
        steps: [
          {
            image: "/assets/img/openai.svg",
            content: (
              <span>
                Search for{" "}
                <span className="border-primary border text-primary p-1 px-2 rounded text-sm">
                  Company Name
                </span>{" "}
                contact
              </span>
            ),
          },
          {
            image: "/assets/img/chrome.svg",
            content: <span>Visit Top Search Results</span>,
          },
          {
            image: "/assets/img/openai.svg",
            content: <span>Get Contact Information</span>,
          },
        ],
      },
    },
    {
      title: "Find LinkedIn Profiles",
      flow: {
        startText: "Column = Company LinkedIn",
        finalText: "Company LinkedIn = John Doe, john@acme..",
        steps: [
          {
            image: "/assets/img/openai.svg",
            content: (
              <span>
                Search for{" "}
                <span className="border-primary border text-primary p-1 px-2 rounded text-sm mr-1">
                  User
                </span>
                <span className="border-primary border text-primary p-1 px-2 rounded text-sm">
                  Company Name
                </span>{" "}
                LinkedIn
              </span>
            ),
          },
          {
            image: "/assets/img/chrome.svg",
            content: <span>Visit Top Search Results</span>,
          },
          {
            image: "/assets/img/openai.svg",
            content: <span>Get LinkedIn Profile</span>,
          },
        ],
      },
    },
    {
      title: "Enrich Latest Company News",
      flow: {
        startText: "Column = Company News",
        finalText: "Company News = Leadership Change...",
        steps: [
          {
            image: "/assets/img/google.svg",
            content: (
              <span>
                Search for{" "}
                <span className="border-primary border text-primary p-1 px-2 rounded text-sm">
                  Company Name
                </span>{" "}
                News
              </span>
            ),
          },
          {
            image: "/assets/img/chrome.svg",
            content: <span>Visit Sub-pages</span>,
          },
          {
            image: "/assets/img/openai.svg",
            content: <span>Get Company News</span>,
          },
        ],
      },
    },
    {
      title: "Find Employee Count",
      flow: {
        startText: "Column = Company Employees",
        finalText: "Company Employees = 100",
        steps: [
          {
            image: "/assets/img/openai.svg",
            content: (
              <span>
                Search for{" "}
                <span className="border-primary border text-primary p-1 px-2 rounded text-sm">
                  Company Name
                </span>{" "}
                employee count
              </span>
            ),
          },
          {
            image: "/assets/img/chrome.svg",
            content: <span>Visit Top Search Results</span>,
          },
          {
            image: "/assets/img/openai.svg",
            content: <span>Get Employee Count</span>,
          },
        ],
      },
    },
  ];

  const [highlightedIndex, setHighlightedIndex] = useState(0);

  useEffect(() => {
    let timer: string | number | NodeJS.Timeout | undefined;
    if (highlightedIndex < useCases.length) {
      timer = setTimeout(() => {
        setHighlightedIndex((prev) => prev + 1);
      }, 6000); // Adjust delay as needed
    }
    if (highlightedIndex >= useCases.length) {
      setHighlightedIndex(0);
    }
    return () => clearTimeout(timer);
  }, [highlightedIndex, useCases.length]);

  const selectedUseCase = useCases[highlightedIndex];

  return (
    <section
      id="use-cases"
      className="w-full py-12 md:py-24 lg:py-32 bg-white border-b"
    >
      <div className="px-4 md:px-6">
        <div className="px-4 md:px-6 max-w-[1000px] container mx-auto">
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-1">
              <h2 className="text-4xl font-bold mb-6">Use Cases</h2>
              <div className="opacity-60 mb-12">
                Some examples of the kind of columns you can enrich
              </div>
              {useCases.map((useCase, index) => (
                <div key={index} className="py-3 flex items-center">
                  {highlightedIndex === index && (
                    <ChevronRightIcon className="h-6 w-6 mr-2 text-primary" />
                  )}
                  <h3
                    className={`text-lg font-medium ${
                      highlightedIndex === index ? "text-primary" : ""
                    }`}
                  >
                    {useCase.title}
                  </h3>
                </div>
              ))}
              <div className="mb-4 py-3 text-lg font-medium">+ Many More</div>
            </div>
            <div className="col-span-1 flex items-center justify-end">
              {selectedUseCase && (
                <InfoFlow
                  key={highlightedIndex}
                  startText={selectedUseCase.flow.startText}
                  finalText={selectedUseCase.flow.finalText}
                  loop
                  steps={selectedUseCase.flow.steps}
                />
              )}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
