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

	if err := t.Execute(w, data); err != nil {
		internalServerError(w, err)
		//internalServerError(s.appLogger, w, errors.NewInternalError(
		//	errors.WithError(err),
		//	errors.WithMessage("Failed to template.Execute()"),
		//), 0)
		return
	}
}
