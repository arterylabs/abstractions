import type { Metadata } from "next";
import { PageHeader } from "@/components/ui/page-header";

export const metadata: Metadata = {
  title: "Privacy Policy",
  description:
    "Read the Abstractions Privacy Policy to understand how we collect, use, and protect your personal data.",
  alternates: { canonical: "/privacy" },
};

export default function Privacy() {
  return (
    <>
      <div className="mx-auto max-w-[1200px] w-full docs-content py-12 px-8">
        <PageHeader
          display
          title="Privacy Policy"
          description={
            <>
              Last updated: May 2026 &nbsp;·&nbsp; Artery Ventures, LLP
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
          <p>
            We, Artery Ventures, LLP (&quot;we&quot;, &quot;us&quot;,
            &quot;our&quot;, &quot;Company&quot;) are the owners of the website{" "}
            <a href="https://useabstractions.com">useabstractions.com</a>{" "}
            (&quot;Website&quot;) and the Abstractions application
            (&quot;Application&quot;). The Website and the Application are
            collectively referred to as the &quot;Platform&quot;. The Website
            showcases our Application and the features of our Application,
            including codebase analysis, automated newsletter generation, and
            developer tools (&quot;Services&quot;).
          </p>
          <p>
            We respect data privacy rights and are committed to protecting
            personal information collected on this Platform. This Privacy Policy
            (&quot;Privacy Policy&quot;) sets forth how we collect, use, and
            protect the personal information collected on our Platform.
          </p>
          <p>
            PLEASE READ THIS PRIVACY POLICY CAREFULLY. BY CONTINUING TO USE THE
            PLATFORM OR PROVIDING US PERSONAL INFORMATION, YOU CONSENT TO OUR
            USE OF YOUR PERSONAL INFORMATION IN ACCORDANCE WITH THE TERMS OF
            THIS PRIVACY POLICY. IF YOU DO NOT AGREE TO THIS PRIVACY POLICY, YOU
            MAY WITHDRAW YOUR CONSENT OR ALTERNATIVELY CHOOSE NOT TO PROVIDE
            YOUR PERSONAL INFORMATION ON THE PLATFORM. SUCH AN INTIMATION TO
            WITHDRAW YOUR CONSENT CAN BE PROVIDED BY EMAIL AT{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            .
          </p>
          <p>
            IF YOU ARE ACCESSING THE PLATFORM ON BEHALF OF A THIRD PARTY, YOU
            REPRESENT THAT YOU HAVE THE AUTHORITY TO BIND SUCH THIRD PARTY TO
            THE TERMS AND CONDITIONS OF THIS PRIVACY POLICY AND, IN SUCH AN
            EVENT, YOUR USE OF THE PLATFORM SHALL REFER TO USE BY SUCH THIRD
            PARTY. IF YOU DO NOT HAVE SUCH AUTHORITY OR DO NOT AGREE TO THE
            TERMS OF THIS PRIVACY POLICY, THEN YOU SHOULD REFRAIN FROM USING THE
            PLATFORM.
          </p>
          <p>
            This Privacy Policy is an electronic record in the form of an
            electronic contract construed in accordance with applicable data
            protection laws.
          </p>
        </div>

        <div>
          <h2>1. Definitions</h2>
          <p>
            &quot;User(s)&quot;, &quot;you&quot;, &quot;your&quot; shall mean
            and include individuals, business organizations, commercial
            establishments, and their permitted users who use the Platform or
            may opt to share Personal Information to avail the Services through
            the Platform.
          </p>
          <p>
            &quot;UK Data Protection Law&quot; means the UK GDPR, the United
            Kingdom Data Protection Act 2018, the Privacy and Electronic
            Communications Regulations, and any regulation superseding any of
            the foregoing.
          </p>
        </div>

        <div>
          <h2>2. Personal information collected</h2>
          <p>
            For purposes of this Privacy Policy, &quot;Personal
            Information&quot; means information that can be used to personally
            identify the User, including but not limited to a User&apos;s name,
            mobile number, email address, and company name.
          </p>
          <p>
            <strong>Personal Information collected on the Website</strong> — We
            may collect Personal Information from you when you (a) visit our
            Website, (b) sign up to create an account, (c) contact us through
            our contact page, or (d) book a product demo. Information collected
            may include your name, email address, phone number, company name,
            and country.
          </p>
          <p>
            <strong>Personal Information collected via the Application</strong>{" "}
            — We may collect your Personal Information including your name,
            email address, and contact number when you sign up, create an
            account, and use the Services on our Platform. We also collect
            information about the code repositories you connect, including
            repository metadata and codebase content, solely to provide the
            Services.
          </p>
          <p>
            <strong>Social Media Platforms</strong> — If you contact us through
            social media platforms such as LinkedIn, Twitter/X, or others, we
            may receive your publicly available profile information such as your
            name, email, and handle.
          </p>
          <p>
            <strong>Log Data</strong> — When you visit the Website or use the
            Application, our servers automatically record information your
            browser sends (&quot;Log Data&quot;). This Log Data may include
            information such as your IP address, browser type, referring
            webpage, pages visited, time spent, and other analytical data.
          </p>
          <p>
            <strong>Payment Information</strong> — If you elect to pay for the
            Services, we may collect your name, email address, and a customer ID
            assigned by our payment processing partner. Payment card details are
            entered directly with our payment processors — we do not access,
            store, or collect your credit/debit card information.
          </p>
        </div>

        <div>
          <h2>3. Cookies</h2>
          <p>
            We use cookies and similar tracking tools to track user traffic
            patterns and improve our Website and Services. You may choose to
            disable cookies through your browser settings. For more information,
            please refer to our cookie policy at{" "}
            <a href="https://useabstractions.com/cookies">
              useabstractions.com/cookies
            </a>
            .
          </p>
        </div>

        <div>
          <h2>4. Accuracy of information</h2>
          <p>
            You shall be solely responsible for the accuracy, correctness, and
            truthfulness of the Personal Information you share with us, whether
            your own or that of a third party. If you share Personal Information
            on behalf of a third party, you represent that you have the
            necessary authority and have obtained appropriate consent from that
            third party. We shall not be responsible for verifying such
            authority or consent.
          </p>
        </div>

        <div>
          <h2>5. Use of personal information</h2>
          <p>We use your Personal Information for the following purposes:</p>
          <ul>
            <li>
              to inform you about our Services and to respond to your requests;
            </li>
            <li>
              for creation or development of business intelligence or data
              analytics in relation to the Services;
            </li>
            <li>
              to provide you with a better experience when you access our
              Platform and to improve the Services;
            </li>
            <li>to maintain and manage our Platform;</li>
            <li>to manage our relationship with you;</li>
            <li>for internal record keeping; and</li>
            <li>to comply with our legal or statutory obligations.</li>
          </ul>
        </div>

        <div>
          <h2>6. Disclosures</h2>
          <p>
            We do not sell, rent, share, distribute, lease, or otherwise provide
            your Personal Information to third parties without your prior
            consent, except in the following cases:
          </p>
          <p>
            <strong>Affiliates</strong> — We may provide your Personal
            Information to our affiliates to enable them to improve the Services
            and respond to queries.
          </p>
          <p>
            <strong>Service Providers</strong> — We may share your Personal
            Information with service providers who work with us in connection
            with operating the Platform and/or providing the Services. All such
            service providers are subject to stringent confidentiality
            restrictions consistent with this Privacy Policy.
          </p>
          <p>
            <strong>Merger or Acquisition</strong> — We may transfer your
            Personal Information if we are acquired by another entity, merge
            with another company, or transfer part of our business to a third
            party. Any such third party that receives your Personal Information
            shall have the right to continue to use it in line with the purposes
            set out herein. We may notify you of such a transfer.
          </p>
          <p>
            <strong>Legal and Regulatory Authorities</strong> — We may disclose
            your Personal Information in order to comply with legal obligations,
            court orders, or requests by government authorities.
          </p>
        </div>

        <div>
          <h2>7. Data retention</h2>
          <p>
            We will retain your Personal Information for as long as required to
            provide the Services. We may also retain and use your Personal
            Information as necessary to comply with legal obligations, resolve
            disputes, and enforce our agreements.
          </p>
        </div>

        <div>
          <h2>8. Security</h2>
          <p>
            Your Personal Information is stored on third-party cloud
            infrastructure. Although we provide appropriate firewalls and
            protections, we cannot guarantee the security of Personal
            Information transmitted over the internet, as these systems are not
            hack-proof. We will take necessary measures to mitigate data
            security incidents.
          </p>
        </div>

        <div>
          <h2>9. Your rights</h2>
          <p>
            You have the right to access Personal Information in our possession,
            the right to have us rectify or modify such Personal Information,
            the right to have us erase or delete your Personal Information, the
            right to restrict our processing of such Personal Information, the
            right to object to our use of Personal Information, and the right to
            withdraw consent at any time where we rely on consent to process
            your Personal Information. If you would like to exercise any of
            these rights, please contact{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            .
          </p>
        </div>

        <div>
          <h2>10. Choice and opt-out</h2>
          <p>
            We may send you communications including notices about your use of
            the Platform, updates, promotional information, and newsletters. You
            may opt out of receiving promotional emails and newsletters by
            following the unsubscribe instructions in those emails.
            Alternatively, you can opt out at any time by emailing{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>{" "}
            with your specific request.
          </p>
        </div>

        <div>
          <h2>11. Information for EU and UK visitors</h2>
          <p>
            Residents of the European Union (&quot;EU&quot;) and United Kingdom
            (&quot;UK&quot;) should note that this Privacy Policy has been
            updated in accordance with the requirements of the EU General Data
            Protection Regulation (&quot;GDPR&quot;) and UK Privacy Laws. We
            shall be considered the Controllers of the Personal Information
            collected as a result of your use or access of the Website and
            Processors of the Personal Information collected as a result of your
            use of the Services.
          </p>
          <p>
            <strong>Legal Basis</strong> — We will not process your Personal
            Information without a lawful basis to do so. We will process your
            Personal Information only on the legal bases of consent, contract,
            or on the basis of our legitimate interests, provided that such
            interests are not overridden by your privacy rights and interests.
          </p>
          <p>
            <strong>Cross-border transfers</strong> — Personal Information of EU
            and UK residents may be processed outside the EU and UK. We collect
            and transfer such Personal Information in accordance with the
            provisions of the GDPR and UK Privacy Laws. If you have questions,
            please contact{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            .
          </p>
          <p>
            <strong>Your rights (EU and UK residents)</strong> — In addition to
            the rights described in Section 9, you have the right to lodge a
            complaint with a data protection authority. UK residents have the
            right to make a complaint to the Information Commissioner&apos;s
            Office (&quot;ICO&quot;) at{" "}
            <a href="https://www.ico.org.uk">www.ico.org.uk</a>. We would
            appreciate the chance to address your concerns before you approach
            the ICO, so please contact us in the first instance at{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>
            .
          </p>
          <p>
            <strong>Governing Laws</strong> — For EU and UK residents, this
            Privacy Policy shall be governed respectively by the provisions of
            the GDPR and UK Privacy Laws.
          </p>
        </div>

        <div>
          <h2>12. Links to other websites</h2>
          <p>
            Our Platform may contain links to other websites or applications. We
            do not control such third-party websites and are not responsible for
            the protection or privacy of any information you provide while
            visiting them. You should review the privacy policy applicable to
            any third-party website you visit.
          </p>
        </div>

        <div>
          <h2>13. Limitation of liability</h2>
          <p>
            To the extent permissible under law, we shall not be liable for any
            direct, indirect, incidental, special, consequential, or exemplary
            damages arising out of this Privacy Policy, including but not
            limited to damages for loss of profits, goodwill, data, or other
            intangible losses.
          </p>
        </div>

        <div>
          <h2>14. Children&apos;s privacy</h2>
          <p>
            This Platform is not intended for children under the age of 18. We
            do not knowingly collect Personal Information from children under
            the age of 18 without prior, verifiable consent of a legal
            representative. If you are under 18, please do not provide any
            Personal Information on this Platform. If a legal representative
            discovers that a child has provided us with Personal Information,
            please contact us at{" "}
            <a className="tou-email" href="mailto:support@useabstractions.com">
              support@useabstractions.com
            </a>{" "}
            to have the information deleted.
          </p>
        </div>

        <div>
          <h2>15. Governing law and dispute resolution</h2>
          <p>
            This Privacy Policy shall in all respects be governed by and
            construed in accordance with the laws of the State of California,
            USA. Any dispute arising under this Privacy Policy shall be subject
            to the exclusive jurisdiction of the courts of California.
          </p>
        </div>

        <div>
          <h2>16. Changes to this policy</h2>
          <p>
            Please revisit this page periodically to stay aware of any changes
            to this Privacy Policy. If we modify this Privacy Policy, we will
            make it available through the Website and indicate the date of the
            latest revision. If such modifications materially alter your rights
            or obligations, we will make reasonable efforts to notify you via
            email or through our Website.
          </p>
        </div>

        <div>
          <h2>17. Contact us</h2>
          <p>
            If you have any questions, concerns, or grievances about this
            Privacy Policy, or wish to withdraw your consent in relation to the
            processing of your Personal Information, please contact us at:
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
            <a href="https://useabstractions.com">
              https://useabstractions.com
            </a>
          </div>
        </div>
      </div>
    </>
  );
}
