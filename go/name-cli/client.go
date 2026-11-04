package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client talks to the name.com Core API (https://docs.name.com/api/v1).
type Client struct {
	BaseURL  string
	Username string
	Token    string
	HTTP     *http.Client
}

func NewClient(baseURL, username, token string) *Client {
	return &Client{
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Username: username,
		Token:    token,
		HTTP:     &http.Client{Timeout: 30 * time.Second},
	}
}

type Record struct {
	ID         int    `json:"id,omitempty"`
	DomainName string `json:"domainName,omitempty"`
	Host       string `json:"host"`
	FQDN       string `json:"fqdn,omitempty"`
	Type       string `json:"type"`
	Answer     string `json:"answer"`
	TTL        int    `json:"ttl,omitempty"`
	Priority   int    `json:"priority,omitempty"`
}

type Domain struct {
	DomainName       string   `json:"domainName"`
	CreateDate       string   `json:"createDate,omitempty"`
	ExpireDate       string   `json:"expireDate,omitempty"`
	AutorenewEnabled bool     `json:"autorenewEnabled"`
	Locked           bool     `json:"locked"`
	PrivacyEnabled   bool     `json:"privacyEnabled"`
	Nameservers      []string `json:"nameservers,omitempty"`
}

type Availability struct {
	DomainName    string  `json:"domainName"`
	Purchasable   bool    `json:"purchasable"`
	Premium       bool    `json:"premium,omitempty"`
	PurchasePrice float64 `json:"purchasePrice,omitempty"`
	RenewalPrice  float64 `json:"renewalPrice,omitempty"`
	PurchaseType  string  `json:"purchaseType,omitempty"`
	Reason        string  `json:"reason,omitempty"`
}

// APIError is a non-2xx response from the API.
type APIError struct {
	Status  int
	Method  string
	Path    string
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e *APIError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = http.StatusText(e.Status)
	}
	s := fmt.Sprintf("%s %s: %d %s", e.Method, e.Path, e.Status, msg)
	if e.Details != "" {
		s += " (" + e.Details + ")"
	}
	return s
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	u := c.BaseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reqBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		apiErr := &APIError{Status: resp.StatusCode, Method: method, Path: path}
		_ = json.Unmarshal(respBody, apiErr) // best effort; body may not be JSON
		return apiErr
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decoding response from %s %s: %w", method, path, err)
	}
	return nil
}

func domainPath(domain string) string {
	return "/core/v1/domains/" + url.PathEscape(domain)
}

// ListRecords returns every record in the zone, following pagination.
func (c *Client) ListRecords(ctx context.Context, domain string) ([]Record, error) {
	var all []Record
	page := 1
	for {
		var resp struct {
			Records  []Record `json:"records"`
			NextPage int      `json:"nextPage"`
		}
		q := url.Values{"perPage": {"1000"}, "page": {strconv.Itoa(page)}}
		if err := c.do(ctx, http.MethodGet, domainPath(domain)+"/records", q, nil, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Records...)
		if resp.NextPage == 0 || resp.NextPage == page {
			return all, nil
		}
		page = resp.NextPage
	}
}

func (c *Client) GetRecord(ctx context.Context, domain string, id int) (Record, error) {
	var r Record
	err := c.do(ctx, http.MethodGet, domainPath(domain)+"/records/"+strconv.Itoa(id), nil, nil, &r)
	return r, err
}

// CreateRecord POSTs to the records collection (no ID in the path).
func (c *Client) CreateRecord(ctx context.Context, domain string, r Record) (Record, error) {
	r.ID, r.FQDN, r.DomainName = 0, "", "" // ignored by the API on create
	var created Record
	err := c.do(ctx, http.MethodPost, domainPath(domain)+"/records", nil, r, &created)
	return created, err
}

// UpdateRecord replaces record id with r (PUT, so send every field; see
// recordsUpdate for read-modify-write).
func (c *Client) UpdateRecord(ctx context.Context, domain string, id int, r Record) (Record, error) {
	r.ID, r.FQDN, r.DomainName = 0, "", "" // the ID goes in the path; the rest is read-only
	var updated Record
	err := c.do(ctx, http.MethodPut, domainPath(domain)+"/records/"+strconv.Itoa(id), nil, r, &updated)
	return updated, err
}

// ListDomains returns every domain in the account, following pagination.
func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	var all []Domain
	page := 1
	for {
		var resp struct {
			Domains  []Domain `json:"domains"`
			NextPage int      `json:"nextPage"`
		}
		q := url.Values{"perPage": {"1000"}, "page": {strconv.Itoa(page)}}
		if err := c.do(ctx, http.MethodGet, "/core/v1/domains", q, nil, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Domains...)
		if resp.NextPage == 0 || resp.NextPage == page {
			return all, nil
		}
		page = resp.NextPage
	}
}

func (c *Client) GetDomain(ctx context.Context, domain string) (Domain, error) {
	var d Domain
	err := c.do(ctx, http.MethodGet, domainPath(domain), nil, nil, &d)
	return d, err
}

// CheckAvailability checks up to 50 domains for registration.
func (c *Client) CheckAvailability(ctx context.Context, domains []string) ([]Availability, error) {
	body := map[string]any{"domainNames": domains, "purchaseType": "registration"}
	var resp struct {
		Results []Availability `json:"results"`
	}
	// The ":" must not be percent-encoded, per the API docs.
	err := c.do(ctx, http.MethodPost, "/core/v1/domains:checkAvailability", nil, body, &resp)
	return resp.Results, err
}
