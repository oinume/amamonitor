package http_server

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/oinume/amamonitor/backend/config"
)

type commonTemplateData struct {
	StaticURL string
	//GoogleAnalyticsID string
	CurrentURL   string
	CanonicalURL string
	//TrackingID        string
	IsUserAgentPC     bool
	IsUserAgentSP     bool
	IsUserAgentTablet bool
	//UserID            string
	//	NavigationItems   []navigationItem
}

func TemplatePath(file string) string {
	return path.Join(config.DefaultVars.TemplateDir(), file)
}

func ParseHTMLTemplates(files ...string) *template.Template {
	f := []string{
		TemplatePath("_base.html"),
	}
	f = append(f, files...)
	return template.Must(template.ParseFiles(f...))
}

func IsUserAgentPC(req *http.Request) bool {
	return !IsUserAgentSP(req) && !IsUserAgentTablet(req)
}

func IsUserAgentSP(req *http.Request) bool {
	ua := strings.ToLower(req.UserAgent())
	return strings.Contains(ua, "iphone") || strings.Contains(ua, "android") || strings.Contains(ua, "ipod")
}

func IsUserAgentTablet(req *http.Request) bool {
	ua := strings.ToLower(req.UserAgent())
	return strings.Contains(ua, "ipad")
}

func (s *server) getCommonTemplateData(req *http.Request, loggedIn bool, userID uint32) commonTemplateData {
	canonicalURL := fmt.Sprintf("%s://%s%s", config.DefaultVars.WebURLScheme(req), req.Host, req.RequestURI)
	canonicalURL = (strings.SplitN(canonicalURL, "?", 2))[0] // TODO: use url.Parse
	data := commonTemplateData{
		StaticURL: config.DefaultVars.StaticURL(),
		//		GoogleAnalyticsID: config.DefaultVars.GoogleAnalyticsID,
		CurrentURL:        req.RequestURI,
		CanonicalURL:      canonicalURL,
		IsUserAgentPC:     IsUserAgentPC(req),
		IsUserAgentSP:     IsUserAgentSP(req),
		IsUserAgentTablet: IsUserAgentTablet(req),
	}

	//if loggedIn {
	//	data.NavigationItems = loggedInNavigationItems
	//} else {
	//	data.NavigationItems = loggedOutNavigationItems
	//}
	//if flashMessageKey := req.FormValue("flashMessageKey"); flashMessageKey != "" {
	//	flashMessage, _ := s.flashMessageStore.Load(flashMessageKey)
	//	data.FlashMessage = flashMessage
	//}
	//data.TrackingID = context_data.MustTrackingID(req.Context())
	//if userID != 0 {
	//	data.UserID = fmt.Sprint(userID)
	//}

	return data
}

//func getRemoteAddress(req *http.Request) string {
//	xForwardedFor := req.Header.Get("X-Forwarded-For")
//	if xForwardedFor == "" {
//		return (strings.Split(req.RemoteAddr, ":"))[0]
//	}
//	return strings.TrimSpace((strings.Split(xForwardedFor, ","))[0])
//}
