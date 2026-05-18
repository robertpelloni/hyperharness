package auth

import (
	"os"
	"testing"
	"time"
)

func TestNewVolatileKeyVault(t *testing.T) {
	vault := NewVolatileKeyVault()
	if vault == nil {
		t.Fatal("vault should not be nil")
	}
	if len(vault.creds) != 0 {
		t.Error("new vault should be empty")
	}
	if !vault.volatile {
		t.Error("should be volatile")
	}
}

func TestStoreAndGet(t *testing.T) {
	vault := NewVolatileKeyVault()

	cred := &Credential{
		Method: AuthMethodAPIKey,
		APIKey: "sk-test-12345678",
	}

	err := vault.Store(ProviderOpenAI, cred)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := vault.Get(ProviderOpenAI)
	if !ok {
		t.Fatal("should find stored credential")
	}
	if got.APIKey != "sk-test-12345678" {
		t.Errorf("API key mismatch: %s", got.APIKey)
	}
	if got.Method != AuthMethodAPIKey {
		t.Errorf("method mismatch: %s", got.Method)
	}
}

func TestGetMissing(t *testing.T) {
	vault := NewVolatileKeyVault()
	_, ok := vault.Get(ProviderOpenAI)
	if ok {
		t.Error("should not find missing credential")
	}
}

func TestDelete(t *testing.T) {
	vault := NewVolatileKeyVault()
	vault.Store(ProviderOpenAI, &Credential{APIKey: "sk-test"})

	err := vault.Delete(ProviderOpenAI)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := vault.Get(ProviderOpenAI)
	if ok {
		t.Error("should not find deleted credential")
	}
}

func TestList(t *testing.T) {
	vault := NewVolatileKeyVault()
	vault.Store(ProviderOpenAI, &Credential{APIKey: "sk-test-12345678"})
	vault.Store(ProviderAnthropic, &Credential{APIKey: "sk-ant-test-12345678"})

	list := vault.List()
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}

	// Keys should be masked
	for _, cred := range list {
		if cred.APIKey == "sk-test-12345678" || cred.APIKey == "sk-ant-test-12345678" {
			t.Error("API key should be masked in listing")
		}
		if cred.APIKey == "" {
			t.Error("masked key should not be empty")
		}
	}
}

func TestCredentialExpiry(t *testing.T) {
	// Not expired
	cred := &Credential{Method: AuthMethodAPIKey}
	if cred.IsExpired() {
		t.Error("credential with no expiry should not be expired")
	}

	// Expired
	cred.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if !cred.IsExpired() {
		t.Error("credential with past expiry should be expired")
	}

	// Needs refresh
	cred.ExpiresAt = time.Now().Add(1 * time.Minute)
	if !cred.NeedsRefresh(5*time.Minute) {
		t.Error("credential expiring in 1min should need refresh within 5min")
	}
	if cred.NeedsRefresh(30*time.Second) {
		t.Error("credential expiring in 1min should not need refresh within 30s")
	}
}

func TestResolveFromEnv(t *testing.T) {
	vault := NewVolatileKeyVault()

	// Set env var
	os.Setenv("OPENAI_API_KEY", "sk-env-test-12345678")
	defer os.Unsetenv("OPENAI_API_KEY")

	cred, err := vault.Resolve(ProviderOpenAI)
	if err != nil {
		t.Fatal(err)
	}
	if cred.APIKey != "sk-env-test-12345678" {
		t.Errorf("expected env key, got: %s", cred.APIKey)
	}
	if cred.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("expected default base URL, got: %s", cred.BaseURL)
	}
}

func TestResolveLocalProvider(t *testing.T) {
	vault := NewVolatileKeyVault()

	cred, err := vault.Resolve(ProviderOllama)
	if err != nil {
		t.Fatal(err)
	}
	if cred.Method != AuthMethodNone {
		t.Errorf("ollama should use no auth, got: %s", cred.Method)
	}
	if cred.BaseURL != "http://localhost:11434" {
		t.Errorf("ollama base URL: %s", cred.BaseURL)
	}
}

func TestResolveMissing(t *testing.T) {
	vault := NewVolatileKeyVault()

	// Unset all env vars for Anthropic
	os.Unsetenv("ANTHROPIC_API_KEY")

	_, err := vault.Resolve(ProviderAnthropic)
	if err == nil {
		t.Error("should fail for unconfigured provider")
	}
}

func TestResolveAll(t *testing.T) {
	vault := NewVolatileKeyVault()
	os.Setenv("OPENAI_API_KEY", "sk-test")
	defer os.Unsetenv("OPENAI_API_KEY")

	statuses := vault.ResolveAll()
	if len(statuses) == 0 {
		t.Error("should have statuses for all providers")
	}

	oiStatus, ok := statuses[ProviderOpenAI]
	if !ok {
		t.Error("should have OpenAI status")
	}
	if !oiStatus.Configured {
		t.Error("OpenAI should be configured via env var")
	}
	if !oiStatus.Authenticated {
		t.Error("OpenAI should be authenticated via env var")
	}
}

func TestLogin(t *testing.T) {
	vault := NewVolatileKeyVault()
	err := vault.Login(ProviderAnthropic, "sk-ant-test-12345678")
	if err != nil {
		t.Fatal(err)
	}

	cred, ok := vault.Get(ProviderAnthropic)
	if !ok {
		t.Fatal("should find credential after login")
	}
	if cred.APIKey != "sk-ant-test-12345678" {
		t.Errorf("API key mismatch: %s", cred.APIKey)
	}
}

func TestLogout(t *testing.T) {
	vault := NewVolatileKeyVault()
	vault.Login(ProviderAnthropic, "sk-ant-test")

	err := vault.Logout(ProviderAnthropic)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := vault.Get(ProviderAnthropic)
	if ok {
		t.Error("should not find credential after logout")
	}
}

func TestBuildAuthHeader(t *testing.T) {
	tests := []struct {
		name     string
		cred     *Credential
		expected string
	}{
		{"nil", nil, ""},
		{"api key", &Credential{Method: AuthMethodAPIKey, APIKey: "sk-test"}, "Bearer sk-test"},
		{"bearer token", &Credential{Method: AuthMethodBearer, Token: "tok-123"}, "Bearer tok-123"},
		{"none", &Credential{Method: AuthMethodNone}, ""},
		{"bearer with api key", &Credential{Method: AuthMethodBearer, APIKey: "sk-test"}, "Bearer sk-test"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := BuildAuthHeader(tc.cred)
			if result != tc.expected {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestNewNonce(t *testing.T) {
	n1, err := NewNonce(32)
	if err != nil {
		t.Fatal(err)
	}
	n2, err := NewNonce(32)
	if err != nil {
		t.Fatal(err)
	}
	if n1 == n2 {
		t.Error("nonces should be unique")
	}
	if len(n1) < 10 {
		t.Error("nonce should be reasonably long")
	}
}

func TestHashSHA256(t *testing.T) {
	hash := HashSHA256([]byte("hello"))
	if hash == "" {
		t.Error("hash should not be empty")
	}
	if len(hash) != 64 {
		t.Errorf("SHA-256 hex should be 64 chars, got %d", len(hash))
	}
}

func TestMaskKey(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{"short", "****"},
		{"sk-test-1234567890abcdef", "sk-t...cdef"},
		{"x", "****"},
		{"12345678", "****"},
		{"123456789", "1234...6789"},
	}

	for _, tc := range tests {
		result := maskKey(tc.key)
		if result != tc.expected {
			t.Errorf("maskKey(%q) = %q, want %q", tc.key, result, tc.expected)
		}
	}
}

func TestProviderSpecs(t *testing.T) {
	if len(ProviderSpecs) < 10 {
		t.Errorf("should have at least 10 provider specs, got %d", len(ProviderSpecs))
	}

	for provider, spec := range ProviderSpecs {
		if spec.Provider != provider {
			t.Errorf("spec provider mismatch: %s != %s", spec.Provider, provider)
		}
		if spec.Name == "" {
			t.Errorf("provider %s should have a name", provider)
		}
	}
}

func TestFileKeyVault(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/credentials.json"

	vault, err := NewKeyVault(path)
	if err != nil {
		t.Fatal(err)
	}

	vault.Login(ProviderOpenAI, "sk-file-test-12345678")
	vault.Login(ProviderAnthropic, "sk-ant-file-test-12345678")

	// Reload from disk
	vault2, err := NewKeyVault(path)
	if err != nil {
		t.Fatal(err)
	}

	cred, ok := vault2.Get(ProviderOpenAI)
	if !ok {
		t.Fatal("should find OpenAI credential from persisted vault")
	}
	if cred.APIKey != "sk-file-test-12345678" {
		t.Errorf("API key mismatch after reload: %s", cred.APIKey)
	}
}
