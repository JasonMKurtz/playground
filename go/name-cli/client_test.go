package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestUpdateRecord(t *testing.T) {
	f := &fakeAPI{status: 200, body: `{"id":42,"host":"hello-aws","fqdn":"hello-aws.example.com.","type":"CNAME","answer":"new-lb.aws.com","ttl":300}`}
	c := newTestClient(t, f)

	got, err := c.UpdateRecord(context.Background(), "example.com", 42, Record{ID: 7, FQDN: "x.", Host: "hello-aws", Type: "CNAME", Answer: "new-lb.aws.com", TTL: 300})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{name: "method", got: f.method, want: http.MethodPut},
		{name: "path has record id", got: f.path, want: "/core/v1/domains/example.com/records/42"},
		{name: "id not sent in body", got: f.reqBody["id"], want: nil},
		{name: "fqdn not sent", got: f.reqBody["fqdn"], want: nil},
		{name: "answer sent", got: f.reqBody["answer"], want: "new-lb.aws.com"},
		{name: "decoded answer", got: got.Answer, want: "new-lb.aws.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

// records update reads the current record and only changes the flags given.
// The fake answers both the GET and the PUT with the current record, and
// remembers the last request (the PUT).
func TestRecordsUpdateKeepsUnsetFields(t *testing.T) {
	current := `{"id":42,"host":"hello-aws","fqdn":"hello-aws.example.com.","type":"CNAME","answer":"old-lb.aws.com","ttl":600}`

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantMethod string
		wantBody   map[string]any
	}{
		{
			name:       "answer only keeps host, type, ttl",
			args:       []string{"records", "update", "example.com", "42", "-answer", "new-lb.aws.com"},
			wantMethod: http.MethodPut,
			wantBody:   map[string]any{"host": "hello-aws", "type": "CNAME", "answer": "new-lb.aws.com", "ttl": float64(600)},
		},
		{
			name:       "flags before positionals, type upper-cased",
			args:       []string{"records", "update", "-ttl", "300", "-type", "cname", "example.com", "42"},
			wantMethod: http.MethodPut,
			wantBody:   map[string]any{"host": "hello-aws", "type": "CNAME", "answer": "old-lb.aws.com", "ttl": float64(300)},
		},
		{
			name:       "explicit empty host moves it to the apex",
			args:       []string{"records", "update", "example.com", "42", "-host", ""},
			wantMethod: http.MethodPut,
			wantBody:   map[string]any{"host": "", "type": "CNAME", "answer": "old-lb.aws.com", "ttl": float64(600)},
		},
		{
			name:       "no flags is an error and sends nothing",
			args:       []string{"records", "update", "example.com", "42"},
			wantErr:    true,
			wantMethod: "",
		},
		{
			name:       "non-numeric id",
			args:       []string{"records", "update", "example.com", "abc", "-ttl", "300"},
			wantErr:    true,
			wantMethod: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fakeAPI{status: 200, body: current}
			srv := httptest.NewServer(f)
			t.Cleanup(srv.Close)
			t.Setenv("NAMECOM_USERNAME", "jason")
			t.Setenv("NAMECOM_TOKEN", "secret")
			t.Setenv("NAMECOM_API_URL", srv.URL)

			err := run(context.Background(), tt.args, io.Discard)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if f.method != tt.wantMethod {
				t.Errorf("last request method = %q, want %q", f.method, tt.wantMethod)
			}
			if tt.wantBody == nil {
				return
			}
			if f.path != "/core/v1/domains/example.com/records/42" {
				t.Errorf("path = %q", f.path)
			}
			for k, want := range tt.wantBody {
				if got := f.reqBody[k]; !reflect.DeepEqual(got, want) {
					t.Errorf("body[%q] = %v, want %v", k, got, want)
				}
			}
		})
	}
}

func TestDeleteRecord(t *testing.T) {
	f := &fakeAPI{status: 200, body: `{}`}
	c := newTestClient(t, f)

	if err := c.DeleteRecord(context.Background(), "example.com", 42); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{name: "method", got: f.method, want: http.MethodDelete},
		{name: "path has record id", got: f.path, want: "/core/v1/domains/example.com/records/42"},
		{name: "no body sent", got: f.reqBody, want: map[string]any(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}

// deleteAPI is a fake zone for records delete. GET /records/<id> returns the
// record until it's deleted, then 404. DELETE of failID returns 500, and the
// first ignore[id] DELETEs of an ID return 200 without deleting (as name.com
// has been seen to do). deletes lists the IDs actually removed, in order;
// calls counts DELETE requests per ID.
type deleteAPI struct {
	failID  string
	ignore  map[string]int
	gone    map[string]bool
	deletes []string
	calls   map[string]int
}

func (d *deleteAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	switch r.Method {
	case http.MethodGet:
		if d.gone[id] {
			w.WriteHeader(404)
			io.WriteString(w, `{"message":"Not Found"}`)
			return
		}
		fmt.Fprintf(w, `{"id":%s,"host":"hello-%s","fqdn":"hello-%s.example.com.","type":"A","answer":"192.0.2.1","ttl":300}`, id, id, id)
	case http.MethodDelete:
		d.calls[id]++
		if id == d.failID {
			w.WriteHeader(500)
			io.WriteString(w, `{"message":"Internal Server Error"}`)
			return
		}
		if d.ignore[id] > 0 {
			d.ignore[id]--
			io.WriteString(w, `{}`) // claims success, keeps the record
			return
		}
		d.gone[id] = true
		d.deletes = append(d.deletes, id)
		io.WriteString(w, `{}`)
	}
}

func TestRecordsDelete(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		failID      string
		ignore      map[string]int
		wantErr     bool
		wantDeletes []string
		wantCalls   map[string]int
		wantOutput  []string
		notInOutput []string
	}{
		{
			name:        "several ids deleted in order and listed",
			args:        []string{"records", "delete", "example.com", "11", "22", "33"},
			wantDeletes: []string{"11", "22", "33"},
			wantCalls:   map[string]int{"11": 1, "22": 1, "33": 1},
			wantOutput:  []string{"Deleted:", "hello-11.example.com.", "hello-33.example.com."},
		},
		{
			name:        "api error stops but lists what was already deleted",
			args:        []string{"records", "delete", "example.com", "11", "22", "33"},
			failID:      "22",
			wantErr:     true,
			wantDeletes: []string{"11"},
			wantCalls:   map[string]int{"11": 1, "22": 1},
			wantOutput:  []string{"Deleted:", "hello-11.example.com."},
			notInOutput: []string{"hello-22.example.com."},
		},
		{
			name:        "delete that silently didn't happen is retried",
			args:        []string{"records", "delete", "example.com", "11", "22"},
			ignore:      map[string]int{"22": 1},
			wantDeletes: []string{"11", "22"},
			wantCalls:   map[string]int{"11": 1, "22": 2},
			wantOutput:  []string{"Deleted:", "hello-11.example.com.", "hello-22.example.com."},
		},
		{
			name:        "record that never goes away is an error, not reported deleted",
			args:        []string{"records", "delete", "example.com", "11", "22"},
			ignore:      map[string]int{"22": 5},
			wantErr:     true,
			wantDeletes: []string{"11"},
			wantCalls:   map[string]int{"11": 1, "22": 2},
			wantOutput:  []string{"Deleted:", "hello-11.example.com."},
			notInOutput: []string{"hello-22.example.com."},
		},
		{
			name:        "non-numeric id deletes nothing",
			args:        []string{"records", "delete", "example.com", "11", "abc"},
			wantErr:     true,
			wantDeletes: nil,
			wantCalls:   map[string]int{},
		},
		{
			name:        "no ids is a usage error",
			args:        []string{"records", "delete", "example.com"},
			wantErr:     true,
			wantDeletes: nil,
			wantCalls:   map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ignore := map[string]int{}
			for k, v := range tt.ignore {
				ignore[k] = v
			}
			api := &deleteAPI{failID: tt.failID, ignore: ignore, gone: map[string]bool{}, calls: map[string]int{}}
			srv := httptest.NewServer(api)
			t.Cleanup(srv.Close)
			t.Setenv("NAMECOM_USERNAME", "jason")
			t.Setenv("NAMECOM_TOKEN", "secret")
			t.Setenv("NAMECOM_API_URL", srv.URL)

			var out strings.Builder
			err := run(context.Background(), tt.args, &out)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(api.deletes, tt.wantDeletes) {
				t.Errorf("deletes = %v, want %v", api.deletes, tt.wantDeletes)
			}
			if !reflect.DeepEqual(api.calls, tt.wantCalls) {
				t.Errorf("DELETE calls = %v, want %v", api.calls, tt.wantCalls)
			}
			for _, s := range tt.wantOutput {
				if !strings.Contains(out.String(), s) {
					t.Errorf("output missing %q:\n%s", s, out.String())
				}
			}
			for _, s := range tt.notInOutput {
				if strings.Contains(out.String(), s) {
					t.Errorf("output lists %q, whose delete didn't happen:\n%s", s, out.String())
				}
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
