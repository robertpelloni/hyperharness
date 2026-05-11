package providers_test

import (
	"testing"
	"github.com/robertpelloni/hyperharness/internal/providers"
)

func TestGetPooledHTTPClient(t *testing.T) {
	client1 := providers.GetPooledHTTPClient()
	client2 := providers.GetPooledHTTPClient()
	if client1 != client2 {
		t.Errorf("Expected identical pooled clients, got different instances")
	}
}
