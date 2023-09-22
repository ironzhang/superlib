package echoutil

import (
	"io"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"github.com/ironzhang/tlog"
)

type Logger struct {
	tlog tlog.Logger
}

var _ echo.Logger = &Logger{}

func NewLogger() *Logger {
	return &Logger{
		tlog: tlog.Named("echo"),
	}
}

func (p *Logger) Output() io.Writer {
	return os.Stdout
}

func (p *Logger) SetOutput(w io.Writer) {
}

func (p *Logger) Prefix() string {
	return ""
}

func (p *Logger) SetPrefix(prefix string) {
}

func (p *Logger) Level() log.Lvl {
	return log.INFO
}

func (p *Logger) SetLevel(v log.Lvl) {
}

func (p *Logger) SetHeader(h string) {
}

func (p *Logger) Print(i ...interface{}) {
	p.tlog.Print(1, tlog.INFO, i...)
}

func (p *Logger) Printf(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.INFO, format, args...)
}

func (p *Logger) Printj(j log.JSON) {
	p.tlog.Printw(1, tlog.INFO, "printj", "data", j)
}

func (p *Logger) Debug(i ...interface{}) {
	p.tlog.Print(1, tlog.DEBUG, i...)
}

func (p *Logger) Debugf(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.DEBUG, format, args...)
}

func (p *Logger) Debugj(j log.JSON) {
	p.tlog.Printw(1, tlog.DEBUG, "debugj", "data", j)
}

func (p *Logger) Info(i ...interface{}) {
	p.tlog.Print(1, tlog.INFO, i...)
}

func (p *Logger) Infof(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.INFO, format, args...)
}

func (p *Logger) Infoj(j log.JSON) {
	p.tlog.Printw(1, tlog.INFO, "infoj", "data", j)
}

func (p *Logger) Warn(i ...interface{}) {
	p.tlog.Print(1, tlog.WARN, i...)
}

func (p *Logger) Warnf(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.WARN, format, args...)
}

func (p *Logger) Warnj(j log.JSON) {
	p.tlog.Printw(1, tlog.WARN, "warnj", "data", j)
}

func (p *Logger) Error(i ...interface{}) {
	p.tlog.Print(1, tlog.ERROR, i...)
}

func (p *Logger) Errorf(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.ERROR, format, args...)
}

func (p *Logger) Errorj(j log.JSON) {
	p.tlog.Printw(1, tlog.ERROR, "errorj", "data", j)
}

func (p *Logger) Panic(i ...interface{}) {
	p.tlog.Print(1, tlog.PANIC, i...)
}

func (p *Logger) Panicf(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.PANIC, format, args...)
}

func (p *Logger) Panicj(j log.JSON) {
	p.tlog.Printw(1, tlog.PANIC, "panicj", "data", j)
}

func (p *Logger) Fatal(i ...interface{}) {
	p.tlog.Print(1, tlog.FATAL, i...)
}

func (p *Logger) Fatalj(j log.JSON) {
	p.tlog.Printw(1, tlog.FATAL, "fatalj", "data", j)
}

func (p *Logger) Fatalf(format string, args ...interface{}) {
	p.tlog.Printf(1, tlog.FATAL, format, args...)
}
