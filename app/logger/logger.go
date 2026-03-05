// Package logger provides structured console logging with timestamp, level (INF/WRN/ERR),
// message, and optional key=value pairs so all errors and events are visible in a consistent format.
// When Color is true and the output is a TTY, level and status codes are colorized.
package logger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

// ANSI color codes for terminal output (only used when Color is true).
const (
	ansiReset  = "\033[0m"
	ansiDim    = "\033[90m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
	ansiCyan   = "\033[36m"
)

// Level is the log level.
type Level string

const (
	LevelInfo Level = "INF"
	LevelWarn Level = "WRN"
	LevelErr  Level = "ERR"
)

// Logger writes structured log lines: "3:04AM LVL message key1=value1 key2=value2".
// When Color is true, level and status codes are emitted with ANSI colors for TTYs.
type Logger struct {
	w     io.Writer
	Color bool
}

// Default returns a logger that writes to os.Stdout with color enabled when stdout is a TTY.
func Default() *Logger {
	return &Logger{w: os.Stdout, Color: isTerminal(os.Stdout)}
}

// New returns a logger that writes to w (no color).
func New(w io.Writer) *Logger {
	return &Logger{w: w, Color: false}
}

// NewWithColor returns a logger that writes to w; when color is true, level and status are colorized.
func NewWithColor(w io.Writer, color bool) *Logger {
	return &Logger{w: w, Color: color}
}

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

// Log writes a line with the given level, message, and key-value pairs (alternating key, value).
// Keys are quoted if they contain spaces; values are space-escaped for readability.
func (l *Logger) Log(level Level, msg string, kv ...interface{}) {
	line := l.format(level, msg, kv)
	_, _ = fmt.Fprintln(l.w, line)
}

func (l *Logger) format(level Level, msg string, kv []interface{}) string {
	ts := time.Now().Format("3:04AM")
	parts := []string{ts, string(level), msg}
	for i := 0; i+1 < len(kv); i += 2 {
		k := fmt.Sprint(kv[i])
		v := fmt.Sprint(kv[i+1])
		if strings.ContainsAny(k, " \t=") {
			k = fmt.Sprintf("%q", k)
		}
		if strings.ContainsAny(v, " \t\n") {
			v = strings.ReplaceAll(v, "\n", " ")
			v = strings.TrimSpace(v)
			if len(v) > 128 {
				v = v[:128] + "..."
			}
			v = fmt.Sprintf("%q", v)
		}
		pair := k + "=" + v
		if l.Color {
			switch k {
			case "status":
				pair = l.colorStatus(pair, v)
			case "method":
				pair = "method=" + ansiCyan + v + ansiReset
			case "path":
				pair = "path=" + ansiCyan + v + ansiReset
			case "request_id", "duration_ms":
				pair = k + "=" + ansiDim + v + ansiReset
			}
		}
		parts = append(parts, pair)
	}
	s := strings.Join(parts, " ")
	if l.Color {
		s = l.colorizeLine(s, level, ts)
	}
	return s
}

// colorizeLine adds ANSI colors for timestamp (dim), level (green/warn/err).
func (l *Logger) colorizeLine(line string, level Level, ts string) string {
	levelColor := ansiGreen
	switch level {
	case LevelWarn:
		levelColor = ansiYellow
	case LevelErr:
		levelColor = ansiRed
	}
	levelStr := string(level)
	// Replace first occurrence of ts and level with colored versions
	afterTs := strings.TrimPrefix(line, ts)
	if len(afterTs) == len(line) {
		return line
	}
	afterLevel := strings.TrimPrefix(strings.TrimLeft(afterTs, " "), levelStr)
	if len(afterLevel) == len(strings.TrimLeft(afterTs, " ")) {
		return line
	}
	rest := strings.TrimLeft(afterLevel, " ")
	return ansiDim + ts + ansiReset + " " + levelColor + levelStr + ansiReset + " " + rest
}

// colorStatus returns key=value with ANSI color on value for HTTP status (2xx green, 4xx yellow, 5xx red).
func (l *Logger) colorStatus(pair, v string) string {
	n, err := strconv.Atoi(v)
	if err != nil {
		return pair
	}
	c := ansiGreen
	if n >= 500 {
		c = ansiRed
	} else if n >= 400 {
		c = ansiYellow
	}
	return "status=" + c + v + ansiReset
}

// Info logs at INF level.
func (l *Logger) Info(msg string, kv ...interface{}) {
	l.Log(LevelInfo, msg, kv...)
}

// Infof logs at INF level with a format string.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warn logs at WRN level.
func (l *Logger) Warn(msg string, kv ...interface{}) {
	l.Log(LevelWarn, msg, kv...)
}

// Warnf logs at WRN level with a format string.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Err logs at ERR level.
func (l *Logger) Err(msg string, kv ...interface{}) {
	l.Log(LevelErr, msg, kv...)
}

// Errf logs at ERR level with a format string.
func (l *Logger) Errf(format string, args ...interface{}) {
	l.Err(fmt.Sprintf(format, args...))
}

// defaultLogger is the package-level logger used by apilog and optional callers.
var defaultLogger = Default()

// SetDefault sets the logger used by package-level apilog. Main can call this to use a test buffer.
func SetDefault(l *Logger) {
	defaultLogger = l
}

// DefaultLogger returns the current default logger.
func DefaultLogger() *Logger {
	return defaultLogger
}
