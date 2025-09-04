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
		return errors.New("URL tidak valid")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("URL harus diawali http:// atau https://")
	}

	host := u.Hostname()
	if !domainRegex.MatchString(host) {
		return errors.New("domain tidak valid (harus format seperti example.com)")
	}

	if _, err := net.LookupHost(host); err != nil {
		return errors.New("domain tidak bisa di-resolve (tidak ditemukan di DNS)")
	}

	return nil
}
