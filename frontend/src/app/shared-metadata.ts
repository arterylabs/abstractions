export const siteUrl =
  "https://" + (process.env.NEXT_PUBLIC_SITE_URL ?? "useabstractions.com");

export const ogImage = {
  url: "/og-image.png",
  width: 1200,
  height: 630,
  alt: "Abstractions - Receive weekly newsletters explaining your codebase",
};

export const ogImages = [ogImage];
export const twitterImages = [ogImage.url];
