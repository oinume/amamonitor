package config

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Vars struct {
	MySQLUser            string `envconfig:"MYSQL_USER"`
	MySQLPassword        string `envconfig:"MYSQL_PASSWORD"`
	MySQLHost            string `envconfig:"MYSQL_HOST"`
	MySQLPort            string `envconfig:"MYSQL_PORT"`
	MySQLDatabase        string `envconfig:"MYSQL_DATABASE"`
	NodeEnv              string `envconfig:"NODE_ENV"`
	ServiceEnv           string `envconfig:"AMAMONITOR_ENV" required:"true"`
	GCPProjectID         string `envconfig:"GCP_PROJECT_ID"`
	GCPServiceAccountKey string `envconfig:"GCP_SERVICE_ACCOUNT_KEY"`
	//EnableFetcherHTTP2        bool   `envconfig:"ENABLE_FETCHER_HTTP2" default:"true"`
	//EnableTrace               bool   `envconfig:"ENABLE_TRACE"`
	//EnableStackdriverProfiler bool   `envconfig:"ENABLE_STACKDRIVER_PROFILER"`
	//GoogleClientID            string `envconfig:"GOOGLE_CLIENT_ID"`
	//GoogleClientSecret        string `envconfig:"GOOGLE_CLIENT_SECRET"`
	//GoogleAnalyticsID         string `envconfig:"GOOGLE_ANALYTICS_ID"`
	HTTPPort int `envconfig:"PORT" default:"5001"`
	//GRPCPort                  int    `envconfig:"GRPC_PORT" default:"4002"`
	//RollbarAccessToken        string `envconfig:"ROLLBAR_ACCESS_TOKEN"`
	VersionHash string `envconfig:"VERSION_HASH"`
	//DebugSQL                  bool   `envconfig:"DEBUG_SQL"`
	//ZipkinReporterURL         string `envconfig:"ZIPKIN_REPORTER_URL"`
	//LocalLocation             *time.Location
}

func Process() (*Vars, error) {
	var vars Vars
	if err := envconfig.Process("", &vars); err != nil {
		return nil, err
	}

	//vars.LocalLocation = asiaTokyo
	//if vars.VersionHash == "" {
	//	vars.VersionHash = timestamp.Format("20060102150405")
	//}
	return &vars, nil
}

func MustProcess() *Vars {
	vars, err := Process()
	if err != nil {
		panic(err)
	}
	return vars
}

var DefaultVars = &Vars{}
var once sync.Once

func MustProcessDefault() {
	once.Do(func() {
		DefaultVars = MustProcess()
	})
}

func (v *Vars) IsProductionEnv() bool {
	return v.ServiceEnv == "production"
}

func (v *Vars) IsDevelopmentEnv() bool {
	return v.ServiceEnv == "development"
}

func (v *Vars) DBURL() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		v.MySQLUser, v.MySQLPassword, v.MySQLHost, v.MySQLPort, v.MySQLDatabase,
	)
}

func (v *Vars) XODBURL() string {
	// mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@$(MYSQL_HOST):$(MYSQL_PORT)/$(MYSQL_DATABASE)?charset=utf8mb4&parseTime=true&loc=UTC
	return fmt.Sprintf(
		"mysql://%s:%s@%s:%s/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		v.MySQLUser, v.MySQLPassword, v.MySQLHost, v.MySQLPort, v.MySQLDatabase,
	)
}

func (v *Vars) StaticURL() string {
	if v.IsProductionEnv() {
		return "https://asset.amamonitor.com/static/" + v.VersionHash
	} else if v.IsDevelopmentEnv() {
		return "http://asset.local.amamonitor.com/static/" + v.VersionHash
	} else {
		return "/static/" + v.VersionHash
	}
}

func (v *Vars) WebURL() string {
	if v.IsProductionEnv() {
		return "https://www.lekcije.com"
	} else if v.IsDevelopmentEnv() {
		return "http://www.local.lekcije.com"
	} else {
		return "http://localhost:5000"
	}
}

func (v *Vars) WebURLScheme(r *http.Request) string {
	if v.IsProductionEnv() {
		return "https"
	}
	if r != nil && r.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}

//func WebURLScheme(r *http.Request) string {
//	return DefaultVars.WebURLScheme(r)
//}
