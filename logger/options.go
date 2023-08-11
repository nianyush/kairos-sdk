package logger

import (
	"github.com/rs/zerolog"
)

type LogOption func(a *KairosLog)

func WithLevel(level Level) func(k *KairosLog) {
	return func(k *KairosLog) {
		k.SetLevel(level)
	}
}

// WithDebugFunction allows to pass a function that will set the logger to debug if the function returns true
// This is done so logger consumers can implement easily the check for debug, be it a file, a flag in cmd or an
// environment variable check
func WithDebugFunction(f func() bool) func(k *KairosLog) {
	return func(k *KairosLog) {
		if f() {
			k.SetLevel(DebugLevel)
		}
	}
}

// WithStringContext overrides the default context of the existing logger to provide extra info in the log
func WithStringContext(key, value string) func(k *KairosLog) {
	return func(k *KairosLog) {
		k.Logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			c.Str(key, value)
			return c
		})
	}
}
