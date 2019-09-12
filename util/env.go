package util

// Env defines for configuration
type (
	// Environment defines public interfaces
	//
	// must be implement in children
	//
	Environment interface {
		// IsProd returns determined product environment mode
		IsProd() bool
		// IsDev returns determined develop environment mode
		IsDev() bool
		// IsDebug returns
		IsDebug() bool
		// IsSentry returns
		IsSentry() bool
		// EnvString gets property variable
		EnvString(prop string) string
		// EnvInt gets property variable
		EnvInt(prop string) int
	}
)
