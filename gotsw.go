package gotsw

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
)

const (
	version     = "1.0.0"
	defaultBase = "https://api.tsw.io"
	userAgent   = "gotsw/" + version
	mediaType   = "application/json"
	rateLimit   = 500 * time.Millisecond
	retryLimit  = 3
)

type Client struct {
	// Http Client used to interact with the TSW API
	client *http.Client

	// BASE URL for APIs
	BaseURL *url.URL

	// User Agent for the client
	UserAgent string

	ImageService  ImageService
	RegionService RegionService
	SshKeyService SshKeyService
	MetalService  MetalService

	// Optional function called after every successful request made to API
	onRequestCompleted RequestCompletionCallback

	// Optional extra HTTP headers to set on every request to the API.
	headers map[string]string
}

type RequestCompletionCallback func(*http.Request, *http.Response)

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string `json:"message"`

	Success bool `json:"success"`
}

func NewFromToken(token string) *Client {
	cleanToken := strings.Trim(strings.TrimSpace(token), "'")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cleanToken})
	return NewClient(oauth2.NewClient(ctx, ts))
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBase)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.ImageService = &ImagesServiceHandler{client: c}
	c.RegionService = &RegionsServiceHandler{client: c}
	c.SshKeyService = &SshKeyServiceHandler{client: c}
	c.MetalService = &MetalServiceHandler{client: c}

	c.headers = make(map[string]string)
	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// SetRequestHeaders sets optional HTTP headers on the client that are
// sent on each HTTP request.
func SetRequestHeaders(headers map[string]string) ClientOpt {
	return func(c *Client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mediaType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// OnRequestCompleted sets the API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	resp, err := DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.Message = string(data)
		}
	}

	return errorResponse
}

func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}
