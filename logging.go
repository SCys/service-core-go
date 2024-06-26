package core

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/phuslu/log"
)

var loggerSupportExtraFields = []string{"client_ip", "_id"}
var loggerRoot = log.Logger{Level: log.DebugLevel}

type SimpleFileWriter struct {
	Filename   string
	Filemode   fs.FileMode
	MaxSize    int64
	MaxBackups int

	Written    int64 // RO 记录写入数量
	File       *os.File
	TimeFormat string
	Mut        sync.Mutex
}

func (w *SimpleFileWriter) WriteEntry(e *log.Entry) (n int, err error) {
	return w.Write(e.Value())
}

func (w *SimpleFileWriter) Write(p []byte) (n int, err error) {
	if w.File == nil {
		return len(p), nil
	}

	n, err = w.File.Write(p)
	if err != nil {
		return
	}

	w.Written += int64(n)
	if w.MaxSize > 0 && w.Written > w.MaxSize && w.Filename != "" {
		err = w.Rotate()
	}

	return
}

func (w *SimpleFileWriter) Open() (err error) {
	err = os.MkdirAll(filepath.Dir(w.Filename), 0755)
	if err == nil {
		w.File, err = os.OpenFile(w.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			w.Mut = sync.Mutex{}
		}
	}

	return
}

func (w *SimpleFileWriter) Rotate() (err error) {
	// rename to new filename
	now := time.Now()
	prefix := filepath.Base(w.Filename)
	name := prefix + "." + now.Format(time.RFC3339) + ".log"
	os.Rename(w.Filename, name)

	// keep old
	old := w.File

	f, err := os.OpenFile(w.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	w.File = f

	go func(old *os.File, prefix string) {
		// change uid/gid
		uid, _ := strconv.Atoi(os.Getenv("SUDO_UID"))
		gid, _ := strconv.Atoi(os.Getenv("SUDO_GID"))
		if uid != 0 && gid != 0 && os.Geteuid() == 0 {
			// _ = os.Lchown(prefix+".log", uid, gid)
			_ = os.Chown(prefix+".log", uid, gid)
		}

		dir := filepath.Dir(w.Filename)
		dirfile, err := os.Open(dir)
		if err != nil {
			return
		}
		infos, err := dirfile.Readdir(-1)
		dirfile.Close()
		if err != nil {
			return
		}

		matches := make([]os.FileInfo, 0)
		for _, info := range infos {
			name := info.Name()
			if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".gz") {
				matches = append(matches, info)
			}
		}
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].ModTime().Unix() < matches[j].ModTime().Unix()
		})

		for i := 0; i < len(matches)-w.MaxBackups-1; i++ {
			os.Remove(filepath.Join(dir, matches[i].Name()))
		}
	}(old, prefix)

	return
}

func NewFileWriter(name string, maxSize int64, maxBackups int) *SimpleFileWriter {
	w := &SimpleFileWriter{Filename: name, MaxSize: maxSize, MaxBackups: maxBackups}

	if err := w.Open(); err != nil {
		loggerRoot.Fatal().Msgf("open log failed:%s", err.Error())
	}

	return w
}

var LoggerInfo = log.Logger{
	Level:  log.InfoLevel,
	Writer: NewFileWriter("log/info.log", 70<<20, 60),
}

var LoggerWarn = log.Logger{
	Level:  log.WarnLevel,
	Writer: NewFileWriter("log/warn.log", 70<<20, 60),
}

var LoggerError = log.Logger{
	Level:  log.ErrorLevel,
	Writer: NewFileWriter("log/error.log", 70<<20, 60),
}

// loggerDebug debug logger
var loggerDebug = log.Logger{
	Level: log.DebugLevel,
	Writer: &log.MultiWriter{
		InfoWriter: &log.FileWriter{
			Filename: "log/debug.log", MaxSize: 50 << 20, MaxBackups: 30,
			ProcessID: false,
			HostName:  false,
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
