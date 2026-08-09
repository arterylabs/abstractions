import type { Metadata } from "next";
import { PageHeader } from "@/components/ui/page-header";

export const metadata: Metadata = {
  title: "Cookies Policy",
  description:
    "Read the Abstractions Cookies Policy to understand how we use cookies and similar tracking technologies.",
  alternates: { canonical: "/cookies" },
};

export default function CookiesPolicy() {
  return (
    <>
      <div className="mx-auto max-w-[1200px] w-full docs-content py-12 px-8">
        <PageHeader
          display
          title="Cookies Policy"
          description={
            <>
              Effective date: May 2026 &nbsp;·&nbsp; Artery Ventures, LLP
              &nbsp;·&nbsp;{" "}
              <a
                className="tou-email"
                href="mailto:support@useabstractions.com"
              >
                support@useabstractions.com
              </a>
            </>
          }
        />
        <div>
          <h2>Introduction</h2>
          <p>
            We, Artery Ventures, LLP (&quot;Company&quot;, &quot;we&quot;,
            &quot;us&quot;, &quot;our&quot;) may use cookies or similar
            technologies to understand how visitors use our website{" "}
            <a href="https://useabstractions.com">
              https://useabstractions.com
            </a>{" "}
            and our Abstractions application (together, the
            &quot;Platform&quot;). We are committed to protecting the privacy
            and security of your personal information. We advise you to
            carefully read this cookie policy (&quot;Policy&quot;), together
            with the Company&apos;s Privacy Policy at{" "}
            <a href="https://useabstractions.com/privacy">
              useabstractions.com/privacy
            </a>{" "}
            so that you are aware of the cookies and technologies used as well
            as how we treat your personal information.
          </p>

          <h2>What are Cookies?</h2>
          <p>
            Cookies are small text files which are placed on your device when
            you visit a website. They can allow a site to recognize your device,
            store preferences, analyze trends, operate and improve services, and
            provide a better experience. Similar technologies include local
            storage and other browser storage the Platform or scripts may use
            for the same purposes.
          </p>

          <h2>Types of Cookies Used</h2>

          <h3>By Duration</h3>
          <p>
            <strong>Permanent (persistent) cookies</strong> — Remain on your
            device for a pre-defined period or until you delete them. They can
            help recognize returning visitors and keep you signed in where
            appropriate.
          </p>
          <p>
            <strong>Session cookies</strong> — Last only for the browsing
            session and are typically erased when you close the browser.
          </p>

          <h3>By Domain</h3>
          <p>
            <strong>First-party cookies:</strong> These cookies are set by the
            Platform you are visiting (for example useabstractions.com) and
            share the same domain.
          </p>
          <p>
            <strong>Third-party cookies:</strong> These cookies are placed by a
            domain other than the Platform you are visiting. Where we use
            analytics, we aim to do so in a manner that limits unnecessary
            third-party tracking.
          </p>

          <h2>How we use cookies on the Platform</h2>
          <p>We aim to keep tracking minimal. The Platform may use:</p>

          <h3>Essential / functional</h3>
          <p>
            Cookies or similar storage that are necessary for the Platform to
            function securely and correctly (for example authentication,
            security, load balancing, or basic preferences). These are typically
            first-party and are not used for advertising.
          </p>

          <h3>Website analytics (PostHog)</h3>
          <p>
            On the hosted product, we may use a third-party analytics service to
            understand how the Platform is used so we can improve content,
            layout, and product experience. Analytics data is handled by our
            analytics provider, <a href="https://posthog.com">PostHog</a>. We
            may collect information about page visits and, where enabled,
            selected interactions with the Platform.
          </p>
          <p>
            Analytics requests are sent through a first-party path on our origin
            so they are not loaded from the vendor hostname directly. We do not
            enable session recording in the current integration. Self-hosted
            deployments of Abstractions may leave analytics disabled entirely.
          </p>
          <p>
            Analytics tools may use cookies or similar browser storage to
            distinguish visits in an aggregated way. Ordinary Website visitors
            are not identified to us as named user accounts solely through this
            analytics. We do not currently rely on advertising pixels as a core
            part of the product. If we add or replace tracking tools, we will
            update this Policy.
          </p>

          <h2>Use of information from cookies</h2>
          <p>
            Information from cookies and similar technologies may be used to:
          </p>
          <ul>
            <li>operate and secure the Platform;</li>
            <li>authenticate users and protect accounts;</li>
            <li>understand how visitors navigate the Platform;</li>
            <li>improve content, performance, and layout; and</li>
            <li>
              measure whether product surfaces and calls to action are useful.
            </li>
          </ul>
          <p>
            We do not sell personal information collected via cookies. We do not
            use cookies on the Platform to personalize third-party advertising
            as part of the current Abstractions product experience.
          </p>

          <h2>Disabling Cookies</h2>
          <p>
            You can decide whether to accept cookies. Most browsers let you
            refuse or delete cookies via settings (see your browser&apos;s help
            documentation). Clearing site data may also remove local storage
            used by analytics. Disabling cookies or blocking scripts may affect
            Platform functionality or prevent analytics from loading. You may
            also use browser privacy features or extensions that limit tracking.
          </p>

          <h2>Changes to this Policy</h2>
          <p>
            Please revisit this page periodically to stay aware of any changes
            to this Policy, which we may update from time to time. If we modify
            this Policy, we will make it available through the website and
            indicate the date of the latest revision.
          </p>

          <h2>Contact Us</h2>
          <p>
            If you have any questions or concerns regarding this Policy, you can
            contact us at{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            .
          </p>
        </div>
      </div>
    </>
  );
}
