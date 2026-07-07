package quantum

// Logger is the minimal logging surface used by the client. It is deliberately
// narrow so any logging library can satisfy it with a trivial adapter, and so
// the client never depends on a concrete logger (dependency inversion).
type Logger interface {
	// Logf logs a formatted message. Semantics mirror fmt.Printf.
	Logf(format string, args ...any)
}

// nopLogger is the default logger; it discards everything.
type nopLogger struct{}

func (nopLogger) Logf(string, ...any) {}

// LoggerFunc adapts a plain function to the Logger interface.
//
//	client, _ := quantum.NewClient(
//		quantum.WithAPIKey(key),
//		quantum.WithCompanyID(id),
//		quantum.WithLogger(quantum.LoggerFunc(log.Printf)),
//	)
type LoggerFunc func(format string, args ...any)

// Logf implements Logger.
func (f LoggerFunc) Logf(format string, args ...any) { f(format, args...) }
