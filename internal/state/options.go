/*
Package state defines the global state and options used by all commands.
*/
package state

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const LogTimeFormat = "15:04:05"

const (
	ColorDebug     = "\033[90m" // Bright black/gray
	ColorInfo      = "\033[36m" // Cyan
	ColorLink      = "\033[34m" // Blue
	ColorBackup    = "\033[35m" // Magenta
	ColorOverwrite = "\033[33m" // Yellow
	ColorWarning   = "\033[33m" // Yellow
	ColorError     = "\033[31m" // Red
	ColorDryRun    = "\033[93m" // Bright yellow
	ColorReset     = "\033[0m"  // Reset
)

// colorize wraps text in ANSI color codes
func colorize(text, color string) string {
	return color + text + ColorReset
}

type TrovlOptions struct {
	Verbose      bool
	Debug        bool
	DryRun       bool
	UseRelative  bool
	OverwriteYes bool
	OverwriteNo  bool
	BackupDir    string
	BackupYes    bool
	BackupNo     bool
}

type TrovlState struct {
	Options *TrovlOptions
	Logger  *slog.Logger
	Level   *slog.LevelVar
}

func New(opts *TrovlOptions) *TrovlState {
	lvl := &slog.LevelVar{}

	if opts == nil {
		opts = &TrovlOptions{}
	}

	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      lvl,
		AddSource:  opts.Debug,
		TimeFormat: LogTimeFormat,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Value.Kind() == slog.KindAny {
				if _, ok := a.Value.Any().(error); ok {
					return tint.Attr(9, a)
				}
			}

			if a.Key == slog.MessageKey && opts.DryRun {
				dryRunTag := colorize("[DRY-RUN]", ColorDryRun)
				return slog.String(a.Key, dryRunTag+" "+a.Value.String())
			}

			return a
		},
	}))

	if opts.DryRun {
		logger = logger.With("dry_run", true)
	}

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

func (s *TrovlState) LogLink(msg string, args ...any) {
	taggedMsg := colorize("[LINK]", ColorLink) + " " + msg
	s.Logger.Info(taggedMsg, args...)
}

func (s *TrovlState) LogBackup(msg string, args ...any) {
	taggedMsg := colorize("[BACKUP]", ColorBackup) + " " + msg
	s.Logger.Info(taggedMsg, args...)
}

func (s *TrovlState) LogOverwrite(msg string, args ...any) {
	taggedMsg := colorize("[OVERWRITE]", ColorOverwrite) + " " + msg
	s.Logger.Info(taggedMsg, args...)
}

func (s *TrovlState) LogDebug(msg string, args ...any) {
	s.Logger.Debug(msg, args...)
}

func (s *TrovlState) LogInfo(msg string, args ...any) {
	s.Logger.Info(msg, args...)
}

func (s *TrovlState) LogWarn(msg string, args ...any) {
	s.Logger.Warn(msg, args...)
}

func (s *TrovlState) LogError(msg string, args ...any) {
	s.Logger.Error(msg, args...)
}
