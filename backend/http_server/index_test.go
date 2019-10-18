package http_server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/amamonitor/backend/service"
)

func Test_server_index(t *testing.T) {
	db := getRealDB()
	server := New(db, service.New(db) /* TODO: mock */)
	tests := map[string]struct {
		method         string
		path           string
		query          map[string]string
		handler        http.HandlerFunc
		wantStatusCode int
	}{
		"index_ok": {
			method:         "GET",
			path:           "/index",
			handler:        server.index,
			wantStatusCode: http.StatusOK,
			// TODO: bodyValidator
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("GET", test.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			//q := req.URL.Query()
			//for k, v := range test.query {
			//	q.Set(k, v)
			//}
			//req.URL.RawQuery = q.Encode()
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()

			test.handler(rr, req)
			result := rr.Result()
			if result.StatusCode != test.wantStatusCode {
				t.Errorf("unexpected status code: got=%v, want=%v", result.StatusCode, test.wantStatusCode)
			}
			body, err := ioutil.ReadAll(result.Body)
			if err != nil {
				t.Fatal(err)
			}
			if len(body) == 0 {
				t.Error("empty response body")
			}
		})
	}
}
