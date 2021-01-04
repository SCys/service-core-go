package core

import (
	"context"

	"github.com/phuslu/log"
)

var loggerSupportExtraFields = []string{"client_ip", "_id", "api"}

var loggerInfo = log.Logger{
	Level:  log.InfoLevel,
	Writer: &log.FileWriter{Filename: "log/info.log", MaxSize: 70 << 20, MaxBackups: 60},
}

var loggerWarn = log.Logger{
	Level:  log.WarnLevel,
	Writer: &log.FileWriter{Filename: "log/warn.log", MaxSize: 70 << 20, MaxBackups: 60},
}

var loggerError = log.Logger{
	Level: log.InfoLevel,
	Writer: &log.MultiWriter{
		InfoWriter:    &log.FileWriter{Filename: "log/error.log", MaxSize: 70 << 20, MaxBackups: 60},
		ConsoleWriter: &log.ConsoleWriter{ColorOutput: true},
		ConsoleLevel:  log.ErrorLevel,
	},
}

// loggerDebug debug logger
var loggerDebug = log.Logger{
	Level: log.DebugLevel,
	Writer: &log.MultiWriter{
		InfoWriter:    &log.FileWriter{Filename: "log/debug.log", MaxSize: 50 << 20, MaxBackups: 30},
		ConsoleWriter: &log.ConsoleWriter{ColorOutput: false},
		ConsoleLevel:  log.DebugLevel,
	},
}

func loggerExtraFieldContext(entry *log.Entry, msg string, params ...interface{}) *log.Entry {
	fields := make([]interface{}, 0, len(params))

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

		case H:
			entry = entry.Fields(v)

		default:
			fields = append(fields, v)
		}
	}

	entry.Msgf(msg, fields...)

	return entry
}

// I logging info message
func I(msg string, params ...interface{}) {
	loggerExtraFieldContext(loggerInfo.Info(), msg, params...)
}

// W logging warning message
func W(msg string, params ...interface{}) {
	loggerExtraFieldContext(loggerWarn.Warn(), msg, params...)
}

// E logging info message
func E(msg string, err error, params ...interface{}) {
	entry := loggerError.Error().Err(err)
	loggerExtraFieldContext(entry, msg, params...)
}

// F logging info message
func F(msg string, params ...interface{}) {
	loggerExtraFieldContext(loggerError.Fatal(), msg, params...)
}

// D logging info message
func D(msg string, params ...interface{}) {
	loggerExtraFieldContext(loggerDebug.Debug(), msg, params...)
}

// func init() {
// 	if log.IsTerminal(os.Stderr.Fd()) {
// 		Logger.Caller = 1
// 		Logger.Writer = &log.ConsoleWriter{
// 			ColorOutput:    true,
// 			QuoteString:    true,
// 			EndWithMessage: true,
// 		}
// 	}
// }
