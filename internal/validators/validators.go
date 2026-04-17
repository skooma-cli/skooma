// Package validators contains validation functions for user input and configuration.
package validators

import (
	"errors"
	"net/mail"
	"net/url"
	"strings"

	"github.com/skooma-cli/skooma/internal/sanitize"
)

// NotEmpty checks that a string is not empty or whitespace-only.
func NotEmpty(label string) func(string) error {
	return func(str string) error {
		if strings.TrimSpace(str) == "" {
			return errors.New(label + " can't be empty")
		}
		return nil
	}
}

// NoSpaces checks that a string contains no spaces.
func NoSpaces(label string) func(string) error {
	return func(str string) error {
		if strings.Contains(str, " ") {
			return errors.New(label + " can't contain spaces")
		}
		return nil
	}
}

// NoUnderscores checks that a string contains no underscores.
func NoUnderscores(label string) func(string) error {
	return func(str string) error {
		if strings.Contains(str, "_") {
			return errors.New(label + " can't contain underscores")
		}
		return nil
	}
}

// ValidURL checks that a string is a valid URL when prefixed with https://.
// Handles input that may or may not include an http/https prefix.
func ValidURL(label string) func(string) error {
	return func(str string) error {
		cleaned := sanitize.StripHTTPPrefix(str)
		u, err := url.ParseRequestURI("https://" + cleaned)
		if err != nil || u.Host == "" || !strings.Contains(u.Host, ".") {
			return errors.New(label + " must be a valid URL (e.g., github.com/user/repo)")
		}
		parts := strings.SplitN(u.Host, ".", 2)
		if parts[0] == "" || parts[1] == "" {
			return errors.New(label + " must be a valid URL (e.g., github.com/user/repo)")
		}
		return nil
	}
}

// RFC5322Address validates the "Name <email>" format using net/mail.
func RFC5322Address(label string) func(string) error {
	return func(str string) error {
		addr, err := mail.ParseAddress(str)
		if err != nil || addr.Name == "" {
			return errors.New(label + " must be in format: Name <email@example.com>")
		}
		return nil
	}
}

// AllowEmpty wraps one or more validators to allow empty strings. If the input
// is empty or whitespace-only, it skips validation and returns nil. Otherwise,
// it runs the validators in order.
func AllowEmpty(validators ...func(string) error) func(string) error {
	return func(str string) error {
		if strings.TrimSpace(str) == "" {
			return nil
		}
		return All(validators...)(str)
	}
}

// All composes multiple validators into one, running them in order and
// stopping at the first error.
func All(validators ...func(string) error) func(string) error {
	return func(str string) error {
		// TODO: consider running all validators instead of stopping at the first error, to provide more comprehensive feedback to the user
		for _, v := range validators {
			if err := v(str); err != nil {
				return err
			}
		}
		return nil
	}
}
