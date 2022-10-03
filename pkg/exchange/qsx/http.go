package qsx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUserAccessDenied = errors.New("you do not have access to the requested resource")
	ErrNotFound         = errors.New("the requested resource not found")
	ErrTooManyRequests  = errors.New("you have exceeded throttle")
)

type RequestError struct {
	//request path
	//method
	//response status code
	//body??
}

type HTTPClient interface {
	Get(ctx context.Context, path string, v interface{}) error
	Post(ctx context.Context, path string, payload interface{}, v interface{}) error
	Put(ctx context.Context, path string, payload interface{}, v interface{}) error
	Delete(ctx context.Context, path string, payload interface{}, v interface{}) error
}

type Options struct {
	ApiURL  string
	Verbose bool
}

type Client struct {
	httpClient  *http.Client
	options     *Options
	timestamp   func() string
	rateLimiter *rate.Limiter
}

func New(httpClient *http.Client, options Options, rateLimiter *rate.Limiter) *Client {
	return &Client{
		httpClient: httpClient,
		options:    &options,
		timestamp: func() string {
			return strconv.FormatInt(time.Now().Unix(), 10)
		},
		rateLimiter: rateLimiter,
	}
}

func (c *Client) Get(ctx context.Context, path string, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", http.MethodGet, err)
	}

	if err = c.doRequest(req, v); err != nil {
		return err
	}

	return nil
}

func (c *Client) Post(ctx context.Context, path string, payload interface{}, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", http.MethodPost, err)
	}

	if err = c.doRequest(req, v); err != nil {
		return err
	}

	return nil
}

func (c *Client) Put(ctx context.Context, path string, payload interface{}, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", http.MethodPut, err)
	}

	if err = c.doRequest(req, v); err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, path string, payload interface{}, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path, payload)
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", http.MethodDelete, err)
	}

	if err = c.doRequest(req, v); err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, payload interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.options.ApiURL, path), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	if c.options.Verbose {
		body, _ := httputil.DumpRequest(req, true)
		log.Println(fmt.Sprintf("%s", string(body)))
	}

	req = req.WithContext(ctx)
	return req, nil
}

func (c *Client) doRequest(r *http.Request, v interface{}) error {
	resp, err := c.do(r)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}
	defer resp.Body.Close()

	if v == nil {
		return nil
	}

	var buf bytes.Buffer
	dec := json.NewDecoder(io.TeeReader(resp.Body, &buf))
	if err = dec.Decode(v); err != nil {
		return fmt.Errorf("could not parse response body: %w [%s:%s] %s", err, r.Method, r.URL.String(), buf.String())
	}

	return nil
}

func (c *Client) do(r *http.Request) (*http.Response, error) {

	// setup rate limiter
	ctx := context.Background()
	err := c.rateLimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to make request [%s:%s]: %w", r.Method, r.URL.String(), err)
	}

	if c.options.Verbose {
		body, _ := httputil.DumpResponse(resp, true)
		log.Println(fmt.Sprintf("%s", string(body)))
	}

	switch resp.StatusCode {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusNoContent:
		return resp, nil
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusUnauthorized,
		http.StatusForbidden:
		return nil, ErrUserAccessDenied
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	}

	return nil, fmt.Errorf("failed to do request, %d status code received", resp.StatusCode)
}

func Query(params []string) string {
	if len(params) == 0 {
		return ""
	}
	return "?" + strings.Join(params, "&")
}

func EncodeURLValues(urlPath string, values url.Values) string {
	u := urlPath
	if len(values) > 0 {
		u += "?" + values.Encode()
	}
	return u
}
