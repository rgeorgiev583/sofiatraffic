package l10n

import (
	"github.com/rgeorgiev583/sofiatraffic/l10n"
)

// Translator maps names of terms in the reference language (i.e. English) to their translation in the local language.
var Translator map[string]string

// InitTranslator initializes the Translator global variable with the appropriate translator for the language defined in the system locale settings.
func InitTranslator() {
	switch l10n.Language {
	case l10n.LanguageCodeBulgarian:
		Translator = BulgarianTranslator

	case l10n.LanguageCodeEnglish:
		Translator = EnglishTranslator
	}
}
