package vault

import (
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretMeta holds metadata about a Vault secret relevant to expiration.
type SecretMeta struct {
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
}

// Client wraps the Vault API client.
type Client struct {
	vc *vaultapi.Client
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}
	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// GetSecretMeta reads a KV v2 secret and returns its expiration metadata.
func (c *Client) GetSecretMeta(path string) (*SecretMeta, error) {
	secret, err := c.vc.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	ttlRaw, ok := secret.Data["ttl"]
	if !ok {
		return nil, fmt.Errorf("no ttl field in secret at %q", path)
	}

	ttlFloat, ok := ttlRaw.(float64)
	if !ok {
		return nil, fmt.Errorf("unexpected ttl type at %q", path)
	}

	ttl := time.Duration(ttlFloat) * time.Second
	expiresAt := time.Now().Add(ttl)

	return &SecretMeta{
		Path:      path,
		ExpiresAt: expiresAt,
		TTL:       ttl,
	}, nil
}
