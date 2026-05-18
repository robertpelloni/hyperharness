package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Message types
const (
	TypeRequest  = "request"
	TypeResponse = "response"
	TypeNotify   = "notify"
)

// Message represents an RPC message.
type Message struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *RPCError       `json:"error,omitempty"`
}

// RPCError represents an error in an RPC response.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Handler processes an RPC request.
type Handler func(ctx context.Context, params json.RawMessage) (json.RawMessage, error)

// Server is an RPC server.
type Server struct {
	mu       sync.RWMutex
	handlers map[string]Handler
	clients  map[net.Conn]struct{}
	closed   atomic.Bool

	listener net.Listener
}

// NewServer creates a new RPC server.
func NewServer() *Server {
	return &Server{
		handlers: make(map[string]Handler),
		clients:  make(map[net.Conn]struct{}),
	}
}

// Handle registers a handler for a method.
func (s *Server) Handle(method string, handler Handler) {
	s.mu.Lock()
	s.handlers[method] = handler
	s.mu.Unlock()
}

// Close shuts down the server.
func (s *Server) Close() error {
	s.closed.Store(true)
	s.mu.Lock()
	for c := range s.clients {
		c.Close()
	}
	s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// Clients returns the list of connected client connections.
func (s *Server) Clients() []net.Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]net.Conn, 0, len(s.clients))
	for c := range s.clients {
		result = append(result, c)
	}
	return result
}

// Broadcast sends a notification to all connected clients.
func (s *Server) Broadcast(method string, params interface{}) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	paramData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	msg := Message{
		ID:     fmt.Sprintf("broadcast_%d", time.Now().UnixNano()),
		Type:   TypeNotify,
		Method: method,
		Params: paramData,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	for c := range s.clients {
		c.Write(data)
	}
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}

		s.mu.Lock()
		s.clients[conn] = struct{}{}
		s.mu.Unlock()

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	dec := json.NewDecoder(conn)
	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			if s.closed.Load() {
				return
			}
			return
		}

		go s.processMessage(conn, msg)
	}
}

func (s *Server) processMessage(conn net.Conn, msg Message) {
	s.mu.RLock()
	handler, ok := s.handlers[msg.Method]
	s.mu.RUnlock()

	resp := Message{
		ID:   msg.ID,
		Type: TypeResponse,
	}

	if !ok {
		resp.Error = &RPCError{Code: -32601, Message: "method not found"}
	} else {
		result, err := handler(context.Background(), msg.Params)
		if err != nil {
			resp.Error = &RPCError{Code: -32603, Message: err.Error()}
		} else {
			resp.Result = result
		}
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	data = append(data, '\n')
	conn.Write(data)
}

// Client is an RPC client.
type Client struct {
	conn    net.Conn
	mu      sync.Mutex
	nextID  atomic.Uint64
	pending map[string]chan *Message
	closed  atomic.Bool
}

// NewClient creates a new RPC client.
func NewClient() *Client {
	return &Client{
		pending: make(map[string]chan *Message),
	}
}

// Connect dials the server at the given address.
func (c *Client) Connect(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readLoop()
	return nil
}

// Close shuts down the client.
func (c *Client) Close() error {
	c.closed.Store(true)
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Call sends an RPC request and waits for the response.
func (c *Client) Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	id := fmt.Sprintf("%d", c.nextID.Add(1))

	paramData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	msg := Message{
		ID:     id,
		Type:   TypeRequest,
		Method: method,
		Params: paramData,
	}

	ch := make(chan *Message, 1)
	c.mu.Lock()
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')

	c.mu.Lock()
	_, err = c.conn.Write(data)
	c.mu.Unlock()
	if err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		if resp.Error != nil {
			return nil, errors.New(resp.Error.Message)
		}
		return resp.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) readLoop() {
	dec := json.NewDecoder(c.conn)
	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			if c.closed.Load() {
				return
			}
			return
		}

		c.mu.Lock()
		ch, ok := c.pending[msg.ID]
		c.mu.Unlock()

		if ok {
			ch <- &msg
		}
	}
}
