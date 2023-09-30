package server

import (
	"strconv"

	"github.com/alexditu/gogymapp/internal/logging"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`

	// TLS
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`

	// path to static html pages
	HtmlPath string `json:"htmlPath"`

	Log logging.Settings `json:"log"`

	// Sign-In With Google
	ClientId string `json:"clientId"`
}

func (s *Settings) InitDefault() {
	s.Host = "localhost"
	s.Port = 8443
	s.CertPath = "./cert.pem"
	s.KeyPath = "./key.pem"
	s.KeyPath = "./html"
	s.Log = logging.Settings{LogToFile: true, FileName: "server.log", Level: log.DebugLevel}
}

func (s *Settings) URL() string {
	return s.Host + ":" + strconv.FormatUint(uint64(s.Port), 10)
}
