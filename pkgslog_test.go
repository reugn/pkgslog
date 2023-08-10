package pkgslog_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"log/slog"

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
	for i := 0; i < b.N; i++ {
		logger.Info("here")
	}
}

func BenchmarkSlog(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
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
