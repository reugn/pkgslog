package pkgslog_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"slices"
	"testing"
	"testing/slogtest"

	"github.com/reugn/pkgslog"
)

// go test -bench=. -benchmem
func BenchmarkPkgSlog(b *testing.B) {
	textHandler := slog.NewTextHandler(io.Discard, nil)
	packageMap := map[string]slog.Level{
		"github.com/reugn/pkgslog/internal":  slog.LevelWarn,
		"github.com/reugn/pkgslog/pkg":       slog.LevelDebug,
		"github.com/reugn/pkgslog/pkg/inner": slog.LevelInfo,
	}
	logger := slog.New(pkgslog.NewPackageHandler(textHandler, packageMap))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("here")
	}
}

func BenchmarkSlog(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("here")
	}
}

func TestPkgSlog(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	textHandler := slog.NewJSONHandler(writer, nil)
	packageMap := map[string]slog.Level{
		"github.com/reugn/pkgslog/internal": slog.LevelWarn,
	}
	logger := pkgslog.NewPackageHandler(textHandler, packageMap)

	if err := slogtest.TestHandler(logger, func() []map[string]any {
		writer.Flush()
		return parseLogEntries(t, buf.Bytes())
	}); err != nil {
		t.Fatal(err)
	}
}

func TestPkgSlogOverrideUpstreamHandlerLogLevel(t *testing.T) { //nolint:funlen
	levels := []slog.Level{
		slog.LevelDebug,
		slog.LevelError,
		slog.LevelInfo,
		slog.LevelWarn,
	}

	tests := []struct {
		name         string
		pkgLevel     slog.Level
		handlerLevel slog.Level
		enabled      []slog.Level
		disabled     []slog.Level
	}{
		// If
		// 		"pkgslog" level is set to DEBUG
		// 		"handler" level is set to WARNING
		// Then
		// 		"pkgslog" should accept log levels appropriate
		// 		to its *minimum* configured log level
		// Meaning
		// 		All log levels should be accepted
		{
			name:         "debug level",
			pkgLevel:     slog.LevelDebug,
			handlerLevel: slog.LevelWarn,
			enabled: []slog.Level{
				slog.LevelDebug,
				slog.LevelError,
				slog.LevelInfo,
				slog.LevelWarn,
			},
			disabled: []slog.Level{},
		},
		// If
		// 		"pkgslog" level is set to INFO
		// 		"handler" level is set to WARNING
		// Then
		// 		"pkgslog" should accept log levels appropriate
		// 		to its *minimum* configured log level
		// Meaning
		// 		DEBUG level should be rejected
		//		All other levels should be accepted
		{
			name:         "info level",
			pkgLevel:     slog.LevelInfo,
			handlerLevel: slog.LevelWarn,
			enabled: []slog.Level{
				slog.LevelError,
				slog.LevelInfo,
				slog.LevelWarn,
			},
			disabled: []slog.Level{
				slog.LevelDebug,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			textHandler := slog.NewJSONHandler(
				writer,
				&slog.HandlerOptions{
					Level: tt.handlerLevel,
				},
			)

			packageMap := map[string]slog.Level{
				"github.com/reugn/pkgslog": tt.pkgLevel,
			}

			logger := pkgslog.NewPackageHandler(textHandler, packageMap)
			for _, checkLevel := range levels {
				if logger.Enabled(context.TODO(), checkLevel) {
					if !slices.Contains(tt.enabled, checkLevel) {
						t.Fatal("expected " + checkLevel.String() + " level to be DISABLED")
					}
					continue
				}

				if !slices.Contains(tt.disabled, checkLevel) {
					t.Fatal("expected " + checkLevel.String() + " level to be DISABLED")
				}
			}
		})
	}
}

func parseLogEntries(t *testing.T, output []byte) []map[string]any {
	var entries []map[string]any
	lines := bytes.Split(output, []byte("\n"))

	for i := 0; i < len(lines)-1; i++ { // last one is empty
		line := lines[i]
		var entry map[string]any
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Error(err)
		}
		entries = append(entries, entry)
	}

	return entries
}
