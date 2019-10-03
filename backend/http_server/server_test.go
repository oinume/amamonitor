package http_server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/oinume/amamonitor/backend/config"
	"github.com/oinume/amamonitor/backend/service"
	"github.com/xo/dburl"

	"github.com/oinume/amamonitor/backend/fetcher"
)

func TestMain(m *testing.M) {
	config.MustProcessDefault()
	os.Exit(m.Run())
}

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

	db, err := dburl.Open(config.DefaultVars.XODBURL())
	if err != nil {
		t.Fatal(err)
	}
	server := New(db, service.New(db) /* TODO: mock */)
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
