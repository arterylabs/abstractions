"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

export default function FAQ() {
  const questions = [
    {
      question: "What will it cost me to enrich one column for 1000 records?",
      answer: (
        <span>
          For example, you&lsquo;re enriching a company news column.
          <br />
          <br />
          First, Harmony runs a Google search for the company news which will
          cost us <strong>$0.005</strong> then we visit the top 5 websites to
          extract the news which will cost us nothing. Then we pass the content
          to OpenAI to summarize which will cost around <strong>$0.0007</strong>{" "}
          (varies on configuration but we&lsquo;ll assume 5k tokens on
          gpt-4o-mini as healthy buffer). <br />
          <br />
          For 1000 records, enriching this one column from Google will cost you
          about <strong>$5</strong>.
          <br />
          If you skip the google search and extract it from the website
          directly, it will cost you about <strong>$0.70</strong> (That&lsquo;s
          less than a dollar!).
        </span>
      ),
    },
    {
      question:
        "What format should my data be in, and what does Harmony create?",
      answer:
        "Harmony accepts CSV files, and after enrichment it returns a CSV file with the enriched record. Your original file is not modified and is kept as is, instead it creates a new files with the enriched data.",
    },
    {
      question: "Can I use my license on multiple devices?",
      answer:
        "No, the account and license can be used on only one device. If you need to use it on another device, you can activate it on your new device and it will be automatically deactivated on the old device.",
    },
    {
      question: "Is my data safe, and how do you handle it?",
      answer:
        "We don't store or have access to your data, everything stays on your device. Required details are sent to the API providers you configure, and the response if processed on your device. Your device only talks to Harmony to verify if the license if active, and shares some usage data to improve the product.",
    },
    {
      question: "How do I get a refund?",
      answer:
        "If you are not satisfied with Harmony, you can request a refund within 30 days of purchase. Refunds are processed within 7 business days",
    },
  ];

  return (
    <section
      id="faq"
      className="w-full py-12 md:py-24 lg:py-32 bg-white border-t border-b"
    >
      <div className="container max-w-[1000px] mx-auto">
        <h2 className="text-4xl font-bold mb-12">Frequently Asked Questions</h2>
        <Accordion type="single" collapsible>
          {questions.map((faq, index) => (
            <AccordionItem key={index} value={`item-${index}`}>
              <AccordionTrigger className="text-left text-lg font-medium">
                {faq.question}
              </AccordionTrigger>
              <AccordionContent className="text-muted-foreground">
                {faq.answer}
              </AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      </div>
    </section>
  );
}
