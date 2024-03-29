# pkgslog

[![Build](https://github.com/reugn/pkgslog/actions/workflows/build.yml/badge.svg)](https://github.com/reugn/pkgslog/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/reugn/pkgslog)](https://pkg.go.dev/github.com/reugn/pkgslog)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/pkgslog)](https://goreportcard.com/report/github.com/reugn/pkgslog)
[![codecov](https://codecov.io/gh/reugn/pkgslog/branch/main/graph/badge.svg)](https://codecov.io/gh/reugn/pkgslog)

A structured log handler for `log/slog` that can be configured at the package level.  
`pkgslog` allows you to set the minimum log level requirement for each package in your application.

## Usage Example

```go
textHandler := slog.NewTextHandler(os.Stdout, nil)
packageMap := map[string]slog.Level{
    "github.com/reugn/pkgslog/internal":  slog.LevelWarn,
    "github.com/reugn/pkgslog/pkg":       slog.LevelDebug,
    "github.com/reugn/pkgslog/pkg/inner": slog.LevelInfo,
}
logger := slog.New(pkgslog.NewPackageHandler(textHandler, packageMap))
```

## Benchmarks

Benchmark results compared to the standard `slog.TextHandler`

```
BenchmarkPkgSlog-16       590906              2029 ns/op             232 B/op          2 allocs/op
BenchmarkSlog-16         1733892               701.6 ns/op             0 B/op          0 allocs/op
```

## License

Licensed under the MIT License.
