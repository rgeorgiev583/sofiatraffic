package l10n

// BulgarianTranslator maps names of terms in the reference language (i.e. English) to their translation in Bulgarian.
var BulgarianTranslator = map[string]string{
	VehicleTypeBus:        "автобус",
	VehicleTypeTrolleybus: "тролейбус",
	VehicleTypeTram:       "трамвай",
	VehicleTypeMetro:      "метро",

	BusLines:        "автобусни линии",
	TrolleybusLines: "тролейбусни линии",
	TramLines:       "трамвайни линии",

	OperationMode:           "режим",
	OperationModeWeekday:    "делник",
	OperationModePreHoliday: "предпразник",
	OperationModeHoliday:    "празник",

	OnRoute: "по маршрут",
}

// ReverseBulgarianTranslator maps translated terms in Bulgarian to their names in the reference language (i.e. English).
var ReverseBulgarianTranslator = map[string]string{
	"автобус":   VehicleTypeBus,
	"тролейбус": VehicleTypeTrolleybus,
	"трамвай":   VehicleTypeTram,
	"метро":     VehicleTypeMetro,

	"делник":      OperationModeWeekday,
	"предпразник": OperationModePreHoliday,
	"празник":     OperationModeHoliday,
}
