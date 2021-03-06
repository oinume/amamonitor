package http_server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/oinume/amamonitor/backend/fetcher"
	"github.com/oinume/amamonitor/backend/model"
	"github.com/oinume/amamonitor/backend/service"
)

type server struct {
	db      *sql.DB
	service *service.Service
}

func New(db *sql.DB, svc *service.Service) *server {
	return &server{
		db:      db,
		service: svc,
	}
}

func (s *server) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/fetcher/{provider}", s.fetcher).Methods("GET")
	r.HandleFunc("/", s.index).Methods("GET")
	return r
}

func (s *server) fetcher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider, ok := vars["provider"]
	if !ok {
		http.Error(w, "no provider", http.StatusBadRequest)
		return
	}

	f, err := fetcher.NewFromProvider(fetcher.Provider(provider))
	if err != nil {
		internalServerError(w, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		internalServerError(w, err)
		return
	}

	options := new(fetcher.FetchOptions)
	if u := r.FormValue("url"); u != "" {
		options.URL = u
	}
	fetchedGiftItems, err := f.Fetch(r.Context(), options)
	if err != nil {
		internalServerError(w, err)
		return
	}

	var (
		fetchResult *model.FetchResult
		giftItems   []*model.GiftItem
	)
	if err := model.Transaction(r.Context(), s.db, nil, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		fetchResult, giftItems, err = s.service.CreateFetchResultGiftItems(
			r.Context(), tx, fetchedGiftItems, time.Now(),
		)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		internalServerError(w, err)
		return
	}

	type body struct {
		FetchResult *model.FetchResult `json:"fetchResult"`
		GiftItems   []*model.GiftItem  `json:"giftItems"`
	}
	writeJSON(w, http.StatusOK, &body{
		FetchResult: fetchResult,
		GiftItems:   giftItems,
	})
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
	}
}

//func writeHTML(w http.ResponseWriter, code int, body string) {
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	w.WriteHeader(code)
//	if _, err := fmt.Fprint(w, body); err != nil {
//		http.Error(w, "Failed to write HTML", http.StatusInternalServerError)
//	}
//}

func writeHTMLWithTemplate(
	w http.ResponseWriter,
	code int,
	t *template.Template,
	data interface{},
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	if err := t.Execute(w, data); err != nil {
		internalServerError(w, err)
	}
}
