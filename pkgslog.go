package pkgslog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// A PackageHandler represents the package level structured log handler.
type PackageHandler struct {
	handler  slog.Handler
	packages map[string]slog.Level
}

var _ slog.Handler = &PackageHandler{}

// NewPackageHandler returns a new PackageHandler given an underlying slog.Handler
// and a map of package names to their minimum slog.Level, e.g.
//
//	textHandler := slog.NewTextHandler(os.Stdout, nil)
//	packageMap := map[string]slog.Level{
//		"github.com/reugn/pkgslog/internal":  slog.LevelWarn,
//		"github.com/reugn/pkgslog/pkg":       slog.LevelDebug,
//		"github.com/reugn/pkgslog/pkg/inner": slog.LevelInfo,
//	}
//	logger := slog.New(pkgslog.NewPackageHandler(textHandler, packageMap))
func NewPackageHandler(h slog.Handler, packages map[string]slog.Level) *PackageHandler {
	if ph, ok := h.(*PackageHandler); ok {
		h = ph.handler
	}
	return &PackageHandler{h, packages}
}

// Enabled reports whether the handler handles records at the given level.
func (h *PackageHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Scan our package log levels and see if we have any
	// packages configured at the requested log level.
	//
	// If we do, we will be handling that log level
	for _, pkgLevel := range h.packages {
		if pkgLevel <= level {
			return true
		}
	}

	// Otherwise defer to upstream handler wether we should handle the
	// log level or not
	return h.handler.Enabled(ctx, level)
}

// Handle handles the Record.
// If the log writing package has a higher minimum log level configured,
// the currenl log record will be discarded.
// It will only be called when Enabled returns true.
func (h *PackageHandler) Handle(ctx context.Context, r slog.Record) error {
	pkg := callerPackageName()

	if configuredLevel, ok := h.packages[pkg]; ok {
		if configuredLevel > r.Level {
			return nil // discard the log record
		}
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *PackageHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewPackageHandler(h.handler.WithAttrs(attrs), h.packages)
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *PackageHandler) WithGroup(name string) slog.Handler {
	return NewPackageHandler(h.handler.WithGroup(name), h.packages)
}

// callerPackageName returns the log writer package name.
func callerPackageName() string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(5, pc)
	if n == 0 {
		// no PCs available
		return "NA"
	}

	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()

	// extract caller package name
	var pkg string
	lastSlashIndex := strings.LastIndexByte(frame.Function, '/')
	if lastSlashIndex < 0 {
		pkg = frame.Function[:strings.IndexByte(frame.Function, '.')]
	} else {
		packageUpperIndex := indexByteFrom(frame.Function, lastSlashIndex, '.')
		pkg = frame.Function[:packageUpperIndex]
	}
	return pkg
}

// printFrames is a runtime.CallersFrames debug method.
func printFrames(pc []uintptr) {
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		fmt.Println(frame.Function)
		if !more {
			break
		}
	}
}

// indexByteFrom returns the index of the first instance of c in s starting
// from start, or -1 if c is not present in s.
func indexByteFrom(s string, start int, c byte) int {
	for i := start; i <= len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
