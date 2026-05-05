// WHATWG HTML5 email spec: enforces proper domain label structure,
// rejects consecutive dots in domain, and supports subdomains (e.g. user@mail.domain.com).
const EMAIL_REGEX =
  /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;

/** Returns true if the string looks like a valid email address. */
export function isValidEmail(email: string): boolean {
  return EMAIL_REGEX.test(email);
}

/**
 * Returns true if the string looks like a valid phone number.
 * Accepts 7–15 digits optionally prefixed with a "+" country code.
 * Common formatting characters (spaces, hyphens, parentheses, dots, slashes)
 * are stripped before the check so e.g. "(+1) 555-1234" is considered valid.
 */
export function isValidPhone(phone: string): boolean {
  const cleaned = phone.replace(/[\s\-()\./]/g, "");
  if (!cleaned) return true;
  if (cleaned.startsWith("+")) {
    return /^\+\d{7,15}$/.test(cleaned);
  }
  return /^\d{7,15}$/.test(cleaned);
}
