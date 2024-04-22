package api

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Index(webappURL string) func(writer http.ResponseWriter, request *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	indexTmpl := template.Must(template.ParseFiles(wd + "internal/server/resources/index.html"))

	return func(writer http.ResponseWriter, request *http.Request) {
		err := indexTmpl.ExecuteTemplate(
			writer,
			wd+"internal/server/resources/index.html",
			struct {
				WebAppURL string
			}{
				WebAppURL: webappURL,
			},
		)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

func Validate(token string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ok, err := ext.ValidateWebAppQuery(request.URL.Query(), token)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("validation failed; error: " + err.Error()))
			return
		}
		if ok {
			writer.Write([]byte("validation success; user is authenticated."))
		} else {
			writer.Write([]byte("validation failed; data cannot be trusted."))
		}
	}
}
