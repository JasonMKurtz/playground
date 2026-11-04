package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// fakeAPI records the last request and replies with a canned status/body.
type fakeAPI struct {
	status   int
	body     string
	pages    map[string]string // optional: response body keyed by ?page=
	method   string
	path     string
	reqBody  map[string]any
	user     string
	password string
}

func (f *fakeAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.method = r.Method
	f.path = r.URL.Path
	f.user, f.password, _ = r.BasicAuth()
	f.reqBody = nil
	if b, _ := io.ReadAll(r.Body); len(b) > 0 {
		_ = json.Unmarshal(b, &f.reqBody)
	}
	w.WriteHeader(f.status)
	if f.pages != nil {
		io.WriteString(w, f.pages[r.URL.Query().Get("page")])
		return
	}
	io.WriteString(w, f.body)
}

func newTestClient(t *testing.T, f *fakeAPI) *Client {
	t.Helper()
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, "jason", "secret")
}

func TestCreateRecord(t *testing.T) {
	f := &fakeAPI{status: 200, body: `{"id":42,"host":"hello-aws","fqdn":"hello-aws.example.com.","type":"CNAME","answer":"lb.aws.com","ttl":300}`}
	c := newTestClient(t, f)

	got, err := c.CreateRecord(context.Background(), "example.com", Record{ID: 7, Host: "hello-aws", Type: "CNAME", Answer: "lb.aws.com", TTL: 300})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{name: "method", got: f.method, want: http.MethodPost},
		{name: "path has no record id", got: f.path, want: "/core/v1/domains/example.com/records"},
		{name: "basic auth user", got: f.user, want: "jason"},
		{name: "basic auth token", got: f.password, want: "secret"},
		{name: "id not sent", got: f.reqBody["id"], want: nil},
		{name: "type sent", got: f.reqBody["type"], want: "CNAME"},
		{name: "decoded id", got: got.ID, want: 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestListRecordsPaginates(t *testing.T) {
	f := &fakeAPI{status: 200, pages: map[string]string{
		"1": `{"records":[{"id":1}],"nextPage":2,"lastPage":2}`,
		"2": `{"records":[{"id":2}]}`,
	}}
	c := newTestClient(t, f)

	recs, err := c.ListRecords(context.Background(), "example.com")
	if err != nil {
		t.Fatal(err)
	}
	var ids []int
	for _, r := range recs {
		ids = append(ids, r.ID)
	}
	if want := []int{1, 2}; !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
}

func TestCheckAvailability(t *testing.T) {
	f := &fakeAPI{status: 200, body: `{"results":[{"domainName":"foo.com","purchasable":true,"purchasePrice":10.99}]}`}
	c := newTestClient(t, f)

	res, err := c.CheckAvailability(context.Background(), []string{"foo.com"})
	if err != nil {
		t.Fatal(err)
	}
	if f.path != "/core/v1/domains:checkAvailability" {
		t.Errorf("path = %q, want unencoded colon", f.path)
	}
	if len(res) != 1 || !res[0].Purchasable || res[0].PurchasePrice != 10.99 {
		t.Errorf("results = %+v", res)
	}
}

func TestAPIErrors(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		wantMessage string
		wantError   string
	}{
		{
			name:        "json error body",
			status:      405,
			body:        `{"message":"Method Not Allowed"}`,
			wantMessage: "Method Not Allowed",
			wantError:   "GET /core/v1/domains/example.com: 405 Method Not Allowed",
		},
		{
			name:        "details included",
			status:      403,
			body:        `{"message":"Permission denied","details":"Failed authentication"}`,
			wantMessage: "Permission denied",
			wantError:   "GET /core/v1/domains/example.com: 403 Permission denied (Failed authentication)",
		},
		{
			name:        "non-json body falls back to status text",
			status:      502,
			body:        `<html>bad gateway</html>`,
			wantMessage: "",
			wantError:   "GET /core/v1/domains/example.com: 502 Bad Gateway",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, &fakeAPI{status: tt.status, body: tt.body})
			_, err := c.GetDomain(context.Background(), "example.com")

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("err = %v, want *APIError", err)
			}
			if apiErr.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", apiErr.Message, tt.wantMessage)
			}
			if err.Error() != tt.wantError {
				t.Errorf("Error() = %q, want %q", err.Error(), tt.wantError)
			}
		})
	}
}

func TestParseInterspersed(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantPos  []string
		wantType string
	}{
		{name: "domain first", args: []string{"example.com", "-type", "A"}, wantPos: []string{"example.com"}, wantType: "A"},
		{name: "flags first", args: []string{"-type", "A", "example.com"}, wantPos: []string{"example.com"}, wantType: "A"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			typ := fs.String("type", "", "")
			pos, err := parseInterspersed(fs, tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(pos, tt.wantPos) || *typ != tt.wantType {
				t.Errorf("pos = %v type = %q, want %v %q", pos, *typ, tt.wantPos, tt.wantType)
			}
		})
	}
}
