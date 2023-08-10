# pkgslog

A package level structured log handler for `log/slog`.  
`pkgslog` adds the ability to set minimum log level requirement per package.

## Example

```go
textHandler := slog.NewTextHandler(os.Stdout, nil)
packageMap := map[string]slog.Level{
    "github.com/reugn/pkgslog/internal":  slog.LevelWarn,
    "github.com/reugn/pkgslog/pkg":       slog.LevelDebug,
    "github.com/reugn/pkgslog/pkg/inner": slog.LevelInfo,
}
logger := slog.New(pkgslog.NewPackageHandler(textHandler, packageMap))
```

## Benchmarking

Benchmark results compared to the standard `slog.TextHandler`

```
BenchmarkPkgSlog
BenchmarkPkgSlog-16       437230              2421 ns/op             232 B/op          2 allocs/op
BenchmarkSlog
BenchmarkSlog-16         1000000              1004 ns/op               0 B/op          0 allocs/op
```

## License

Licensed under the MIT License.
