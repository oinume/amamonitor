package http_server

import (
	"net/http"
)

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	t := ParseHTMLTemplates(TemplatePath("index.html"))
	type Data struct {
		commonTemplateData
	}
	data := &Data{
		commonTemplateData: s.getCommonTemplateData(r, false, 0),
	}
	writeHTMLWithTemplate(w, http.StatusOK, t, data)
}
