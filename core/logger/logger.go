package logger

import (
	"fmt"
	"time"
)

// Logger é a interface que o app usa. Pluga o que quiser aqui.
type Logger interface {
	WithAppTag(tag string) Logger
	Logf(ctx, format string, args ...any)
	Warnf(ctx, format string, args ...any)
	Errorf(ctx, format string, args ...any)
	Debugf(ctx, format string, args ...any)
}

// ===================== Implementação default (console) =====================

type ConsoleLogger struct {
	appTag string
	now    func() time.Time
}

func NewConsole() *ConsoleLogger {
	return &ConsoleLogger{
		now: func() time.Time { return time.Now() },
	}
}

func (l *ConsoleLogger) WithAppTag(tag string) Logger {
	cp := *l
	cp.appTag = tag
	return &cp
}

func (l *ConsoleLogger) Logf(ctx, format string, args ...any)   { l.out("LOG", ctx, format, args...) }
func (l *ConsoleLogger) Warnf(ctx, format string, args ...any)  { l.out("WARN", ctx, format, args...) }
func (l *ConsoleLogger) Errorf(ctx, format string, args ...any) { l.out("ERROR", ctx, format, args...) }
func (l *ConsoleLogger) Debugf(ctx, format string, args ...any) { l.out("DEBUG", ctx, format, args...) }

func (l *ConsoleLogger) out(level, ctx, format string, args ...any) {
	ts := l.now().Format("01/02/2006, 03:04:05 PM")
	// Formato estilo Nest:
	// [app] <ts>   LOG [Context] message
	fmt.Printf("[%s] %s\t%s [%s] %s\n", l.appTag, ts, level, ctx, fmt.Sprintf(format, args...))
}
