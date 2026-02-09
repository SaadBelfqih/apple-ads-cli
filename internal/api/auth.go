package api

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/SaadBelfqih/apple-ads-cli/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenURL    = "https://appleid.apple.com/auth/oauth2/token"
	tokenScope  = "searchadsorg"
	jwtLifetime = 180 * 24 * time.Hour
	// Refresh token 60 seconds before expiry
	tokenRefreshBuffer = 60 * time.Second
)

// TokenSource manages OAuth2 token lifecycle.
type TokenSource struct {
	cfg        *config.Config
	privateKey *ecdsa.PrivateKey

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

// NewTokenSource creates a TokenSource from config.
func NewTokenSource(cfg *config.Config) (*TokenSource, error) {
	keyData, err := os.ReadFile(cfg.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %s", cfg.PrivateKeyPath)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		ecKey, ecErr := x509.ParseECPrivateKey(block.Bytes)
		if ecErr != nil {
			return nil, fmt.Errorf("parse private key: %w (also tried EC: %v)", err, ecErr)
		}
		key = ecKey
	}

	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not ECDSA")
	}

	return &TokenSource{
		cfg:        cfg,
		privateKey: ecKey,
	}, nil
}

// Token returns a valid access token, refreshing if necessary.
func (ts *TokenSource) Token() (string, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.accessToken != "" && time.Now().Before(ts.expiresAt.Add(-tokenRefreshBuffer)) {
		return ts.accessToken, nil
	}

	return ts.refresh()
}

func (ts *TokenSource) refresh() (string, error) {
	clientSecret, err := ts.buildJWT()
	if err != nil {
		return "", fmt.Errorf("build JWT: %w", err)
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {ts.cfg.ClientID},
		"client_secret": {clientSecret},
		"scope":         {tokenScope},
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Apple returns OAuth2-style error objects here, for example:
		// {"error":"invalid_client","error_description":"..."}.
		var oauthErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if json.Unmarshal(body, &oauthErr) == nil && oauthErr.Error != "" {
			if oauthErr.ErrorDescription != "" {
				return "", fmt.Errorf("token exchange failed (%d): %s: %s", resp.StatusCode, oauthErr.Error, oauthErr.ErrorDescription)
			}
			return "", fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, oauthErr.Error)
		}

		return "", fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}

	if tr.AccessToken == "" {
		return "", fmt.Errorf("no access_token in response: %s", string(body))
	}

	ts.accessToken = tr.AccessToken
	ts.expiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)

	return ts.accessToken, nil
}

// Invalidate clears the cached access token so the next Token() call refreshes.
func (ts *TokenSource) Invalidate() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.accessToken = ""
	ts.expiresAt = time.Time{}
}

func (ts *TokenSource) buildJWT() (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    ts.cfg.TeamID,
		Subject:   ts.cfg.ClientID,
		Audience:  jwt.ClaimStrings{"https://appleid.apple.com"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(jwtLifetime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = ts.cfg.KeyID

	return token.SignedString(ts.privateKey)
}
