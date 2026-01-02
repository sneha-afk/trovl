package state

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const logTimeFormat = "2006-01-02 15:04:05"

type TrovlOptions struct {
	Verbose      bool
	Debug        bool
	UseRelative  bool
	OverwriteYes bool
	OverwriteNo  bool
	DryRun       bool
}

type TrovlState struct {
	Options *TrovlOptions
	Logger  *slog.Logger
	Level   *slog.LevelVar
}

func New(opts *TrovlOptions) *TrovlState {
	lvl := &slog.LevelVar{}

	// handlerOpts := &slog.HandlerOptions{
	// 	Level:     lvl,
	// 	AddSource: opts.Debug, // Show file/line only during debug
	// 	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr { // [LEVEL] file:line msg key=val
	// 		if a.Key == slog.TimeKey {
	// 			return slog.Attr{}
	// 		}
	// 		if a.Key == slog.LevelKey {
	// 			return slog.Attr{
	// 				Key:   "",
	// 				Value: slog.StringValue(fmt.Sprintf("[%s]", a.Value.String())),
	// 			}
	// 		}
	//
	// 		if a.Key == slog.SourceKey {
	// 			source := a.Value.Any().(*slog.Source)
	// 			return slog.Attr{
	// 				Key:   "",
	// 				Value: slog.StringValue(fmt.Sprintf("[%s:%d]", source.File, source.Line)),
	// 			}
	// 		}
	//
	// 		return a
	// 	},
	// }

	if opts == nil {
		opts = &TrovlOptions{}
	}

	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:     lvl,
		AddSource: opts.Debug,
		// TimeFormat: time.RFC3339,
		TimeFormat: logTimeFormat,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindAny {
				if _, ok := a.Value.Any().(error); ok {
					return tint.Attr(9, a)
				}
			}
			return a
		},
	}))

	state := TrovlState{
		Options: opts,
		Logger:  logger,
		Level:   lvl,
	}
	state.SetLogLevel()

	return &state
}

func DefaultState() *TrovlState {
	return New(&TrovlOptions{})
}

func (s *TrovlState) Verbose() bool {
	return s.Options.Verbose
}

func (s *TrovlState) Debug() bool {
	return s.Options.Debug
}

// SetLogLevel should be called after setting or changing the log level
func (s *TrovlState) SetLogLevel() {
	switch {
	case s.Options.Debug:
		s.Level.Set(slog.LevelDebug)
	case s.Options.Verbose, s.Options.DryRun:
		s.Level.Set(slog.LevelInfo)
	default:
		s.Level.Set(slog.LevelWarn)
	}
}
