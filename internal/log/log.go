package log

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Level zerolog.Level

const (
	ErrorLevel Level = Level(zerolog.ErrorLevel)
	WarnLevel  Level = Level(zerolog.WarnLevel)
	InfoLevel  Level = Level(zerolog.InfoLevel)
	DebugLevel Level = Level(zerolog.DebugLevel)
)

var isLogInit = false
var logger zerolog.Logger

func getLogger() *zerolog.Logger {
	if isLogInit {
		return &logger
	} else {
		return &log.Logger
	}
}

func Init(level Level, logFolder string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.Level(level))

	err := os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		Error(fmt.Errorf("Error in creating log file. %v", err))
	}

	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s.log", logFolder, time.Now().Format("02-01-2006")), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Error(fmt.Errorf("Error in creating log file. %v", err))
	}

	w := zerolog.MultiLevelWriter(logFile, zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123})
	logger = zerolog.New(w).With().Stack().Timestamp().Caller().Logger()

	isLogInit = true
}

func Panic(err error, msg string) {
	getLogger().Panic().Err(err).Msg(msg)
}

func Error(err error) {
	getLogger().Error().Err(err).Send()
}

func Info(msg string, args ...any) {
	getLogger().Info().Msg(fmt.Sprintf(msg, args...))
}

func Debug(msg string, args ...any) {
	getLogger().Debug().Msg(fmt.Sprintf(msg, args...))
}

func LoggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// log end request
				getLogger().Info().
					Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					}).
					Msg("request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
