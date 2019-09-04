package http_server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/amamonitor/backend/fetcher"
)

//func Test_internalServerError(t *testing.T) {
//	type args struct {
//		w   http.ResponseWriter
//		err error
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//		})
//	}
//}

func Test_server_fetcher(t *testing.T) {
	gifts := []fetcher.AmatenGift{
		{
			ID:        123,
			FaceValue: 10000,
			Price:     8710,
			Rate:      "87.1",
		},
		{
			ID:        456,
			FaceValue: 1000,
			Price:     900,
			Rate:      "90.0",
		},
	}
	fakeHandler := fetcher.NewFakeAmatenAPIGiftsHandler(t, gifts)
	ts := httptest.NewServer(http.HandlerFunc(fakeHandler))
	defer ts.Close()

	server := New()
	tests := map[string]struct {
		method         string
		path           string
		query          map[string]string
		handler        http.HandlerFunc
		wantStatusCode int
	}{
		"fetcher_ok": {
			method: "GET",
			path:   "/fetcher/all",
			query: map[string]string{
				"amatenUrl": ts.URL,
			},
			handler:        server.fetcher,
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
			q := req.URL.Query()
			for k, v := range test.query {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
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
