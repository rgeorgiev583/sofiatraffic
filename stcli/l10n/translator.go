package l10n

import "github.com/rgeorgiev583/sofiatraffic/i18n"

// Translator maps names of terms in the reference language (i.e. English) to their translation in the local language.
var Translator map[string]string

// InitTranslator initializes the Translator global variable with the appropriate translator for the local language.
func InitTranslator() {
	switch i18n.Language {
	case i18n.LanguageCodeBulgarian:
		Translator = BulgarianTranslator

	case i18n.LanguageCodeEnglish:
		Translator = EnglishTranslator
	}
}
