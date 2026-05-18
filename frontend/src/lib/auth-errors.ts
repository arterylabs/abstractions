export const oauthErrorMessages: Record<string, string> = {
  invalid_state:
    "Authentication failed due to an invalid state. Please try again.",
  unable_to_create_user: "We couldn't create your account. Please try again.",
  unable_to_create_authenticator:
    "We couldn't link your account. Please try again.",
  unable_to_create_workspace:
    "We couldn't set up your workspace. Please try again.",
  unable_to_set_default_workspace:
    "We couldn't configure your workspace. Please try again.",
  unable_to_find_user: "We couldn't find your account. Please try signing up.",
  unable_to_find_authenticators:
    "We couldn't find your login method. Please try again.",
  unable_to_update_authenticator:
    "We couldn't update your login method. Please try again.",
  unable_to_update_user: "We couldn't update your account. Please try again.",
  no_email:
    "Your account has no email address. Please use a different login method.",
};

export const oauthErrorFallback = "Authentication failed. Please try again.";

export function getOAuthErrorMessage(errorCode: string | null): string {
  if (!errorCode) return oauthErrorFallback;
  return oauthErrorMessages[errorCode] ?? oauthErrorFallback;
}
