"use client";

export function NewsletterPreviewFrame({
  htmlContent,
}: {
  htmlContent: string;
}) {
  return (
    <iframe
      srcDoc={htmlContent}
      sandbox="allow-same-origin"
      title="Newsletter Preview"
      // colorScheme forces light mode inside the iframe regardless of OS preference
      style={{ colorScheme: "light" }}
      className="w-full min-h-[600px] border rounded-lg bg-white"
    />
  );
}
