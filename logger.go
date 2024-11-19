package consumer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	Logger   zerolog.Logger
	multi    zerolog.LevelWriter
	logPath  = "./logs"
	fileName = fmt.Sprintf(
		"%s/%s.log", logPath, strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-"),
	)
)

func init() {
	// Get log level from environment variable
	logLevelStr := strings.ToLower(os.Getenv("CONSUMER_LOG_LEVEL"))

	logLevel := zerolog.WarnLevel

	if logLevelStr != "" {
		switch logLevelStr {
		case zerolog.LevelDebugValue:
			logLevel = zerolog.DebugLevel
		case zerolog.LevelInfoValue:
			logLevel = zerolog.InfoLevel
		case zerolog.LevelWarnValue:
			logLevel = zerolog.WarnLevel
		case zerolog.LevelErrorValue:
			logLevel = zerolog.ErrorLevel
		case zerolog.LevelFatalValue:
			logLevel = zerolog.FatalLevel
		case zerolog.LevelPanicValue:
			logLevel = zerolog.PanicLevel
		}
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	// Log time always in UTC
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
	// Default level is info
	zerolog.SetGlobalLevel(logLevel)
	if os.Getenv("LOG_TO_FILE") != "" {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.TimeOnly,
		}
		_ = os.MkdirAll(logPath, 0755)
		runLogFile, _ := os.OpenFile(
			fileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		multi = zerolog.MultiLevelWriter(consoleWriter, runLogFile)
		Logger = zerolog.New(multi).With().Caller().Timestamp().Logger()
	} else {
		Logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	}
}

func GetLogger(component string) zerolog.Logger {
	return Logger.With().Str("component", component).Logger()
}
