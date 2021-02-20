package logger

type (
	// SentryErrorLogger is satisfied of elastic.Logger
	SentryErrorLogger struct{}
	// SentryInfoLogger is satisfied of elastic.Logger
	SentryInfoLogger struct{}
)

// Printf prints out message as error
func (a *SentryErrorLogger) Printf(format string, v ...interface{}) {
	errdeps(4, format, v...)
}

// Printf prints out message as info
func (a *SentryInfoLogger) Printf(format string, v ...interface{}) {
	infodeps(4, format, v...)
}
