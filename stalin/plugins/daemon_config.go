package plugins

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var RuntimeStalinConfig = initDaemonConfig()

type DaemonConfig struct {
	Logger   *logrus.Logger
	logLevel string `json:"log_level"`
}

func initDaemonConfig() *DaemonConfig {
	// init logger
	logger := logrus.New()
	formater := &logrus.TextFormatter{ForceColors: true, DisableColors: false}
	logger.Formatter = formater
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout
	return &DaemonConfig{Logger: logger}
}

func (d *DaemonConfig) SetConfig(config interface{}) error {
	daemon_config := &DaemonConfig{}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, daemon_config)
	if err != nil {
		return err
	}
	d.Logger.Level = logLevelFromString(daemon_config.logLevel)
	return nil
}

func logLevelFromString(log_level string) logrus.Level {
	switch strings.ToUpper(log_level) {
	case "INFO":
		LogInfo("[Stalin]: Set LogLevel 'INFO'")
		return logrus.InfoLevel
	case "ERROR", "ERR":
		LogInfo("[Stalin]: Set LogLevel 'ERROR'")
		return logrus.ErrorLevel
	case "WARN", "WARNING":
		LogInfo("[Stalin]: Set LogLevel 'WARN'")
		return logrus.WarnLevel
	default:
		LogInfo("[Stalin]: Set LogLevel 'DEBUG'")
		return logrus.DebugLevel
	}
}
