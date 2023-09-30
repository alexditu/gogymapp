package server

import (
	"io"
	"net/http"

	"log"
)

func newHttpServer(s *server, w *io.PipeWriter) *http.Server {
	return &http.Server{
		Addr:    s.setts.URL(),
		Handler: s.router,
		// create a stdlib log.Logger that writes to logrus.Logger.
		ErrorLog: log.New(w, "", 0),
	}
}
