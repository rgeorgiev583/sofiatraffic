package l10n

// BulgarianTranslator maps names of terms in the reference language (i.e. English) to their translation in Bulgarian.
var BulgarianTranslator = map[string]string{
	VehicleTypeBus:        "автобус",
	VehicleTypeTrolleybus: "тролейбус",
	VehicleTypeTram:       "трамвай",

	AirConditioningAbbreviation:         "К",
	WheelchairAccessibilityAbbreviation: "И",

	GenerationTime: "време на генериране",
}

// ReverseBulgarianTranslator maps translated terms in Bulgarian to their names in the reference language (i.e. English).
var ReverseBulgarianTranslator = map[string]string{
	"автобус":   VehicleTypeBus,
	"тролейбус": VehicleTypeTrolleybus,
	"трамвай":   VehicleTypeTram,
}
