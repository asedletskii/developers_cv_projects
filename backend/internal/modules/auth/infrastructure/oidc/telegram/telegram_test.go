package telegram

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"developers-cv/backend/internal/modules/auth/application/ports"
	"developers-cv/backend/internal/modules/auth/domain"
	"developers-cv/backend/internal/modules/auth/domain/valueobjects"
)

type countingRoundTripper struct {
	base  http.RoundTripper
	count *atomic.Int32
}

func (r countingRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	r.count.Add(1)
	return r.base.RoundTrip(request)
}

func TestAuthorizationURLContainsRequiredOIDCAndPKCEParameters(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	server := newOIDCTestServer(t, privateKey, func(string) string { return "" })
	defer server.Close()

	adapter, err := New(Config{
		ClientID: "client-id", ClientSecret: "client-secret",
		RedirectURL: server.URL + "/callback", IssuerURL: server.URL,
	}, http.DefaultClient, nil)
	require.NoError(t, err)
	challenge := "challenge-value"

	rawURL, err := adapter.AuthorizationURL(context.Background(), ports.TelegramAuthorizationParams{
		State:         "state-value",
		Nonce:         "nonce-value",
		CodeChallenge: challenge,
		RedirectURI:   server.URL + "/callback",
	})
	require.NoError(t, err)

	parsed, err := url.Parse(rawURL)
	require.NoError(t, err)
	query := parsed.Query()
	require.Equal(t, server.URL+"/auth", rawURL[:strings.Index(rawURL, "?")])
	require.Equal(t, "code", query.Get("response_type"))
	require.Equal(t, "client-id", query.Get("client_id"))
	require.Equal(t, server.URL+"/callback", query.Get("redirect_uri"))
	require.Equal(t, "openid profile", query.Get("scope"))
	require.Equal(t, "state-value", query.Get("state"))
	require.Equal(t, "nonce-value", query.Get("nonce"))
	require.Equal(t, challenge, query.Get("code_challenge"))
	require.Equal(t, "S256", query.Get("code_challenge_method"))
}

func TestNewRejectsInvalidConfig(t *testing.T) {
	_, err := New(Config{}, &http.Client{Timeout: time.Second}, nil)
	require.Error(t, err)
}

func TestAuthorizationURLRejectsRedirectURIOverride(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	server := newOIDCTestServer(t, privateKey, func(string) string { return "" })
	defer server.Close()

	adapter, err := New(Config{
		ClientID: "client-id", ClientSecret: "client-secret",
		RedirectURL: server.URL + "/callback", IssuerURL: server.URL,
	}, http.DefaultClient, nil)
	require.NoError(t, err)

	_, err = adapter.AuthorizationURL(context.Background(), ports.TelegramAuthorizationParams{
		State: "state", RedirectURI: server.URL + "/unexpected-callback",
	})
	require.ErrorContains(t, err, "redirect uri does not match")
}

func TestNewDoesNotPerformDiscoveryUntilFirstUse(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	server := newOIDCTestServer(t, privateKey, func(issuer string) string { return "" })
	defer server.Close()

	var discoveryRequests atomic.Int32
	client := &http.Client{
		Timeout:   time.Second,
		Transport: countingRoundTripper{base: http.DefaultTransport, count: &discoveryRequests},
	}
	lazy, err := New(Config{
		ClientID: "client-id", ClientSecret: "client-secret",
		RedirectURL: server.URL + "/callback", IssuerURL: server.URL,
	}, client, nil)
	require.NoError(t, err)
	require.EqualValues(t, 0, discoveryRequests.Load())

	_, err = lazy.AuthorizationURL(context.Background(), ports.TelegramAuthorizationParams{
		State: "state", RedirectURI: server.URL + "/callback",
	})
	require.NoError(t, err)
	require.Greater(t, discoveryRequests.Load(), int32(0))
	requestsAfterInitialization := discoveryRequests.Load()
	_, err = lazy.AuthorizationURL(context.Background(), ports.TelegramAuthorizationParams{
		State: "state-again", RedirectURI: server.URL + "/callback",
	})
	require.NoError(t, err)
	require.Equal(t, requestsAfterInitialization, discoveryRequests.Load(), "successful discovery should be cached")
}

func TestLazyProviderDiscoveryFailureIsReportedAsUnavailableAndBackedOff(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "temporarily unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	var discoveryRequests atomic.Int32
	client := &http.Client{
		Timeout:   time.Second,
		Transport: countingRoundTripper{base: http.DefaultTransport, count: &discoveryRequests},
	}
	lazy, err := New(Config{
		ClientID: "client-id", ClientSecret: "client-secret",
		RedirectURL: server.URL + "/callback", IssuerURL: server.URL,
	}, client, nil)
	require.NoError(t, err)
	params := ports.TelegramAuthorizationParams{State: "state", RedirectURI: server.URL + "/callback"}

	_, err = lazy.AuthorizationURL(context.Background(), params)
	require.ErrorIs(t, err, ports.ErrProviderUnavailable)
	requestsAfterFirstFailure := discoveryRequests.Load()
	require.Greater(t, requestsAfterFirstFailure, int32(0))
	_, err = lazy.AuthorizationURL(context.Background(), params)
	require.ErrorIs(t, err, ports.ErrProviderUnavailable)
	require.Equal(t, requestsAfterFirstFailure, discoveryRequests.Load(), "backoff should prevent an immediate discovery retry")
}

func TestExchangeVerifiesIDTokenAndNonce(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return signedIDToken(t, privateKey, map[string]any{
			"iss":   issuer,
			"sub":   "telegram-subject",
			"aud":   "client-id",
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(5 * time.Minute).Unix(),
			"nonce": "nonce-value",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	result, err := adapter.Exchange(context.Background(), "authorization-code", server.URL+"/callback", "verifier", nonceHash(t, "nonce-value"))
	require.NoError(t, err)
	require.Equal(t, domain.ProviderTelegram, result.Provider)
	require.Equal(t, "telegram-subject", result.ProviderSubject.String())
}

func TestExchangeRejectsNonceMismatch(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return signedIDToken(t, privateKey, map[string]any{
			"iss":   issuer,
			"sub":   "telegram-subject",
			"aud":   "client-id",
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(5 * time.Minute).Unix(),
			"nonce": "provider-nonce",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	_, err = adapter.Exchange(context.Background(), "authorization-code", server.URL+"/callback", "verifier", nonceHash(t, "stored-nonce"))
	require.ErrorContains(t, err, "nonce mismatch")
}

func TestExchangeRejectsRedirectURIOverride(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return signedIDToken(t, privateKey, map[string]any{
			"iss":   issuer,
			"sub":   "telegram-subject",
			"aud":   "client-id",
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(5 * time.Minute).Unix(),
			"nonce": "nonce-value",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	_, err = adapter.Exchange(context.Background(), "authorization-code", server.URL+"/unexpected-callback", "verifier", nonceHash(t, "nonce-value"))
	require.ErrorContains(t, err, "redirect uri does not match")
}

func TestExchangeRejectsWrongIssuer(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return signedIDToken(t, privateKey, map[string]any{
			"iss":   issuer + "/wrong",
			"sub":   "telegram-subject",
			"aud":   "client-id",
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(5 * time.Minute).Unix(),
			"nonce": "nonce-value",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	_, err = adapter.Exchange(context.Background(), "authorization-code", server.URL+"/callback", "verifier", nonceHash(t, "nonce-value"))
	require.Error(t, err)
}

func TestExchangeRejectsWrongAudience(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return signedIDToken(t, privateKey, map[string]any{
			"iss":   issuer,
			"sub":   "telegram-subject",
			"aud":   "another-client",
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(5 * time.Minute).Unix(),
			"nonce": "nonce-value",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	_, err = adapter.Exchange(context.Background(), "authorization-code", server.URL+"/callback", "verifier", nonceHash(t, "nonce-value"))
	require.Error(t, err)
}

func TestExchangeRejectsExpiredIDToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return signedIDToken(t, privateKey, map[string]any{
			"iss":   issuer,
			"sub":   "telegram-subject",
			"aud":   "client-id",
			"iat":   time.Now().Add(-10 * time.Minute).Unix(),
			"exp":   time.Now().Add(-5 * time.Minute).Unix(),
			"nonce": "nonce-value",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	_, err = adapter.Exchange(context.Background(), "authorization-code", server.URL+"/callback", "verifier", nonceHash(t, "nonce-value"))
	require.Error(t, err)
}

func TestExchangeRejectsNonRS256IDToken(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	server := newOIDCTestServer(t, privateKey, func(issuer string) string {
		return unsignedStyleHS256Token(t, map[string]any{
			"iss":   issuer,
			"sub":   "telegram-subject",
			"aud":   "client-id",
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(5 * time.Minute).Unix(),
			"nonce": "nonce-value",
		})
	})
	defer server.Close()

	adapter := newTestAdapter(t, server.URL)
	_, err = adapter.Exchange(context.Background(), "authorization-code", server.URL+"/callback", "verifier", nonceHash(t, "nonce-value"))
	require.Error(t, err)
}

func nonceHash(t *testing.T, value string) valueobjects.NonceHash {
	t.Helper()
	hash := sha256.Sum256([]byte(value))
	result, err := valueobjects.NewNonceHash(hash[:])
	require.NoError(t, err)
	return result
}

func newTestAdapter(t *testing.T, issuer string) ports.TelegramOIDC {
	t.Helper()
	adapter, err := New(Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  issuer + "/callback",
		IssuerURL:    issuer,
	}, http.DefaultClient, nil)
	require.NoError(t, err)
	return adapter
}

func newOIDCTestServer(t *testing.T, privateKey *rsa.PrivateKey, idToken func(issuer string) string) *httptest.Server {
	t.Helper()
	publicKey := &privateKey.PublicKey
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	issuer := server.URL
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{
			"issuer":                                issuer,
			"authorization_endpoint":                issuer + "/auth",
			"token_endpoint":                        issuer + "/token",
			"jwks_uri":                              issuer + "/jwks",
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "authorization-code", r.FormValue("code"))
		require.Equal(t, "verifier", r.FormValue("code_verifier"))
		writeJSON(t, w, map[string]any{
			"access_token": "test-access-token",
			"token_type":   "Bearer",
			"id_token":     idToken(issuer),
		})
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, map[string]any{"keys": []any{rsaJWK(publicKey)}})
	})
	return server
}

func signedIDToken(t *testing.T, key *rsa.PrivateKey, claims map[string]any) string {
	t.Helper()
	return signJWT(t, key, "RS256", claims)
}

func unsignedStyleHS256Token(t *testing.T, claims map[string]any) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString(mustJSON(t, map[string]any{"alg": "HS256", "typ": "JWT"}))
	payload := base64.RawURLEncoding.EncodeToString(mustJSON(t, claims))
	return header + "." + payload + ".invalid-signature"
}

func signJWT(t *testing.T, key *rsa.PrivateKey, algorithm string, claims map[string]any) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString(mustJSON(t, map[string]any{"alg": algorithm, "kid": "test-key", "typ": "JWT"}))
	payload := base64.RawURLEncoding.EncodeToString(mustJSON(t, claims))
	input := header + "." + payload
	hash := sha256.Sum256([]byte(input))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	require.NoError(t, err)
	return input + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func rsaJWK(key *rsa.PublicKey) map[string]any {
	return map[string]any{
		"kty": "RSA",
		"kid": "test-key",
		"use": "sig",
		"alg": "RS256",
		"n":   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString([]byte{1, 0, 1}),
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	return data
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(value))
}
