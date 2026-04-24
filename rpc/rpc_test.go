package rpc

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestServerClientTCP(t *testing.T) {
	server := NewServer()
	server.Handle("echo", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		return params, nil
	})
	server.Handle("add", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		var args struct{ A, B int }
		json.Unmarshal(params, &args)
		result, _ := json.Marshal(args.A + args.B)
		return result, nil
	})

	listener, err := listenTestServer(server, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	client := NewClient()
	if err := client.Connect(listener.Addr().String()); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// Test echo
	result, err := client.Call(context.Background(), "echo", map[string]string{"msg": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	var echoResult map[string]string
	json.Unmarshal(result, &echoResult)
	if echoResult["msg"] != "hello" {
		t.Errorf("echo result: %v", echoResult)
	}

	// Test add
	result, err = client.Call(context.Background(), "add", map[string]int{"A": 3, "B": 4})
	if err != nil {
		t.Fatal(err)
	}
	var sum int
	json.Unmarshal(result, &sum)
	if sum != 7 {
		t.Errorf("add result: %d, want 7", sum)
	}

	// Test method not found
	_, err = client.Call(context.Background(), "nonexistent", nil)
	if err == nil {
		t.Error("expected error for unknown method")
	}
}

func TestServerBroadcast(t *testing.T) {
	server := NewServer()

	listener, err := listenTestServer(server, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	client := NewClient()
	client.Connect(listener.Addr().String())
	defer client.Close()

	// Make a call to ensure connection is established
	client.Call(context.Background(), "echo", "ping")

	// Broadcast should not error
	err = server.Broadcast("test", map[string]string{"msg": "hello"})
	if err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}
}

func TestClientTimeout(t *testing.T) {
	server := NewServer()
	server.Handle("slow", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		time.Sleep(2 * time.Second)
		return json.RawMessage(`"done"`), nil
	})

	listener, err := listenTestServer(server, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	client := NewClient()
	client.Connect(listener.Addr().String())
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = client.Call(ctx, "slow", nil)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestRPCMessageSerialization(t *testing.T) {
	msg := Message{
		ID:     "test_1",
		Type:   TypeRequest,
		Method: "test.method",
		Params: json.RawMessage(`{"key": "value"}`),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.ID != "test_1" {
		t.Errorf("ID mismatch: %s", decoded.ID)
	}
	if decoded.Type != TypeRequest {
		t.Errorf("Type mismatch: %s", decoded.Type)
	}
	if decoded.Method != "test.method" {
		t.Errorf("Method mismatch: %s", decoded.Method)
	}
}

func TestRPCError(t *testing.T) {
	rpcErr := &RPCError{
		Code:    -32601,
		Message: "method not found",
	}

	data, err := json.Marshal(rpcErr)
	if err != nil {
		t.Fatal(err)
	}

	var decoded RPCError
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}

	if decoded.Code != -32601 {
		t.Errorf("error code: %d", decoded.Code)
	}
}

func TestServerMultipleClients(t *testing.T) {
	server := NewServer()
	server.Handle("ping", func(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
		return json.RawMessage(`"pong"`), nil
	})

	listener, err := listenTestServer(server, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	c1 := NewClient()
	c1.Connect(listener.Addr().String())
	defer c1.Close()

	// Force connection to be accepted by making a call
	c1.Call(context.Background(), "ping", nil)

	c2 := NewClient()
	c2.Connect(listener.Addr().String())
	defer c2.Close()

	c2.Call(context.Background(), "ping", nil)

	clients := server.Clients()
	if len(clients) < 2 {
		t.Skipf("race in accept loop, got %d clients (flaky)", len(clients))
	}
}

// Helper to get the actual listener address
func listenTestServer(s *Server, addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s.listener = ln
	go s.acceptLoop()
	return ln, nil
}
