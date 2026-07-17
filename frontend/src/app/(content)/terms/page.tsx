import type { Metadata } from "next";
import { PageHeader } from "@/components/ui/page-header";

export const metadata: Metadata = {
  title: "Terms of Service",
  description:
    "Read the Abstractions Terms of Service to understand your rights and obligations when using our platform.",
  alternates: { canonical: "/terms" },
};

export default function Terms() {
  return (
    <>
      <div className="mx-auto max-w-[1200px] w-full docs-content py-12 px-8">
        <PageHeader
          title="Terms"
          description={
            <>
              {" "}
              Effective date: June 2025 &nbsp;·&nbsp; Artery Ventures, LLP
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
          <h2>1. Agreement to terms</h2>
          <p>
            These Terms of Use constitute a legally binding agreement made
            between you, whether personally or on behalf of an entity
            (&quot;you&quot;), and Artery Ventures, LLP (&quot;Company&quot;,
            &quot;we&quot;, &quot;us&quot;, or &quot;our&quot;), concerning your
            access to and use of the{" "}
            <a href="https://useabstractions.com">
              https://useabstractions.com
            </a>{" "}
            website and the Abstractions application, as well as any other media
            form, media channel, mobile website, or mobile application related,
            linked, or otherwise connected thereto (collectively, the
            &quot;Site&quot; or &quot;Service&quot;).
          </p>
          <p>
            You agree that by accessing the Site, you have read, understood, and
            agree to be bound by all of these Terms of Use. IF YOU DO NOT AGREE
            WITH ALL OF THESE TERMS OF USE, THEN YOU ARE EXPRESSLY PROHIBITED
            FROM USING THE SITE AND YOU MUST DISCONTINUE USE IMMEDIATELY.
          </p>
          <p>
            Supplemental terms and conditions or documents that may be posted on
            the Site from time to time are hereby expressly incorporated herein
            by reference. We reserve the right, in our sole discretion, to make
            changes or modifications to these Terms of Use at any time and for
            any reason. We will alert you about any changes by updating the
            &quot;Last updated&quot; date of these Terms of Use. Your continued
            use of the Site after any revised Terms of Use are posted
            constitutes your acceptance of those changes.
          </p>
          <p>
            The Site is intended for users who are at least 18 years old.
            Persons under the age of 18 are not permitted to use or register for
            the Site.
          </p>
        </div>

        <div>
          <h2>2. Description of service</h2>
          <p>
            Abstractions is a software-as-a-service platform that allows users
            to connect their code repositories from third-party code hosting
            platforms, index their codebase, and receive periodic newsletters
            and insights about the concepts, patterns, and technical content
            within their codebase. The Service may include, without limitation,
            codebase analysis, automated newsletter generation, and related
            developer tools.
          </p>
          <p>
            We reserve the right to modify, suspend, or discontinue the Service
            (or any part thereof) at any time with or without notice. We shall
            not be liable to you or any third party for any modification,
            suspension, or discontinuation of the Service.
          </p>
        </div>

        <div>
          <h2>3. Intellectual property rights</h2>
          <p>
            Unless otherwise indicated, the Site is our proprietary property and
            all source code, databases, functionality, software, website
            designs, audio, video, text, photographs, and graphics on the Site
            (collectively, the &quot;Content&quot;) and the trademarks, service
            marks, and logos contained therein (the &quot;Marks&quot;) are owned
            or controlled by us or licensed to us, and are protected under
            applicable intellectual property laws of India and international
            conventions.
          </p>
          <p>
            The Content and the Marks are provided on the Site &quot;AS IS&quot;
            for your information and personal use only. Except as expressly
            provided in these Terms of Use, no part of the Site and no Content
            or Marks may be copied, reproduced, aggregated, republished,
            uploaded, posted, publicly displayed, encoded, translated,
            transmitted, distributed, sold, licensed, or otherwise exploited for
            any commercial purpose whatsoever, without our express prior written
            permission.
          </p>
          <p>
            Provided that you are eligible to use the Site, you are granted a
            limited, non-exclusive, non-transferable license to access and use
            the Site solely for your personal or internal business use. We
            reserve all rights not expressly granted to you.
          </p>
          <p>
            Separately from these Terms, the Abstractions application source
            code is made available as free software under the{" "}
            <a
              href="https://github.com/arterylabs/abstractions/blob/master/LICENSE"
              target="_blank"
              rel="noopener noreferrer"
            >
              GNU Affero General Public License v3.0
            </a>
            . Rights to use, modify, and redistribute that source code are
            governed by the AGPL-3.0 license (and any additional notices in the
            repository), not by the limited Site license in this section.
            Trademarks and hosted branding for useabstractions.com are not
            licensed under the AGPL.
          </p>
        </div>

        <div>
          <h2>4. User representations</h2>
          <p>By using the Site, you represent and warrant that:</p>
          <ol>
            <li>
              All registration information you submit will be true, accurate,
              current, and complete.
            </li>
            <li>
              You will maintain the accuracy of such information and promptly
              update your registration information as necessary.
            </li>
            <li>
              You have the legal capacity and agree to comply with these Terms
              of Use.
            </li>
            <li>
              You are not a minor in the jurisdiction in which you reside.
            </li>
            <li>
              You will not access the Site through automated or non-human means,
              whether through a bot, script, or otherwise, except as expressly
              permitted by the Service.
            </li>
            <li>
              You will not use the Site for any illegal or unauthorized purpose.
            </li>
            <li>
              Your use of the Site will not violate any applicable law or
              regulation.
            </li>
          </ol>
          <p>
            If you provide any information that is untrue, inaccurate, not
            current, or incomplete, we have the right to suspend or terminate
            your account and refuse any and all current or future use of the
            Site.
          </p>
        </div>

        <div>
          <h2>5. User registration</h2>
          <p>
            You may be required to register with the Site to access certain
            features, including connecting your code repositories from supported
            code hosting platforms. You agree to keep your credentials
            confidential and will be responsible for all use of your account. We
            reserve the right to remove, reclaim, or change a username you
            select if we determine, in our sole discretion, that such username
            is inappropriate, obscene, or otherwise objectionable.
          </p>
        </div>

        <div>
          <h2>6. Code hosting platform integration and codebase access</h2>
          <p>
            The Service requires you to connect your code hosting platform
            account and authorize Abstractions to access certain repositories.
            By connecting your code hosting platform account, you:
          </p>
          <ol>
            <li>
              Authorize us to access, read, and index the repositories you
              explicitly select or authorize.
            </li>
            <li>
              Represent and warrant that you have the right to grant such
              access, including for repositories belonging to organizations or
              third parties.
            </li>
            <li>
              Acknowledge that we will store indexed data derived from your
              codebase solely to provide and improve the Service.
            </li>
            <li>
              Understand that we do not claim ownership over your code or any
              intellectual property contained therein.
            </li>
          </ol>
          <p>
            You may revoke our access to your code hosting platform account at
            any time through your code hosting platform account settings or
            through the Site. Upon revocation, we will cease indexing your
            repositories, though cached or derived data may be retained for a
            limited period as described in our Privacy Policy.
          </p>
          <p>
            You are solely responsible for ensuring that any third-party code,
            libraries, or proprietary information in your connected repositories
            does not violate the rights of any third party by being shared with
            us. We are not responsible for any unauthorized disclosure of
            third-party intellectual property contained in your repositories.
          </p>
        </div>

        <div>
          <h2>7. Fees and payment</h2>
          <p>
            We accept the following forms of payment: Visa, Mastercard, American
            Express, and other major payment methods as indicated on the Site.
            You may be required to purchase or pay a fee to access certain
            features or tiers of the Service.
          </p>
          <p>
            You agree to provide current, complete, and accurate purchase and
            account information for all purchases made via the Site. We bill you
            through an online billing account for purchases made via the Site.
            All prices are displayed in U.S. Dollars or Indian Rupees as
            applicable, and we may change prices at any time.
          </p>
          <p>
            If your purchase is subject to recurring charges, then you consent
            to our charging your payment method on a recurring basis without
            requiring your prior approval for each recurring charge, until you
            notify us of your cancellation. We reserve the right to correct any
            errors in pricing and to refuse any order placed through the Site.
          </p>
        </div>

        <div>
          <h2>8. Free trial</h2>
          <p>
            We may offer a free trial period to new users who register with the
            Site. The account will not be charged during the free trial period.
            At the end of the free trial, the account will be suspended unless
            upgraded to a paid plan. Trial terms, including duration, will be
            specified at the time of registration.
          </p>
        </div>

        <div>
          <h2>9. Cancellation</h2>
          <p>
            You can cancel your subscription at any time by logging into your
            account settings or by contacting us at{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            . Your cancellation will take effect at the end of the current paid
            billing term, and you will continue to have access to the Service
            until that date.
          </p>
          <p>
            If you are unsatisfied with our services, please contact us at{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            .
          </p>
        </div>

        <div>
          <h2>10. Prohibited activities</h2>
          <p>
            You may not access or use the Site for any purpose other than that
            for which we make the Site available. As a user of the Site, you
            agree not to:
          </p>
          <ol>
            <li>
              Use the Service to access, index, or process repositories for
              which you do not have authorization.
            </li>
            <li>
              Attempt to reverse engineer, decompile, disassemble, or otherwise
              derive the source code of the Abstractions platform or its
              underlying algorithms.
            </li>
            <li>
              Use the Service to generate newsletters, content, or insights for
              the purpose of competing with Abstractions or any of its
              affiliates.
            </li>
            <li>
              Systematically retrieve data or other content from the Site to
              create or compile a collection, compilation, database, or
              directory without written permission from us.
            </li>
            <li>
              Make any unauthorized use of the Site, including collecting
              usernames or email addresses by electronic or other means for the
              purpose of sending unsolicited email.
            </li>
            <li>
              Circumvent, disable, or otherwise interfere with security-related
              features of the Site.
            </li>
            <li>
              Engage in any automated use of the system, such as using scripts
              to send comments or messages, or using any data mining, robots, or
              similar data gathering and extraction tools, except as expressly
              permitted.
            </li>
            <li>
              Interfere with, disrupt, or create an undue burden on the Site or
              the networks or services connected to the Site.
            </li>
            <li>
              Attempt to impersonate another user or person or use the username
              of another user.
            </li>
            <li>
              Use any information obtained from the Site to harass, abuse, or
              harm another person.
            </li>
            <li>
              Upload or transmit viruses, Trojan horses, or other malicious
              material that interferes with any party&apos;s use and enjoyment
              of the Site.
            </li>
            <li>
              Use the Site in a manner inconsistent with any applicable laws or
              regulations, including but not limited to the Information
              Technology Act, 2000 (India) and its rules and amendments.
            </li>
          </ol>
        </div>

        <div>
          <h2>11. User generated contributions</h2>
          <p>
            The Site may allow you to submit feedback, suggestions, bug reports,
            or other content (collectively, &quot;Contributions&quot;). When you
            create or make available any Contributions, you represent and
            warrant that your Contributions are not false, inaccurate, or
            misleading, do not violate the intellectual property rights of any
            third party, and comply with applicable law.
          </p>
          <p>
            We have the right, in our sole and absolute discretion, to edit,
            redact, or otherwise change any Contributions, or to remove any
            Contributions at any time and for any reason, without notice.
          </p>
        </div>

        <div>
          <h2>12. Contribution license</h2>
          <p>
            By posting your Contributions to the Site, you grant us an
            unrestricted, unlimited, irrevocable, perpetual, non-exclusive,
            transferable, royalty-free, fully-paid, worldwide right and license
            to use, copy, reproduce, disclose, distribute, and prepare
            derivative works of such Contributions for any purpose, commercial
            or otherwise.
          </p>
          <p>
            You retain full ownership of all of your Contributions and any
            intellectual property rights associated with your Contributions. We
            do not assert any ownership over your Contributions.
          </p>
        </div>

        <div>
          <h2>13. Submissions</h2>
          <p>
            Any questions, comments, suggestions, ideas, feedback, or other
            information regarding the Site (&quot;Submissions&quot;) provided by
            you to us are non-confidential and shall become our sole property.
            We shall own exclusive rights, including all intellectual property
            rights, and shall be entitled to the unrestricted use and
            dissemination of these Submissions for any lawful purpose, without
            acknowledgment or compensation to you.
          </p>
        </div>

        <div>
          <h2>14. Third-party websites and content</h2>
          <p>
            The Site may contain links to third-party websites, including code
            hosting platforms and other developer tools (&quot;Third-Party
            Websites&quot;). Such Third-Party Websites are not investigated,
            monitored, or checked for accuracy, appropriateness, or completeness
            by us. We are not responsible for any Third-Party Websites accessed
            through the Site. Your use of Third-Party Websites is governed
            solely by the terms and policies of those third parties.
          </p>
        </div>

        <div>
          <h2>15. Site management</h2>
          <p>
            We reserve the right to: (1) monitor the Site for violations of
            these Terms of Use; (2) take appropriate legal action against anyone
            who, in our sole discretion, violates the law or these Terms of Use;
            (3) refuse, restrict access to, limit the availability of, or
            disable any of your Contributions or any portion thereof; and (4)
            otherwise manage the Site in a manner designed to protect our rights
            and property and to facilitate the proper functioning of the Site.
          </p>
        </div>

        <div>
          <h2>16. Privacy policy</h2>
          <p>
            We care about data privacy and security. Please review our Privacy
            Policy at{" "}
            <a className="tou-email" href="https://useabstractions.com/privacy">
              https://useabstractions.com/privacy
            </a>
            . By using the Site, you agree to be bound by our Privacy Policy,
            which is incorporated into these Terms of Use. The Site is hosted on
            servers that may be located outside India. If you access the Site
            from any region with laws governing personal data collection, use,
            or disclosure that differ from applicable Indian laws, you agree to
            your data being transferred to and processed in such locations.
          </p>
        </div>

        <div>
          <h2>17. Data protection and codebase data</h2>
          <p>
            You acknowledge and agree that to the extent any personal data or
            proprietary code is shared by you with us through the Service:
          </p>
          <ol>
            <li>
              You shall ensure you have the right to share such data or code
              with us, and that doing so does not violate any third-party rights
              or confidentiality obligations.
            </li>
            <li>
              We will handle your data in accordance with our Privacy Policy and
              applicable data protection laws, including the Digital Personal
              Data Protection Act, 2023 (India) as applicable.
            </li>
            <li>
              You are solely responsible for ensuring that repositories you
              connect do not contain sensitive personal data of third parties
              beyond what is necessary for the Service.
            </li>
            <li>
              If any government authority or court of competent jurisdiction
              requires us to disclose data, we shall have the right to do so in
              accordance with applicable law.
            </li>
          </ol>
        </div>

        <div>
          <h2>18. Term and termination</h2>
          <p>
            These Terms of Use shall remain in full force and effect while you
            use the Site. We reserve the right to, in our sole discretion and
            without notice or liability, deny access to and use of the Site
            (including blocking certain IP addresses), delete your account and
            any content or information that you posted, for any reason,
            including breach of these Terms of Use or any applicable law or
            regulation.
          </p>
          <p>
            If we terminate or suspend your account for any reason, you are
            prohibited from registering and creating a new account under your
            name, a fake or borrowed name, or the name of any third party.
          </p>
        </div>

        <div>
          <h2>19. Modifications and interruptions</h2>
          <p>
            We reserve the right to change, modify, or remove the contents of
            the Site at any time or for any reason at our sole discretion
            without notice. We also reserve the right to modify or discontinue
            all or part of the Site without notice at any time. We will not be
            liable to you or any third party for any modification, price change,
            suspension, or discontinuance of the Site.
          </p>
          <p>
            We cannot guarantee the Site will be available at all times. We may
            experience hardware, software, or other problems or need to perform
            maintenance related to the Site, resulting in interruptions, delays,
            or errors. Nothing in these Terms of Use will be construed to
            obligate us to maintain and support the Site or to supply any
            corrections, updates, or releases in connection therewith.
          </p>
        </div>

        <div>
          <h2>20. Governing law</h2>
          <p>
            These Terms shall be governed by and construed in accordance with
            the laws of India. Artery Ventures, LLP and you irrevocably consent
            that the courts of competent jurisdiction in Mumbai, Maharashtra,
            India shall have exclusive jurisdiction to resolve any dispute which
            may arise in connection with these Terms.
          </p>
        </div>

        <div>
          <h2>21. Dispute resolution</h2>
          <p>
            To expedite resolution and control the cost of any dispute,
            controversy, or claim related to these Terms of Use (each a
            &quot;Dispute&quot;), the parties agree to first attempt to
            negotiate any Dispute informally for at least thirty (30) days
            before initiating formal legal proceedings. Such informal
            negotiations commence upon written notice from one party to the
            other.
          </p>
          <p>
            If informal negotiations fail to resolve a Dispute, the parties
            agree to submit to the exclusive jurisdiction of the courts in
            Mumbai, Maharashtra, India. Notwithstanding the foregoing, either
            party may seek injunctive or other equitable relief in any court of
            competent jurisdiction to prevent the actual or threatened
            infringement of intellectual property rights.
          </p>
        </div>

        <div>
          <h2>22. Disclaimer</h2>
          <p className="caps">
            THE SITE IS PROVIDED ON AN AS-IS AND AS-AVAILABLE BASIS. YOU AGREE
            THAT YOUR USE OF THE SITE AND OUR SERVICES WILL BE AT YOUR SOLE
            RISK. TO THE FULLEST EXTENT PERMITTED BY LAW, WE DISCLAIM ALL
            WARRANTIES, EXPRESS OR IMPLIED, IN CONNECTION WITH THE SITE AND YOUR
            USE THEREOF, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTIES
            OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND
            NON-INFRINGEMENT. WE MAKE NO WARRANTIES OR REPRESENTATIONS ABOUT THE
            ACCURACY OR COMPLETENESS OF THE SITE&apos;S CONTENT OR THE ACCURACY
            OF ANY CODEBASE ANALYSIS, NEWSLETTER CONTENT, OR INSIGHTS GENERATED
            BY THE SERVICE.
          </p>
        </div>

        <div>
          <h2>23. Limitations of liability</h2>
          <p className="caps">
            IN NO EVENT WILL WE OR OUR DIRECTORS, EMPLOYEES, OR AGENTS BE LIABLE
            TO YOU OR ANY THIRD PARTY FOR ANY DIRECT, INDIRECT, CONSEQUENTIAL,
            EXEMPLARY, INCIDENTAL, SPECIAL, OR PUNITIVE DAMAGES, INCLUDING LOST
            PROFIT, LOST REVENUE, LOSS OF DATA, OR OTHER DAMAGES ARISING FROM
            YOUR USE OF THE SITE OR THE SERVICE, EVEN IF WE HAVE BEEN ADVISED OF
            THE POSSIBILITY OF SUCH DAMAGES. NOTWITHSTANDING ANYTHING TO THE
            CONTRARY CONTAINED HEREIN, OUR LIABILITY TO YOU FOR ANY CAUSE
            WHATSOEVER WILL AT ALL TIMES BE LIMITED TO THE AMOUNT PAID, IF ANY,
            BY YOU TO US IN THE THREE (3) MONTHS PRECEDING THE EVENT GIVING RISE
            TO THE CLAIM.
          </p>
          <p className="caps">
            WE ARE NOT LIABLE FOR ANY LOSS OR DAMAGE ARISING FROM UNAUTHORIZED
            ACCESS TO YOUR CODE REPOSITORIES OR ANY THIRD-PARTY SERVICES
            CONNECTED TO THE SITE, OR FROM ANY INACCURACIES IN CODEBASE ANALYSIS
            OR NEWSLETTER CONTENT GENERATED BY THE SERVICE.
          </p>
        </div>

        <div>
          <h2>24. Indemnification</h2>
          <p>
            You agree to defend, indemnify, and hold us harmless, including our
            subsidiaries, affiliates, officers, agents, partners, and employees,
            from and against any loss, damage, liability, claim, or demand,
            including reasonable attorneys&apos; fees and expenses, made by any
            third party due to or arising out of: (1) your Contributions; (2)
            use of the Site; (3) breach of these Terms of Use; (4) any breach of
            your representations and warranties set forth in these Terms of Use;
            (5) your violation of the rights of a third party, including
            intellectual property rights; or (6) your unauthorized sharing of
            third-party code or data through the Service.
          </p>
        </div>

        <div>
          <h2>25. Electronic communications, transactions, and signatures</h2>
          <p>
            Visiting the Site, sending us emails, and completing online forms
            constitute electronic communications. You consent to receive
            electronic communications, and you agree that all agreements,
            notices, disclosures, and other communications we provide to you
            electronically, via email and on the Site, satisfy any legal
            requirement that such communication be in writing, to the extent
            permitted by applicable law.
          </p>
        </div>

        <div>
          <h2>26. Corrections</h2>
          <p>
            There may be information on the Site that contains typographical
            errors, inaccuracies, or omissions, including descriptions, pricing,
            availability, and various other information. We reserve the right to
            correct any errors, inaccuracies, or omissions and to change or
            update the information on the Site at any time, without prior
            notice.
          </p>
        </div>

        <div>
          <h2>27. Miscellaneous</h2>
          <p>
            These Terms of Use and any policies or operating rules posted by us
            on the Site constitute the entire agreement and understanding
            between you and us. Our failure to exercise or enforce any right or
            provision of these Terms of Use shall not operate as a waiver of
            such right or provision. If any provision or part of a provision of
            these Terms of Use is determined to be unlawful, void, or
            unenforceable, that provision or part of the provision is deemed
            severable from these Terms of Use and does not affect the validity
            and enforceability of any remaining provisions. There is no joint
            venture, partnership, employment, or agency relationship created
            between you and us as a result of these Terms of Use or use of the
            Site.
          </p>
        </div>

        <div>
          <h2>28. Contact us</h2>
          <p>
            To resolve a complaint regarding the Site or to receive further
            information regarding use of the Site, please contact us at:
          </p>
          <div className="tou-addr">
            <strong>Artery Ventures, LLP</strong>
            <br />
            FO-02, 4th Floor, 28/A, 80 Feet Rd
            <br />
            Indiranagar, Bengaluru, Karnataka 560038
            <br />
            India
            <br />
            <br />
            Email:{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            <br />
            Website:{" "}
            <a className="tou-email" href="https://useabstractions.com">
              https://useabstractions.com
            </a>
          </div>
        </div>
      </div>
    </>
  );
}
