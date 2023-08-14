package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// KairosLog implements the bridge between zerolog and the logger.Interface that yip needs.
// We also use that interface across all kairos libs, so its easier to bridge it here
type KairosLog struct {
	zerolog.Logger
}

func (k KairosLog) Infof(tpl string, args ...interface{}) {
	k.Logger.Info().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Info(args ...interface{}) {
	k.Logger.Info().Msg(fmt.Sprint(args...))
}
func (k KairosLog) Warnf(tpl string, args ...interface{}) {
	k.Logger.Warn().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Warn(args ...interface{}) {
	k.Logger.Warn().Msg(fmt.Sprint(args...))
}
func (k KairosLog) Debugf(tpl string, args ...interface{}) {
	k.Logger.Debug().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Debug(args ...interface{}) {
	k.Logger.Debug().Msg(fmt.Sprint(args...))
}
func (k KairosLog) Errorf(tpl string, args ...interface{}) {
	k.Logger.Error().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Error(args ...interface{}) {
	k.Logger.Error().Msg(fmt.Sprint(args...))
}
func (k KairosLog) Fatalf(tpl string, args ...interface{}) {
	k.Logger.Fatal().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Fatal(args ...interface{}) {
	k.Logger.Fatal().Msg(fmt.Sprint(args...))
}
func (k KairosLog) Panicf(tpl string, args ...interface{}) {
	k.Logger.Panic().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Panic(args ...interface{}) {
	k.Logger.Panic().Msg(fmt.Sprint(args...))
}
func (k KairosLog) Tracef(tpl string, args ...interface{}) {
	k.Logger.Trace().Msg(fmt.Sprintf(tpl, args...))
}
func (k KairosLog) Trace(args ...interface{}) {
	k.Logger.Trace().Msg(fmt.Sprint(args...))
}
func (k KairosLog) SetLevel(level Level) {
	k.Logger.Level(zerolog.Level(level))
}
func (k KairosLog) IsDebugLevel() bool {
	return k.Logger.GetLevel() == zerolog.DebugLevel
}

// Fix to set a decent time format in the console output otherwise its a very simple output
func TimeFormatConsole(w *zerolog.ConsoleWriter) {
	w.TimeFormat = time.RFC3339
}

// NewKairosLog provides a normal console log
func NewKairosLog(opts ...LogOption) (*KairosLog, error) {
	k := &KairosLog{zerolog.New(zerolog.NewConsoleWriter(TimeFormatConsole)).With().Timestamp().Logger()}
	for _, o := range opts {
		o(k)
	}
	return k, nil
}

// NewKairosLogToFile provides a log that logs to a given file
// file is closed automatically by GC on application exit
func NewKairosLogToFile(logfile string, opts ...LogOption) (*KairosLog, error) {
	f, err := os.Open(logfile)
	if err != nil {
		return nil, err
	}
	k := &KairosLog{zerolog.New(f).With().Timestamp().Logger()}
	for _, o := range opts {
		o(k)
	}
	return k, nil
}

// NewKairosMultiLog provides a logger that logs to both console and a given file
func NewKairosMultiLog(logfile string, opts ...LogOption) (*KairosLog, error) {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	multiWriter := zerolog.MultiLevelWriter(zerolog.NewConsoleWriter(TimeFormatConsole), zerolog.NewConsoleWriter(TimeFormatConsole, func(w *zerolog.ConsoleWriter) {
		w.Out = f
	}))
	k := &KairosLog{zerolog.New(multiWriter).With().Timestamp().Logger()}
	for _, o := range opts {
		o(k)
	}
	return k, nil
}

// NewKairosNullLog provides a logger that discards all output
func NewKairosNullLog(opts ...LogOption) (*KairosLog, error) {
	k := &KairosLog{zerolog.New(nil).With().Timestamp().Logger().Level(zerolog.Level(Disabled))}
	for _, o := range opts {
		o(k)
	}
	return k, nil
}

// NewKairosBufferLog will return a logger that stores all logs in a buffer, used mainly for testing
func NewKairosBufferLog(b *bytes.Buffer, opts ...LogOption) (*KairosLog, error) {
	k := &KairosLog{zerolog.New(zerolog.NewConsoleWriter(TimeFormatConsole, func(w *zerolog.ConsoleWriter) {
		w.Out = b
	})).With().Timestamp().Logger()}
	for _, o := range opts {
		o(k)
	}
	return k, nil
}

// NewKairosMultiCustomLogTargets provides a multi log with custom writers so we can override the writers easily
// The burden of generating the writers is on the caller
// This should be called with writers that implement the zerolog.LevelWriter interface so zerolog can properly write to
// them based on the level and such. Otherwise it will use the Writer method directly and miss a lot fo stuff
func NewKairosMultiCustomLogTargets(writers ...io.Writer) (*KairosLog, error) {
	multiWriter := zerolog.MultiLevelWriter(writers...)
	k := &KairosLog{zerolog.New(multiWriter).With().Timestamp().Logger()}
	return k, nil
}
