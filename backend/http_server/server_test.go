package http_server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_internalServerError(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_server_fetcher(t *testing.T) {
	server := New()
	tests := map[string]struct {
		path    string
		handler http.HandlerFunc
	}{
		"fetcher_ok": {
			path:    "/fetcher",
			handler: server.fetcher,
			// TODO: status code
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("GET", test.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			test.handler(rr, req)
			result := rr.Result()
			if result.StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: got=%v, want=%v", result.StatusCode, http.StatusOK)
			}
		})
	}
}

func Test_writeJSON(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		code int
		body interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
