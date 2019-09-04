package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oinume/amamonitor/backend/fetcher"
)

type server struct{}

func NewServer() *server {
	return &server{}
}

func (s *server) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/fetcher/{type}", s.fetcher).Methods("GET")
	return r
}

func (s *server) fetcher(w http.ResponseWriter, r *http.Request) {
	amaten, err := fetcher.NewAmatenClient()
	if err != nil {
		internalServerError(w, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		internalServerError(w, err)
	}

	options := new(fetcher.FetchOptions)
	if amatenURL := r.FormValue("amatenUrl"); amatenURL != "" {
		options.URL = amatenURL
	}
	_, err = amaten.Fetch(r.Context(), options)
	if err != nil {
		internalServerError(w, err)
		return
	}

	// TODO: Write gifts to DB

	w.WriteHeader(http.StatusOK)
}

func internalServerError(w http.ResponseWriter, err error) {
	//switch _ := errors.Cause(err).(type) { // TODO:
	//default:
	// unknown error
	//sUserID := ""
	//if userID == 0 {
	//	sUserID = fmt.Sprint(userID)
	//}
	//util.SendErrorToRollbar(err, sUserID)
	//fields := []zapcore.Field{
	//	zap.Error(err),
	//}
	//if e, ok := err.(errors.StackTracer); ok {
	//	b := &bytes.Buffer{}
	//	for _, f := range e.StackTrace() {
	//		fmt.Fprintf(b, "%+v\n", f)
	//	}
	//	fields = append(fields, zap.String("stacktrace", b.String()))
	//}
	//if appLogger != nil {
	//	appLogger.Error("internalServerError", fields...)
	//}

	http.Error(w, fmt.Sprintf("Internal Server Error\n\n%v", err), http.StatusInternalServerError)
	//if !config.IsProductionEnv() {
	//	fmt.Fprintf(w, "----- stacktrace -----\n")
	//	if e, ok := err.(errors.StackTracer); ok {
	//		for _, f := range e.StackTrace() {
	//			fmt.Fprintf(w, "%+v\n", f)
	//		}
	//	}
	//}
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, `{ "status": "Failed to Encode as writeJSON" }`, http.StatusInternalServerError)
		return
	}
}
