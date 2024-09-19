package core

import (
	"context"
	"strings"

	"github.com/phuslu/log"
)

var loggerSupportExtraFields = []string{"client_ip", "_id"}

var LoggerInfo = log.Logger{
	Level: log.InfoLevel,
	Writer: &log.FileWriter{
		Filename:   "log/info.log",
		MaxSize:    25 << 20,
		MaxBackups: 10,
		ProcessID:  false,
		HostName:   false,
	},
}

var LoggerWarn = log.Logger{
	Level: log.WarnLevel,
	Writer: &log.FileWriter{
		Filename:   "log/warn.log",
		MaxSize:    25 << 20,
		MaxBackups: 10,
		ProcessID:  false,
		HostName:   false,
	},
}

var LoggerError = log.Logger{
	Level: log.ErrorLevel,
	Writer: &log.FileWriter{
		Filename:   "log/error.log",
		MaxSize:    25 << 20,
		MaxBackups: 10,
		ProcessID:  false,
		HostName:   false,
	},
}

// loggerDebug debug logger
var loggerDebug = log.Logger{
	Level: log.DebugLevel,
	Writer: &log.MultiWriter{
		InfoWriter: &log.FileWriter{
			Filename:   "log/debug.log",
			MaxSize:    50 << 20,
			MaxBackups: 30,
			ProcessID:  false,
			HostName:   false,
		},
		ConsoleWriter: &log.ConsoleWriter{ColorOutput: false},
		ConsoleLevel:  log.DebugLevel,
	},
}

// LogWithExtraFields adds extra fields to a log entry based on the provided parameters.
func LogWithExtraFields(entry *log.Entry, msg string, params ...any) *log.Entry {
	fields := make([]any, 0, len(params))

	for _, it := range params {
		switch v := it.(type) {
		case context.Context:
			for _, key := range loggerSupportExtraFields {
				if value, ok := v.Value(key).(string); ok {
					entry = entry.Str(key, value)
				}
			}

			if user, ok := v.Value("user").(BasicFields); ok {
				entry = entry.Str("user", user.ID)
			}

			if err, ok := v.Value("error").(error); ok {
				entry = entry.Err(err)
			}

			if s, ok := v.Value("api").(string); ok {
				var sb strings.Builder

				sb.WriteString("[")
				sb.WriteString(s)
				sb.WriteString("]")
				sb.WriteString(msg)

				msg = sb.String()
			}

		case map[string]any:
			entry = entry.Fields(v)

		case H:
			entry = entry.Fields(v.Fields())

		default:
			fields = append(fields, v)
		}
	}

	entry.Msgf(msg, fields...)

	return entry
}

// I logging info message
func I(msg string, params ...any) {
	LogWithExtraFields(LoggerInfo.Info(), msg, params...)
}

// W logging warning message
func W(msg string, params ...any) {
	LogWithExtraFields(LoggerWarn.Warn(), msg, params...)
}

// E logging info message
func E(msg string, err error, params ...any) {
	entry := LoggerError.Error().Err(err)
	LogWithExtraFields(entry, msg, params...)
}

// F logging info message
func F(msg string, params ...any) {
	LogWithExtraFields(LoggerError.Fatal(), msg, params...)
}

// D logging info message
func D(msg string, params ...any) {
	LogWithExtraFields(loggerDebug.Debug(), msg, params...)
}
