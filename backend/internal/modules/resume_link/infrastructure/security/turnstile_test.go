package security

import (
	"context"
	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeHTTPClient struct {
	handler func(*http.Request) (*http.Response, error)
	calls   int
}

func (c *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.calls++
	return c.handler(req)
}

func TestNewTurnstileVerifierValidatesConfig(t *testing.T) {
	client := &fakeHTTPClient{}
	_, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)

	_, err = NewTurnstileVerifier(Config{Timeout: time.Second}, client)
	require.Error(t, err)

	_, err = NewTurnstileVerifier(Config{Secret: "secret"}, client)
	require.Error(t, err)

	_, err = NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, nil)
	require.Error(t, err)

	_, err = NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second, SiteverifyURL: "http://example.com"}, client)
	require.Error(t, err)
}

func TestTurnstileVerifierSuccessAndRequestPayload(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			require.Equal(t, http.MethodPost, req.Method)
			require.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			values := string(body)
			require.Contains(t, values, "secret=secret")
			require.Contains(t, values, "response=token")
			require.Contains(t, values, "remoteip=192.0.2.10")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":true,"hostname":"cv.example.com","action":"contact_view"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{
		Secret:         "secret",
		Timeout:        time.Second,
		ExpectedHost:   "cv.example.com",
		ExpectedAction: "contact_view",
	}, client)
	require.NoError(t, err)

	require.NoError(t, verifier.Verify(context.Background(), "token", "192.0.2.10"))
	require.Equal(t, 1, client.calls)
}

func TestTurnstileVerifierOmitsEmptyRemoteIP(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.NotContains(t, string(body), "remoteip=")
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":true}`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)
	require.NoError(t, verifier.Verify(context.Background(), "token", ""))
}

func TestTurnstileVerifierRejectsInvalidTokenShapeWithoutHTTPCall(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("must not be called")
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)

	require.ErrorIs(t, verifier.Verify(context.Background(), "", ""), ports.ErrTokenVerificationFailed)
	require.ErrorIs(t, verifier.Verify(context.Background(), strings.Repeat("a", ports.MaxTokenLength+1), ""), ports.ErrTokenVerificationFailed)
	require.Equal(t, 0, client.calls)
}

func TestTurnstileVerifierRejectsProviderFailureWithoutLeakingDetails(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":false,"error-codes":["invalid-input-secret"]}`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)

	err = verifier.Verify(context.Background(), "token", "")
	require.ErrorIs(t, err, ports.ErrTokenVerificationFailed)
	require.NotContains(t, err.Error(), "invalid-input-secret")
	require.NotContains(t, err.Error(), "secret")
}

func TestTurnstileVerifierRejectsWrongHostname(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":true,"hostname":"other.example.com"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second, ExpectedHost: "cv.example.com"}, client)
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(context.Background(), "token", ""), ports.ErrTokenVerificationFailed)
}

func TestTurnstileVerifierRejectsWrongAction(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"success":true,"action":"other_action"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second, ExpectedAction: "contact_view"}, client)
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(context.Background(), "token", ""), ports.ErrTokenVerificationFailed)
}

func TestTurnstileVerifierRejectsNon2xx(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader(`{"success":false}`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(context.Background(), "token", ""), ErrVerificationHTTP)
}

func TestTurnstileVerifierRejectsMalformedResponse(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`not-json`)),
				Header:     make(http.Header),
			}, nil
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(context.Background(), "token", ""), ErrVerificationHTTP)
}

func TestTurnstileVerifierMapsNetworkError(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network failure")
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)
	require.ErrorIs(t, verifier.Verify(context.Background(), "token", ""), ErrVerificationHTTP)
}

func TestTurnstileVerifierTimeout(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			<-req.Context().Done()
			return nil, req.Context().Err()
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: 5 * time.Millisecond}, client)
	require.NoError(t, err)

	err = verifier.Verify(context.Background(), "token", "")
	require.ErrorIs(t, err, ErrVerificationHTTP)
}

func TestTurnstileVerifierCallerCancellation(t *testing.T) {
	client := &fakeHTTPClient{
		handler: func(req *http.Request) (*http.Response, error) {
			<-req.Context().Done()
			return nil, req.Context().Err()
		},
	}
	verifier, err := NewTurnstileVerifier(Config{Secret: "secret", Timeout: time.Second}, client)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = verifier.Verify(ctx, "token", "")
	require.ErrorIs(t, err, ErrVerificationHTTP)
}
