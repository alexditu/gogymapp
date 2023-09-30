package server

import (
	"strconv"

	"github.com/alexditu/gogymapp/internal/logging"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Host     string
	Port     uint16
	CertPath string
	KeyPath  string
	Log      logging.Settings
}

func (s *Settings) InitDefault() {
	s.Host = "localhost"
	s.Port = 8443
	s.CertPath = "./cert.pem"
	s.KeyPath = "./key.pem"
	s.Log = logging.Settings{LogToFile: true, FileName: "server.log", Level: log.DebugLevel}
}

func (s *Settings) URL() string {
	return s.Host + ":" + strconv.FormatUint(uint64(s.Port), 10)
}
