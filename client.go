package gotsw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

const (
	// SessionTokenHeader is the header to use for authentication.
	SessionTokenHeader = "Authorization"
)

// loggableMimeTypes is a list of MIME types that are safe to log
// the output of. This is useful for debugging or testing.
var loggableMimeTypes = map[string]struct{}{
	"application/json": {},
	"text/plain":       {},
	// lots of webserver error pages are HTML
	"text/html": {},
}

// New creates a Coder client for the provided URL.
func New(auth string) *Client {
	serverURL, err := url.Parse("https://api.tsw.io/v2/")
	if err != nil {
		panic(err)
	}

	return &Client{
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		authorization: auth,
		URL:           serverURL,
		HTTPClient:    &http.Client{},
	}
}

// Client is an HTTP caller for methods to the Coder API.
// @typescript-ignore Client
type Client struct {
	// mu protects the fields sessionToken, logger, and logBodies. These
	// need to be safe for concurrent access.
	mu            sync.RWMutex
	authorization string
	logger        *slog.Logger
	logBodies     bool

	HTTPClient *http.Client
	URL        *url.URL

	// SessionTokenHeader is an optional custom header to use for setting tokens. By
	// default 'Coder-Session-Token' is used.
	SessionTokenHeader string

	// PlainLogger may be set to log HTTP traffic in a human-readable form.
	// It uses the LogBodies option.
	PlainLogger io.Writer
}

func (c *Client) SetClient(client *http.Client) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.HTTPClient = client
	return c
}

// Logger returns the logger for the client.
func (c *Client) Logger() *slog.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// SetLogger sets the logger for the client.
func (c *Client) SetLogger(logger *slog.Logger) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
	return c
}

// SetPlainLogger may be set to log HTTP traffic in a human-readable form.  It
// uses the LogBodies option.
func (c *Client) SetPlainLogger(plainLogger io.Writer) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.PlainLogger = plainLogger
	return c
}

// LogBodies returns whether requests and response bodies are logged.
func (c *Client) LogBodies() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logBodies
}

// SetLogBodies sets whether to log request and response bodies.
func (c *Client) SetLogBodies(logBodies bool) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logBodies = logBodies
	return c
}

// SessionToken returns the currently set token for the client.
func (c *Client) Authorization() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.authorization
}

// SetSessionToken returns the currently set token for the client.
func (c *Client) SetAuthorization(token string) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.authorization = token
	return c
}

func prefixLines(prefix, s []byte) []byte {
	ss := bytes.NewBuffer(make([]byte, 0, len(s)*2))
	for _, line := range bytes.Split(s, []byte("\n")) {
		_, _ = ss.Write(prefix)
		_, _ = ss.Write(line)
		_ = ss.WriteByte('\n')
	}
	return ss.Bytes()
}

// Request performs a HTTP request with the body provided. The caller is
// responsible for closing the response body.
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) (*http.Response, error) {
	logger := c.Logger()
	if ctx == nil {
		return nil, fmt.Errorf("context should not be nil")
	}

	serverURL, err := c.URL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	var r io.Reader
	if body != nil {
		switch data := body.(type) {
		case io.Reader:
			r = data
		case []byte:
			r = bytes.NewReader(data)
		default:
			// Assume JSON in all other cases.
			buf := bytes.NewBuffer(nil)
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			err = enc.Encode(body)
			if err != nil {
				return nil, fmt.Errorf("encode body: %w", err)
			}
			r = buf
		}
	}

	// Copy the request body so we can log it.
	var reqBody []byte
	c.mu.RLock()
	logBodies := c.logBodies
	c.mu.RUnlock()
	if r != nil && logBodies {
		reqBody, err = io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
		r = bytes.NewReader(reqBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, serverURL.String(), r)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "gotsw/v2")

	tokenHeader := c.SessionTokenHeader
	if tokenHeader == "" {
		tokenHeader = SessionTokenHeader
	}
	req.Header.Set(tokenHeader, "Bearer "+c.Authorization())

	if r != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, opt := range opts {
		opt(req)
	}

	// We already capture most of this information in the span (minus
	// the request body which we don't want to capture anyways).
	logger = logger.With(
		"method", req.Method,
		"url", req.URL.String(),
	)
	logger.Debug("sdk request", "body", string(reqBody))

	resp, err := c.HTTPClient.Do(req)

	// We log after sending the request because the HTTP Transport may modify
	// the request within Do, e.g. by adding headers.
	if resp != nil && c.PlainLogger != nil {
		out, err := httputil.DumpRequest(resp.Request, logBodies)
		if err != nil {
			return nil, fmt.Errorf("dump request: %w", err)
		}
		out = prefixLines([]byte("http --> "), out)
		_, _ = c.PlainLogger.Write(out)
	}

	if err != nil {
		return nil, err
	}

	if c.PlainLogger != nil {
		out, err := httputil.DumpResponse(resp, logBodies)
		if err != nil {
			return nil, fmt.Errorf("dump response: %w", err)
		}
		out = prefixLines([]byte("http <-- "), out)
		_, _ = c.PlainLogger.Write(out)
	}

	// Copy the response body so we can log it if it's a loggable mime type.
	var respBody []byte
	if resp.Body != nil && logBodies {
		mimeType := parseMimeType(resp.Header.Get("Content-Type"))
		if _, ok := loggableMimeTypes[mimeType]; ok {
			respBody, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("copy response body for logs: %w", err)
			}
			err = resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("close response body: %w", err)
			}
			resp.Body = io.NopCloser(bytes.NewReader(respBody))
		}
	}

	logger.Debug("sdk response",
		"status", resp.StatusCode,
		"body", string(respBody),
	)

	return resp, err
}

func parseMimeType(contentType string) string {
	mimeType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		mimeType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	}

	return mimeType
}

// RequestOption is a function that can be used to modify an http.Request.
type RequestOption func(*http.Request)

// WithQueryParam adds a query parameter to the request.
func WithQueryParam(key, value string) RequestOption {
	return func(r *http.Request) {
		if value == "" {
			return
		}
		q := r.URL.Query()
		q.Add(key, value)
		r.URL.RawQuery = q.Encode()
	}
}

// HeaderTransport is a http.RoundTripper that adds some headers to all requests.
type HeaderTransport struct {
	Transport http.RoundTripper
	Header    http.Header
}

var _ http.RoundTripper = &HeaderTransport{}

func (h *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	if h.Transport == nil {
		h.Transport = http.DefaultTransport
	}
	return h.Transport.RoundTrip(req)
}

func (h *HeaderTransport) CloseIdleConnections() {
	type closeIdler interface {
		CloseIdleConnections()
	}
	if tr, ok := h.Transport.(closeIdler); ok {
		tr.CloseIdleConnections()
	}
}
