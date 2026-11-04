package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Verifier проверяет access token Keycloak по ключам JWKS.
type Verifier struct {
	issuer   string
	jwksURL  string
	clientID string
	client   *http.Client

	mu   sync.Mutex
	keys map[string]*rsa.PublicKey
}

// NewVerifier создает проверку токенов издателя issuer.
func NewVerifier(issuer, jwksURL, clientID string) *Verifier {
	return &Verifier{
		issuer:   issuer,
		jwksURL:  jwksURL,
		clientID: clientID,
		client:   &http.Client{Timeout: 5 * time.Second},
		keys:     map[string]*rsa.PublicKey{},
	}
}

// Verify проверяет подпись, издателя, срок и аудиторию токена.
func (v *Verifier) Verify(ctx context.Context, raw string) error {
	_, err := jwt.Parse(raw, v.keyFunc(ctx),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.clientID),
	)
	if err != nil {
		if errors.Is(err, errUnavailable) {
			return errUnavailable
		}
		return errInvalidToken
	}
	return nil
}

// keyFunc возвращает открытый ключ, которым подписан токен.
func (v *Verifier) keyFunc(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, errInvalidToken
		}
		return v.publicKey(ctx, kid)
	}
}

// publicKey находит ключ по kid и один раз обновляет набор, если ключа еще нет.
func (v *Verifier) publicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	if key, ok := v.cachedKey(kid); ok {
		return key, nil
	}
	if err := v.refresh(ctx); err != nil {
		return nil, err
	}
	if key, ok := v.cachedKey(kid); ok {
		return key, nil
	}
	return nil, errInvalidToken
}

// cachedKey возвращает ранее загруженный ключ.
func (v *Verifier) cachedKey(kid string) (*rsa.PublicKey, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	key, ok := v.keys[kid]
	return key, ok
}

// refresh загружает актуальный набор ключей JWKS.
func (v *Verifier) refresh(ctx context.Context) error {
	keys, err := v.fetchKeys(ctx)
	if err != nil {
		return err
	}
	v.mu.Lock()
	v.keys = keys
	v.mu.Unlock()
	return nil
}
