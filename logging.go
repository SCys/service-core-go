package core

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	LSupportExtraFields = []string{"client_ip", "_id"} // extra fields that can be added to the log entry
	LMain               = logrus.New()                 // main logger
	LogPath             = "./log/main.log"             // 日志文件路径，可配置
)

func InitLog() {
	// LMain.SetOutput(os.Stdout)
	LMain.SetLevel(logrus.DebugLevel)
	LMain.SetFormatter(&logrus.TextFormatter{})

	// 确保日志目录存在
	logDir := "./log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		LMain.WithError(err).Error("Failed to create log directory")
	}

	// os multi writer
	mw := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   LogPath,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	})
	LMain.SetOutput(mw)
}

// LogWithExtraFields adds extra fields to a log entry based on the provided parameters.
func LogWithExtraFields(entry *logrus.Entry, msg string, params ...any) {
	for idx, it := range params {
		switch v := it.(type) {
		case context.Context:
			for _, key := range LSupportExtraFields {
				if value, ok := v.Value(key).(string); ok {
					entry = entry.WithField(key, value)
				}
			}

			if user, ok := v.Value("user").(BasicFields); ok {
				entry = entry.WithField("user", user.ID)
			}

			if s, ok := v.Value("api").(string); ok {
				var sb strings.Builder

				sb.WriteString("[")
				sb.WriteString(s)
				sb.WriteString("]")
				sb.WriteString(msg)

				msg = sb.String()
			}

			entry = entry.WithContext(v)
			if idx+1 < len(params) {
				params = params[idx+1:]
			} else {
				params = params[:0]
			}
		case H:
			entry = entry.WithFields(logrus.Fields(v))
			if idx+1 < len(params) {
				params = params[idx+1:]
			} else {
				params = params[:0]
			}
		}
	}

	// entry.Logf(entry, msg, params...)
	entry.Printf(msg, params...)
}

// I logging info message
func I(msg string, params ...any) {
	entry := LMain.WithFields(logrus.Fields{})
	entry.Level = logrus.InfoLevel
	LogWithExtraFields(entry, msg, params...)
}

// W logging warning message
func W(msg string, params ...any) {
	entry := LMain.WithFields(logrus.Fields{})
	entry.Level = logrus.WarnLevel
	LogWithExtraFields(entry, msg, params...)
}

// E logging info message
func E(msg string, err error, params ...any) {
	entry := LMain.WithError(err)
	entry.Level = logrus.ErrorLevel
	LogWithExtraFields(entry, msg, append(params, err)...)
}

// F logging info message
func F(msg string, params ...any) {
	entry := LMain.WithFields(logrus.Fields{})
	entry.Level = logrus.FatalLevel
	LogWithExtraFields(entry, msg, params...)
}

// D logging info message
func D(msg string, params ...any) {
	entry := LMain.WithFields(logrus.Fields{})
	entry.Level = logrus.DebugLevel
	LogWithExtraFields(entry, msg, params...)
}
