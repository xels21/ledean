package ledean

import log "github.com/sirupsen/logrus"

func SetLogger(logLevelStr string) error {
	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)
	return nil
}
