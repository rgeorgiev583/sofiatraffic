package l10n

// EnglishTranslator maps names of terms in the reference language (i.e. English) to their translation in English.
var EnglishTranslator = map[string]string{
	VehicleTypeBus:        "bus",
	VehicleTypeTrolleybus: "trolleybus",
	VehicleTypeTram:       "tram",

	AirConditioningAbbreviation:         "A",
	WheelchairAccessibilityAbbreviation: "W",

	GenerationTime: "generation time",
}

// ReverseEnglishTranslator maps translated terms in English to their names in the reference language (i.e. English).
var ReverseEnglishTranslator = map[string]string{
	"bus":        VehicleTypeBus,
	"trolleybus": VehicleTypeTrolleybus,
	"tram":       VehicleTypeTram,
}
