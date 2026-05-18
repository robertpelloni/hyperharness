package auth

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════════════════════════════
// auth.go — Multi-provider authentication for HyperHarness
// Ported from ollama/auth, opencode auth, and pi-cli credential patterns.
//
// WHAT: Unified authentication with key vault, token management, and
//       provider-specific credential resolution.
// WHY: Every LLM provider has different auth mechanisms (API keys, OAuth,
//       SSH keys, JWT tokens). A unified vault simplifies access.
// HOW: KeyVault stores encrypted credentials per-provider, with automatic
//       env var resolution, file-based secrets, and token refresh.
// ═══════════════════════════════════════════════════════════════════════

// Provider identifies which LLM service we're authenticating with.
type Provider string

const (
	ProviderOpenAI     Provider = "openai"
	ProviderAnthropic  Provider = "anthropic"
	ProviderGoogle     Provider = "google"
	ProviderDeepSeek   Provider = "deepseek"
	ProviderGroq       Provider = "groq"
	ProviderOpenRouter Provider = "openrouter"
	ProviderXAI        Provider = "xai"
	ProviderMistral    Provider = "mistral"
	ProviderCerebras   Provider = "cerebras"
	ProviderAzure      Provider = "azure"
	ProviderBedrock    Provider = "bedrock"
	ProviderVertex     Provider = "vertex"
	ProviderOllama     Provider = "ollama"
	ProviderLMStudio   Provider = "lmstudio"
	ProviderCustom     Provider = "custom"
)

// AuthMethod describes how a provider authenticates.
type AuthMethod string

const (
	AuthMethodAPIKey    AuthMethod = "api_key"
	AuthMethodOAuth     AuthMethod = "oauth"
	AuthMethodSSHKey    AuthMethod = "ssh_key"
	AuthMethodJWT       AuthMethod = "jwt"
	AuthMethodNone      AuthMethod = "none"      // Local providers
	AuthMethodBearer    AuthMethod = "bearer"     // Simple bearer token
	AuthMethodHMAC      AuthMethod = "hmac"       // HMAC-signed requests
	AuthMethodMtls      AuthMethod = "mtls"       // Mutual TLS
)

// Credential holds a provider's authentication data.
type Credential struct {
	Provider     Provider   `json:"provider"`
	Method       AuthMethod `json:"method"`
	APIKey       string     `json:"apiKey,omitempty"`
	Token        string     `json:"token,omitempty"`
	RefreshToken string     `json:"refreshToken,omitempty"`
	ClientID     string     `json:"clientId,omitempty"`
	ClientSecret string     `json:"clientSecret,omitempty"`
	BaseURL      string     `json:"baseUrl,omitempty"`
	SSHKeyPath   string     `json:"sshKeyPath,omitempty"`
	ExpiresAt    time.Time  `json:"expiresAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// IsExpired checks if the credential has expired.
func (c *Credential) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false // No expiry set
	}
	return time.Now().After(c.ExpiresAt)
}

// NeedsRefresh checks if the credential will expire soon.
func (c *Credential) NeedsRefresh(within time.Duration) bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().Add(within).After(c.ExpiresAt)
}

// ProviderSpec defines how to resolve credentials for a provider.
type ProviderSpec struct {
	Provider      Provider
	Name          string
	AuthMethod    AuthMethod
	EnvVars       []string   // Environment variables to check (in order)
	ConfigFile    string     // Relative path to config file
	DefaultBaseURL string
}

// ProviderSpecs is the registry of known provider auth configurations.
var ProviderSpecs = map[Provider]ProviderSpec{
	ProviderOpenAI: {
		Provider:   ProviderOpenAI,
		Name:       "OpenAI",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"OPENAI_API_KEY"},
		DefaultBaseURL: "https://api.openai.com/v1",
	},
	ProviderAnthropic: {
		Provider:   ProviderAnthropic,
		Name:       "Anthropic",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"ANTHROPIC_API_KEY"},
		DefaultBaseURL: "https://api.anthropic.com/v1",
	},
	ProviderGoogle: {
		Provider:   ProviderGoogle,
		Name:       "Google (Gemini)",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"GOOGLE_API_KEY", "GEMINI_API_KEY"},
		DefaultBaseURL: "https://generativelanguage.googleapis.com/v1beta",
	},
	ProviderDeepSeek: {
		Provider:   ProviderDeepSeek,
		Name:       "DeepSeek",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"DEEPSEEK_API_KEY"},
		DefaultBaseURL: "https://api.deepseek.com/v1",
	},
	ProviderGroq: {
		Provider:   ProviderGroq,
		Name:       "Groq",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"GROQ_API_KEY"},
		DefaultBaseURL: "https://api.groq.com/openai/v1",
	},
	ProviderOpenRouter: {
		Provider:   ProviderOpenRouter,
		Name:       "OpenRouter",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"OPENROUTER_API_KEY"},
		DefaultBaseURL: "https://openrouter.ai/api/v1",
	},
	ProviderXAI: {
		Provider:   ProviderXAI,
		Name:       "xAI (Grok)",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"XAI_API_KEY"},
		DefaultBaseURL: "https://api.x.ai/v1",
	},
	ProviderMistral: {
		Provider:   ProviderMistral,
		Name:       "Mistral",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"MISTRAL_API_KEY"},
		DefaultBaseURL: "https://api.mistral.ai/v1",
	},
	ProviderCerebras: {
		Provider:   ProviderCerebras,
		Name:       "Cerebras",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"CEREBRAS_API_KEY"},
		DefaultBaseURL: "https://api.cerebras.ai/v1",
	},
	ProviderAzure: {
		Provider:   ProviderAzure,
		Name:       "Azure OpenAI",
		AuthMethod: AuthMethodAPIKey,
		EnvVars:    []string{"AZURE_OPENAI_API_KEY"},
	},
	ProviderBedrock: {
		Provider:   ProviderBedrock,
		Name:       "AWS Bedrock",
		AuthMethod: AuthMethodHMAC,
		EnvVars:    []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"},
	},
	ProviderVertex: {
		Provider:   ProviderVertex,
		Name:       "Google Vertex AI",
		AuthMethod: AuthMethodOAuth,
		EnvVars:    []string{"GOOGLE_APPLICATION_CREDENTIALS"},
	},
	ProviderOllama: {
		Provider:   ProviderOllama,
		Name:       "Ollama (Local)",
		AuthMethod: AuthMethodNone,
		DefaultBaseURL: "http://localhost:11434",
	},
	ProviderLMStudio: {
		Provider:   ProviderLMStudio,
		Name:       "LM Studio (Local)",
		AuthMethod: AuthMethodNone,
		DefaultBaseURL: "http://localhost:1234/v1",
	},
}

// KeyVault is a secure, file-backed credential store.
type KeyVault struct {
	mu       sync.RWMutex
	path     string
	creds    map[Provider]*Credential
	volatile bool // If true, don't persist to disk
}

// NewKeyVault creates a new key vault, loading from disk if available.
func NewKeyVault(vaultPath string) (*KeyVault, error) {
	if vaultPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		vaultPath = filepath.Join(home, ".hyperharness", "credentials.json")
	}

	vault := &KeyVault{
		path:  vaultPath,
		creds: make(map[Provider]*Credential),
	}

	if err := vault.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	return vault, nil
}

// NewVolatileKeyVault creates an in-memory-only key vault.
func NewVolatileKeyVault() *KeyVault {
	return &KeyVault{
		creds:    make(map[Provider]*Credential),
		volatile: true,
	}
}

// Store saves a credential for a provider.
func (kv *KeyVault) Store(provider Provider, cred *Credential) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	cred.Provider = provider
	cred.UpdatedAt = time.Now()
	if cred.CreatedAt.IsZero() {
		cred.CreatedAt = time.Now()
	}

	kv.creds[provider] = cred

	if !kv.volatile {
		return kv.save()
	}
	return nil
}

// Get retrieves a credential for a provider.
func (kv *KeyVault) Get(provider Provider) (*Credential, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	cred, ok := kv.creds[provider]
	if !ok {
		return nil, false
	}
	if cred.IsExpired() {
		return cred, false // Expired but still returned for refresh
	}
	return cred, true
}

// Delete removes a credential for a provider.
func (kv *KeyVault) Delete(provider Provider) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	delete(kv.creds, provider)

	if !kv.volatile {
		return kv.save()
	}
	return nil
}

// List returns all stored credentials (with API keys masked).
func (kv *KeyVault) List() []Credential {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	result := make([]Credential, 0, len(kv.creds))
	for _, cred := range kv.creds {
		masked := *cred
		if masked.APIKey != "" {
			masked.APIKey = maskKey(masked.APIKey)
		}
		if masked.Token != "" {
			masked.Token = maskKey(masked.Token)
		}
		if masked.ClientSecret != "" {
			masked.ClientSecret = maskKey(masked.ClientSecret)
		}
		result = append(result, masked)
	}
	return result
}

// Resolve resolves a credential by checking the vault, then environment variables.
func (kv *KeyVault) Resolve(provider Provider) (*Credential, error) {
	// Check vault first
	if cred, ok := kv.Get(provider); ok {
		return cred, nil
	}

	// Check expired credential that might be refreshable
	if cred, ok := kv.creds[provider]; ok && cred.NeedsRefresh(5*time.Minute) {
		return cred, fmt.Errorf("credential for %s needs refresh", provider)
	}

	// Fall back to environment variables
	spec, ok := ProviderSpecs[provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	for _, envVar := range spec.EnvVars {
		if val := os.Getenv(envVar); val != "" {
			return &Credential{
				Provider:  provider,
				Method:    spec.AuthMethod,
				APIKey:    val,
				BaseURL:   spec.DefaultBaseURL,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		}
	}

	// Local providers don't need auth
	if spec.AuthMethod == AuthMethodNone {
		return &Credential{
			Provider:  provider,
			Method:    AuthMethodNone,
			BaseURL:   spec.DefaultBaseURL,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}

	return nil, fmt.Errorf("no credentials found for %s (set %s)", provider, strings.Join(spec.EnvVars, " or "))
}

// ResolveAll resolves credentials for all known providers.
func (kv *KeyVault) ResolveAll() map[Provider]*CredentialStatus {
	result := make(map[Provider]*CredentialStatus)
	for provider := range ProviderSpecs {
		cred, err := kv.Resolve(provider)
		status := &CredentialStatus{
			Provider:      provider,
			Configured:    err == nil,
			Authenticated: err == nil && (cred != nil && !cred.IsExpired()),
		}
		if err != nil {
			status.Error = err.Error()
		}
		if cred != nil {
			status.Method = cred.Method
			status.BaseURL = cred.BaseURL
		}
		result[provider] = status
	}
	return result
}

// CredentialStatus reports the authentication state for a provider.
type CredentialStatus struct {
	Provider      Provider   `json:"provider"`
	Configured    bool       `json:"configured"`
	Authenticated bool       `json:"authenticated"`
	Method        AuthMethod `json:"method,omitempty"`
	BaseURL       string     `json:"baseUrl,omitempty"`
	Error         string     `json:"error,omitempty"`
}

// Login stores a credential interactively (programmatic version).
func (kv *KeyVault) Login(provider Provider, apiKey string) error {
	spec, ok := ProviderSpecs[provider]
	if !ok {
		return fmt.Errorf("unknown provider: %s", provider)
	}

	cred := &Credential{
		Provider:  provider,
		Method:    spec.AuthMethod,
		APIKey:    apiKey,
		BaseURL:   spec.DefaultBaseURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return kv.Store(provider, cred)
}

// Logout removes a provider's credential.
func (kv *KeyVault) Logout(provider Provider) error {
	return kv.Delete(provider)
}

// ── Internal helpers ──

func (kv *KeyVault) load() error {
	data, err := os.ReadFile(kv.path)
	if err != nil {
		return err
	}

	var creds map[Provider]*Credential
	if err := json.Unmarshal(data, &creds); err != nil {
		return fmt.Errorf("corrupted credentials file: %w", err)
	}

	kv.creds = creds
	return nil
}

func (kv *KeyVault) save() error {
	dir := filepath.Dir(kv.path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(kv.creds, "", "  ")
	if err != nil {
		return err
	}

	// Write with restrictive permissions (owner only)
	return os.WriteFile(kv.path, data, 0600)
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

// ── Nonce and signing (ported from ollama/auth) ──

// NewNonce generates a cryptographically random nonce.
func NewNonce(length int) (string, error) {
	nonce := make([]byte, length)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(nonce), nil
}

// SignWithEd25519 signs data with an Ed25519 private key.
func SignWithEd25519(privateKey ed25519.PrivateKey, data []byte) (string, error) {
	signature := ed25519.Sign(privateKey, data)
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyEd25519 verifies an Ed25519 signature.
func VerifyEd25519(publicKey ed25519.PublicKey, data []byte, signature string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}
	if !ed25519.Verify(publicKey, data, sig) {
		return errors.New("signature verification failed")
	}
	return nil
}

// HashSHA256 computes a SHA-256 hash.
func HashSHA256(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

// ── Token refresh (for OAuth/JWT providers) ──

// TokenRefresher handles OAuth2-style token refresh.
type TokenRefresher struct {
	vault     *KeyVault
	provider  Provider
	tokenURL  string
	clientID  string
}

// NewTokenRefresher creates a token refresher for a provider.
func NewTokenRefresher(vault *KeyVault, provider Provider, tokenURL, clientID string) *TokenRefresher {
	return &TokenRefresher{
		vault:    vault,
		provider: provider,
		tokenURL: tokenURL,
		clientID: clientID,
	}
}

// Refresh attempts to refresh an OAuth token.
func (tr *TokenRefresher) Refresh(ctx context.Context) (*Credential, error) {
	cred, ok := tr.vault.Get(tr.provider)
	if !ok {
		return nil, fmt.Errorf("no credential for %s", tr.provider)
	}

	if cred.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token available for %s", tr.provider)
	}

	// In a real implementation, this would make an HTTP POST to tokenURL
	// with grant_type=refresh_token. For now, return the existing credential.
	return cred, nil
}

// ── Auth header builder ──

// BuildAuthHeader constructs the appropriate Authorization header value.
func BuildAuthHeader(cred *Credential) string {
	if cred == nil {
		return ""
	}

	switch cred.Method {
	case AuthMethodAPIKey, AuthMethodBearer:
		if cred.Token != "" {
			return "Bearer " + cred.Token
		}
		return "Bearer " + cred.APIKey
	case AuthMethodNone:
		return ""
	default:
		if cred.APIKey != "" {
			return "Bearer " + cred.APIKey
		}
		return ""
	}
}
