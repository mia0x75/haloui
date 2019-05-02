package g

import (
	log "github.com/sirupsen/logrus"
)

// InitLog TODO
func InitLog() {
	cfg := Config()
	if level, err := log.ParseLevel(cfg.Log.Level); err == nil {
		log.SetLevel(level)
	}
}
