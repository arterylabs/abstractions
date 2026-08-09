package email_utils

import (
	"regexp"
	"strings"
)

func init() {
	updateBlockedDomainsCache()
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEmailValid checks if the email provided passes the required structure and length.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func GetWebsiteFromEmail(e string) string {
	tld := strings.Split(e, "@")
	if len(tld) < 2 {
		return ""
	}
	return tld[1]
}

var blockedDomainsCache = make(map[string]bool)

func updateBlockedDomainsCache() {
	for _, domain := range DISPOSABLE_EMAIL_DOMAINS {
		blockedDomainsCache[domain] = true
	}
}

// IsEmailDomainBlocked returns true when the email's domain is a known disposable provider.
// Enforced for all email send providers via the shared sendEmail path in clients/email.
func IsEmailDomainBlocked(e string) bool {
	tld := GetWebsiteFromEmail(strings.ToLower(e))
	return blockedDomainsCache[tld]
}

// Small curated list of common disposable-email domains. Prefer a short explicit set
// over a multi-thousand-line dump that is hard to audit.
var DISPOSABLE_EMAIL_DOMAINS = []string{
	"mailinator.com",
	"guerrillamail.com",
	"tempmail.com",
	"10minutemail.com",
	"throwaway.email",
	"yopmail.com",
	"trashmail.com",
	"temp-mail.org",
	"fakeinbox.com",
	"getnada.com",
}
