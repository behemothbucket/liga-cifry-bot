package ui

import (
	"net/http"
	"os"
	"path/filepath"
	"telegram-bot/internal/logger"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Index(webappURL string) func(writer http.ResponseWriter, request *http.Request) {
	cwd, _ := os.Getwd()
	return func(writer http.ResponseWriter, request *http.Request) {
		template := template.Must(
			template.ParseFiles(filepath.Join(cwd, "./internal/web_app/ui/index.html")),
		)
		err := template.ExecuteTemplate(
			writer,
			"index.html",
			struct {
				WebAppURL string
			}{
				WebAppURL: webappURL,
			},
		)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = writer.Write([]byte(err.Error()))
			if err != nil {
				logger.Fatal("failed to write templates", "err", err)
			}
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
