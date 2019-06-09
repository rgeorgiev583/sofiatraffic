package l10n

// EnglishTranslator maps names of terms in the reference language (i.e. English) to their translation in English.
var EnglishTranslator = map[string]string{
	VehicleTypeBus:        "bus",
	VehicleTypeTrolleybus: "trolleybus",
	VehicleTypeTram:       "tram",
	VehicleTypeMetro:      "metro",

	BusLines:        "bus lines",
	TrolleybusLines: "trolleybus lines",
	TramLines:       "tram lines",

	OperationMode:           "operation mode",
	OperationModeWeekday:    "weekday",
	OperationModePreHoliday: "pre-holiday",
	OperationModeHoliday:    "holiday",

	OnRoute: "on route",
}

// ReverseEnglishTranslator maps translated terms in English to their names in the reference language (i.e. English).
var ReverseEnglishTranslator = map[string]string{
	"bus":        VehicleTypeBus,
	"trolleybus": VehicleTypeTrolleybus,
	"tram":       VehicleTypeTram,
	"metro":      VehicleTypeMetro,

	"weekday":     OperationModeWeekday,
	"pre-holiday": OperationModePreHoliday,
	"holiday":     OperationModeHoliday,
}
