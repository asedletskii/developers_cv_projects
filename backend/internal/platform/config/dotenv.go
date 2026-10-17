package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var dotenvKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func loadDotEnv() error {
	if isNonLocalEnvironment(os.Getenv("APP_ENV")) {
		return nil
	}

	data, err := os.ReadFile(".env")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read .env: %w", err)
	}
	return loadDotEnvContent(data)
}

func isNonLocalEnvironment(value string) bool {
	switch strings.TrimSpace(value) {
	case "production", "staging":
		return true
	default:
		return false
	}
}

func loadDotEnvContent(data []byte) error {
	if isNonLocalEnvironment(os.Getenv("APP_ENV")) {
		return nil
	}

	if strings.TrimSpace(os.Getenv("APP_ENV")) == "" {
		for _, line := range strings.Split(string(data), "\n") {
			key, value, ok, err := parseDotEnvLine(line)
			if err != nil || !ok || key != "APP_ENV" {
				continue
			}
			if isNonLocalEnvironment(value) {
				return errors.New("APP_ENV=production/staging must be set in the process environment; .env is for local development only")
			}
			break
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		key, value, ok, err := parseDotEnvLine(scanner.Text())
		if err != nil {
			return fmt.Errorf(".env line %d: %w", lineNumber, err)
		}
		if !ok {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read .env: %w", err)
	}
	return nil
}

func parseDotEnvLine(rawLine string) (key, value string, ok bool, err error) {
	line := strings.TrimSpace(rawLine)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false, nil
	}
	if strings.HasPrefix(line, "export ") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
	}
	key, rawValue, found := strings.Cut(line, "=")
	if !found {
		return "", "", false, errors.New("expected KEY=VALUE assignment")
	}
	key = strings.TrimSpace(key)
	if !dotenvKeyPattern.MatchString(key) {
		return "", "", false, fmt.Errorf("invalid variable name %q", key)
	}
	rawValue = strings.TrimSpace(rawValue)
	if rawValue == "" {
		return key, "", true, nil
	}

	if rawValue[0] == '\'' || rawValue[0] == '"' {
		quote := rawValue[0]
		closing := -1
		for i := 1; i < len(rawValue); i++ {
			if quote == '"' && rawValue[i] == '\\' {
				i++
				continue
			}
			if rawValue[i] == quote {
				closing = i
				break
			}
		}
		if closing < 0 {
			return "", "", false, errors.New("unterminated quoted value")
		}
		remainder := rawValue[closing+1:]
		trimmedRemainder := strings.TrimSpace(remainder)
		if trimmedRemainder != "" && (len(remainder) == 0 || (remainder[0] != ' ' && remainder[0] != '\t') || !strings.HasPrefix(trimmedRemainder, "#")) {
			return "", "", false, errors.New("unexpected characters after quoted value")
		}
		quoted := rawValue[:closing+1]
		if quote == '\'' {
			return key, quoted[1 : len(quoted)-1], true, nil
		}
		unquoted, unquoteErr := strconv.Unquote(quoted)
		if unquoteErr != nil {
			return "", "", false, fmt.Errorf("invalid double-quoted value: %w", unquoteErr)
		}
		return key, unquoted, true, nil
	}

	for i := 1; i < len(rawValue); i++ {
		if rawValue[i] == '#' && (rawValue[i-1] == ' ' || rawValue[i-1] == '\t') {
			rawValue = strings.TrimRight(rawValue[:i], " \t")
			break
		}
	}
	return key, rawValue, true, nil
}
