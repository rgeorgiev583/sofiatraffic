package l10n

import (
	"os"
	"strings"
)

const (
	LanguageCodeEnglish   = "en"
	LanguageCodeBulgarian = "bg"
)

// Language indicates the local language (which is determined by the system locale settings).
var Language string

// Translator maps names of terms in the reference language (i.e. English) to their translation in the local language.
var Translator map[string]string

// ReverseTranslator maps translated terms in the local language to their names in the reference language (i.e. English).
var ReverseTranslator map[string]string

// InitTranslatorWithLocaleName initializes the Translator global variable with the appropriate translator for the specified locale.
func InitTranslatorWithLocaleName(locale string) {
	if strings.HasPrefix(locale, LanguageCodeBulgarian+"_") {
		Language = LanguageCodeBulgarian
		Translator = BulgarianTranslator
		ReverseTranslator = ReverseBulgarianTranslator
	} else {
		Language = LanguageCodeEnglish
		Translator = EnglishTranslator
		ReverseTranslator = ReverseEnglishTranslator
	}
}

// InitTranslator initializes the Translator global variable with the appropriate translator for the language defined in the system locale settings.
// Determining of the language from locale settings is only implemented for Unix-like platforms so far (i.e. only POSIX locales are supported). TODO: implement for Windows.
func InitTranslator() {
	locale := os.Getenv("LC_ALL")
	if locale != "" {
		InitTranslatorWithLocaleName(locale)
		return
	}

	locale = os.Getenv("LC_MESSAGES")
	if locale != "" {
		InitTranslatorWithLocaleName(locale)
		return
	}

	locale = os.Getenv("LANG")
	if locale != "" {
		InitTranslatorWithLocaleName(locale)
		return
	}
}
