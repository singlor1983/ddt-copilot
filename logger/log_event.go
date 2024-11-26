package logger

import (
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type LogEvent struct {
	*zerolog.Event
}

type DebugStr interface {
	DebugStr() string
}

func (e *LogEvent) DebugStr(key string, msg DebugStr) *LogEvent {
	if e.Event != nil && msg != nil {
		e.Str(key, msg.DebugStr())
	}
	return e
}

func (e *LogEvent) HttpRequest(r *http.Request) *LogEvent {
	if e.Event != nil {
		e.Str("method", r.Method)
		e.Str("url", r.URL.String())
		e.Str("remoteaddr", r.RemoteAddr)
		e.Str("cookie", r.Header.Get("Cookie"))
		e.Str("content-type", r.Header.Get("Content-Type"))
		e.Str("host", r.Header.Get("Host"))
		e.Str("x-real-ip", r.Header.Get("X-Real-IP"))
		e.Str("x-forwarded-for", r.Header.Get("X-Forwarded-For"))
	}
	return e
}

func (e *LogEvent) Timestamp() *LogEvent {
	e.Int64("timestamp_ms", time.Now().UnixMilli())
	return e
}

func (e *LogEvent) UID(uid int64) *LogEvent {
	e.Int64("uid", uid)
	return e
}
