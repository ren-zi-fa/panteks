package validator

import (
	"errors"
	"net"
	"net/url"
	"regexp"
)

var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)

func ValidateTarget(targetUrl string) error {
	u, err := url.ParseRequestURI(targetUrl)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL must start with http:// or https://")
	}

	host := u.Hostname()
	if !domainRegex.MatchString(host) {
		return errors.New("invalid domain format (e.g., example.com)")
	}

	if _, err := net.LookupHost(host); err != nil {
		return errors.New("domain could not be resolved (not found in DNS)")
	}

	return nil
}
