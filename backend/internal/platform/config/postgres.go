package config

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

func (c PostgresConfig) ConnectionString() string {
	values := url.Values{}
	values.Set("sslmode", c.SSLMode)
	var options []string
	if c.StatementTimeout > 0 {
		options = append(options, "-c statement_timeout="+strconv.FormatInt(c.StatementTimeout.Milliseconds(), 10))
	}
	if c.LockTimeout > 0 {
		options = append(options, "-c lock_timeout="+strconv.FormatInt(c.LockTimeout.Milliseconds(), 10))
	}
	if c.IdleInTransactionTimeout > 0 {
		options = append(options, "-c idle_in_transaction_session_timeout="+strconv.FormatInt(c.IdleInTransactionTimeout.Milliseconds(), 10))
	}
	if len(options) > 0 {
		values.Set("options", strings.Join(options, " "))
	}

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     "/" + c.Database,
		RawQuery: values.Encode(),
	}).String()
}

func isValidPostgresSSLMode(value string) bool {
	switch value {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
		return true
	default:
		return false
	}
}
