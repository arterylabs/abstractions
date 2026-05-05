This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

Smart enrich is a desktop app, I want to design the architecture of the app.

Here are the primary components.

1. The API - backend which takes care of authentication and licenses
2. The Desktop app - the actual desktop app built in electron
3. The web app - gives a UI for license management, onboarding, and documentation

Here's what the user flow for the desktop app looks like,

Login -> Add your OpenAI, or Google Keys -> Upload a CSV File -> Configure Enriched Columns -> Run an Enrichment Job -> Get Enriched CSV File

All of these things will be configurable. We should be able to support more AI providers, search providers in the future, we also want to allow people to make their own API calls in the application. We should also be able to support excel files for imports. The API calls we make could also be proxied through a vpn provider like Brightdata, or in future allow users to connect to clearbit or apollo.io for enrichment

Design an architecture of the application optimizing for configurability and flexibility of the entire system. Focus on the desktop app as that contains the core business logic, the API and the web app are straightforward.
