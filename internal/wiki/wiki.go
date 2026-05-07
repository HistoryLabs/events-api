// Package wiki provides a thin client for the Wikipedia action=parse API.
package wiki

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL   = "https://en.wikipedia.org/w/api.php"
	defaultUserAgent = "EventsAPI/1.0 (https://github.com/HistoryLabs/events-api)"
	defaultTimeout   = 10 * time.Second
)

var ErrUpstreamStatus = errors.New("wiki: non-2xx upstream response")

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	UserAgent  string
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		BaseURL:    defaultBaseURL,
		UserAgent:  defaultUserAgent,
	}
}

var DefaultClient = NewClient()

// FetchOpts configures a single Fetch call.
// Section is the Wikipedia section number (1 = first content section).
// Redirects, when true, asks the API to follow page redirects.
type FetchOpts struct {
	Section   int
	Redirects bool
}

// Fetch calls action=parse for the given page and returns the raw JSON body.
// The body is returned undecoded because callers regex over the JSON-escaped HTML directly.
func (c *Client) Fetch(ctx context.Context, page string, opts FetchOpts) ([]byte, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("wiki: parse base url: %w", err)
	}
	q := url.Values{}
	q.Set("action", "parse")
	q.Set("format", "json")
	q.Set("page", page)
	q.Set("section", strconv.Itoa(opts.Section))
	if opts.Redirects {
		q.Set("redirects", "true")
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("wiki: build request: %w", err)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wiki: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: %d", ErrUpstreamStatus, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("wiki: read body: %w", err)
	}
	return body, nil
}

// Convenience wrapper around DefaultClient.Fetch.
func Fetch(ctx context.Context, page string, opts FetchOpts) ([]byte, error) {
	return DefaultClient.Fetch(ctx, page, opts)
}
