package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// MessageType identifies the type of RPC message.
type MessageType string

const (
	TypeRequest  MessageType = "request"
	TypeResponse MessageType = "response"
	TypeNotify   MessageType = "notify"
	TypeError    MessageType = "error"
)

// Message is the wire format for all RPC communication.
type Message struct {
	ID      string          `json:"id"`
	Type    MessageType     `json:"type"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	Source  string          `json:"source,omitempty"`
	Time    time.Time       `json:"time"`
}

// RPCError represents an error in RPC communication.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// Handler processes incoming RPC requests.
type Handler func(ctx context.Context, params json.RawMessage) (json.RawMessage, error)

// Server is an RPC server that listens for connections.
type Server struct {
	listener net.Listener
	handlers map[string]Handler
	clients  map[string]*ClientConn
	mu       sync.RWMutex
	notify   []Handler
	closed   bool
}

// ClientConn represents a connected client.
type ClientConn struct {
	ID       string
	Conn     net.Conn
	encoder  *json.Encoder
	decoder  *json.Decoder
	server   *Server
	mu       sync.Mutex
}

// NewServer creates a new RPC server.
func NewServer() *Server {
	return &Server{
		handlers: make(map[string]Handler),
		clients:  make(map[string]*ClientConn),
	}
}

// Handle registers a handler for a method.
func (s *Server) Handle(method string, handler Handler) {
	s.mu.Lock()
	s.handlers[method] = handler
	s.mu.Unlock()
}

// OnNotify registers a handler for notifications.
func (s *Server) OnNotify(handler Handler) {
	s.mu.Lock()
	s.notify = append(s.notify, handler)
	s.mu.Unlock()
}

// Listen starts listening on the given address.
func (s *Server) Listen(addr string) error {
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go s.acceptLoop()
	return nil
}

// ListenUnix starts listening on a Unix socket.
func (s *Server) ListenUnix(path string) error {
	os.Remove(path) // Remove stale socket

	var err error
	s.listener, err = net.Listen("unix", path)
	if err != nil {
		return err
	}

	go s.acceptLoop()
	return nil
}

// acceptLoop accepts incoming connections.
func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed {
				return
			}
			continue
		}

		client := &ClientConn{
			ID:      fmt.Sprintf("client_%d", time.Now().UnixNano()),
			Conn:    conn,
			encoder: json.NewEncoder(conn),
			decoder: json.NewDecoder(conn),
			server:  s,
		}

		s.mu.Lock()
		s.clients[client.ID] = client
		s.mu.Unlock()

		go client.readLoop()
	}
}

// Close shuts down the server.
func (s *Server) Close() error {
	s.closed = true
	if s.listener != nil {
		s.listener.Close()
	}
	s.mu.Lock()
	for _, client := range s.clients {
		client.Conn.Close()
	}
	s.mu.Unlock()
	return nil
}

// Broadcast sends a notification to all connected clients.
func (s *Server) Broadcast(method string, params interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	msg := Message{
		ID:     fmt.Sprintf("notify_%d", time.Now().UnixNano()),
		Type:   TypeNotify,
		Method: method,
		Params: data,
		Time:   time.Now(),
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		client.send(msg)
	}
	return nil
}

// Clients returns the list of connected client IDs.
func (s *Server) Clients() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ids []string
	for id := range s.clients {
		ids = append(ids, id)
	}
	return ids
}

// readLoop reads messages from a client connection.
func (c *ClientConn) readLoop() {
	defer func() {
		c.server.mu.Lock()
		delete(c.server.clients, c.ID)
		c.server.mu.Unlock()
		c.Conn.Close()
	}()

	for {
		var msg Message
		if err := c.decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		switch msg.Type {
		case TypeRequest:
			c.handleRequest(msg)
		case TypeNotify:
			c.handleNotify(msg)
		}
	}
}

// handleRequest processes an incoming request.
func (c *ClientConn) handleRequest(msg Message) {
	c.server.mu.RLock()
	handler, ok := c.server.handlers[msg.Method]
	c.server.mu.RUnlock()

	if !ok {
		c.send(Message{
			ID:   msg.ID,
			Type: TypeError,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("method not found: %s", msg.Method),
			},
			Time: time.Now(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := handler(ctx, msg.Params)
	if err != nil {
		c.send(Message{
			ID:   msg.ID,
			Type: TypeError,
			Error: &RPCError{
				Code:    -32603,
				Message: err.Error(),
			},
			Time: time.Now(),
		})
		return
	}

	c.send(Message{
		ID:     msg.ID,
		Type:   TypeResponse,
		Result: result,
		Time:   time.Now(),
	})
}

// handleNotify processes an incoming notification.
func (c *ClientConn) handleNotify(msg Message) {
	c.server.mu.RLock()
	handlers := c.server.notify
	c.server.mu.RUnlock()

	ctx := context.Background()
	for _, handler := range handlers {
		handler(ctx, msg.Params)
	}
}

// send sends a message to the client.
func (c *ClientConn) send(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.encoder.Encode(&msg)
}

// ---- Client ----

// Client is an RPC client.
type Client struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
	pending map[string]chan *Message
	mu      sync.Mutex
	nextID  int
	closed  bool
}

// NewClient creates a new RPC client.
func NewClient() *Client {
	return &Client{
		pending: make(map[string]chan *Message),
	}
}

// Connect connects to an RPC server.
func (c *Client) Connect(addr string) error {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)
	go c.readLoop()
	return nil
}

// ConnectUnix connects to a Unix socket RPC server.
func (c *Client) ConnectUnix(path string) error {
	conn, err := net.DialTimeout("unix", path, 10*time.Second)
	if err != nil {
		return err
	}
	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)
	go c.readLoop()
	return nil
}

// Close closes the client connection.
func (c *Client) Close() error {
	c.closed = true
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Call sends a request and waits for the response.
func (c *Client) Call(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.nextID++
	id := fmt.Sprintf("req_%d", c.nextID)
	ch := make(chan *Message, 1)
	c.pending[id] = ch
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
	}()

	msg := Message{
		ID:     id,
		Type:   TypeRequest,
		Method: method,
		Params: data,
		Time:   time.Now(),
	}

	if err := c.encoder.Encode(&msg); err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		if resp.Error != nil {
			return nil, fmt.Errorf("RPC error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Notify sends a notification (no response expected).
func (c *Client) Notify(method string, params interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	msg := Message{
		ID:     fmt.Sprintf("notify_%d", time.Now().UnixNano()),
		Type:   TypeNotify,
		Method: method,
		Params: data,
		Time:   time.Now(),
	}

	return c.encoder.Encode(&msg)
}

// readLoop reads responses from the server.
func (c *Client) readLoop() {
	for {
		var msg Message
		if err := c.decoder.Decode(&msg); err != nil {
			if c.closed || err == io.EOF {
				return
			}
			continue
		}

		c.mu.Lock()
		ch, ok := c.pending[msg.ID]
		c.mu.Unlock()

		if ok {
			ch <- &msg
		}
	}
}
