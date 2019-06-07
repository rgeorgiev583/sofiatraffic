package i18n

import (
	"os"
	"strings"
)

const (
	// LanguageCodeEnglish represents the English language.
	LanguageCodeEnglish = "en"
	// LanguageCodeBulgarian represents the Bulgarian language.
	LanguageCodeBulgarian = "bg"
)

// Language indicates the local language (which is determined by the system locale settings).
var Language string

// InitWithLocaleName initializes the internationalization subsystem to use the appropriate local language as determined by the specified locale.
func InitWithLocaleName(locale string) {
	if strings.HasPrefix(locale, LanguageCodeBulgarian+"_") {
		Language = LanguageCodeBulgarian
	} else {
		Language = LanguageCodeEnglish
	}
}

// Init initializes the internationalization subsystem to use the appropriate local language as determined by the system locale settings.
// Determining of the language from locale settings is only implemented for Unix-like platforms so far (i.e. only POSIX locales are supported). TODO: implement for Windows.
func Init() {
	locale := os.Getenv("LC_ALL")
	if locale != "" {
		InitWithLocaleName(locale)
		return
	}

	locale = os.Getenv("LC_MESSAGES")
	if locale != "" {
		InitWithLocaleName(locale)
		return
	}

	locale = os.Getenv("LANG")
	if locale != "" {
		InitWithLocaleName(locale)
		return
	}
}
