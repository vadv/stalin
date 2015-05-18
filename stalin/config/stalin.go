package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	. "stalin/log"
	"strings"
)

type StalinConfig struct {
	Logger   *logrus.Logger
	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`
}

var Stalin = initStalinConfig()

func initStalinConfig() *StalinConfig {
	logger := logrus.New()
	formater := &logrus.TextFormatter{ForceColors: true, DisableColors: false}
	logger.Formatter = formater
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout
	Log.SetLogger(logger)
	return &StalinConfig{Logger: logger, LogFile: "STDOUT", LogLevel: "DEBUG"}
}

func logFileFromString(log string) *os.File {
	switch strings.ToUpper(log) {
	case "STDOUT", "":
		Log.Info("[Stalin]: Set LogFile to 'STDOUT'")
		return os.Stdout
	default:
		Log.Info("[Stalin]: Set LogFile to '%s'", log)
		f, err := os.Create(log)
		if err != nil {
			Log.Error("[Stalin]: Set LogFile to '%s' error: %v", log, err)
			os.Exit(2)
		}
		return f
	}
}

func logLevelFromString(log_level string) logrus.Level {
	switch strings.ToUpper(log_level) {
	case "INFO":
		Log.Info("[Stalin]: Set LogLevel to 'INFO'")
		return logrus.InfoLevel
	case "ERROR", "ERR":
		Log.Info("[Stalin]: Set LogLevel to 'ERROR'")
		return logrus.ErrorLevel
	case "WARN", "WARNING":
		Log.Info("[Stalin]: Set LogLevel to 'WARN'")
		return logrus.WarnLevel
	default:
		Log.Info("[Stalin]: Set LogLevel to 'DEBUG'")
		return logrus.DebugLevel
	}
}

func (s *StalinConfig) Configure(config interface{}) error {
	daemon_config := &StalinConfig{}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, daemon_config)
	if err != nil {
		return err
	}
	s.Logger.Out = logFileFromString(daemon_config.LogFile)
	s.Logger.Level = logLevelFromString(daemon_config.LogLevel)
	Log.SetLogger(s.Logger)
	return nil
}
