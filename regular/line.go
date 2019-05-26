package regular

import "strings"

// Line represents an urban transit line.
type Line struct {
	VehicleType string      `json:"vehicle_type"` // type of the vehicle (either "bus", "trolley" or "tram")
	Code        string      `json:"name"`         // numerical code of the line
	Arrivals    ArrivalList `json:"arrivals"`     // list of expected vehicle arrivals
}

const (
	// VehicleTypeBus represents a bus.
	VehicleTypeBus = "bus"
	// VehicleTypeTrolleybus represents a trolleybus.
	VehicleTypeTrolleybus = "trolley"
	// VehicleTypeTram represents a tram.
	VehicleTypeTram = "tram"
)

// VehicleTypeTranslator translates the names of vehicles from English to the local language.
var VehicleTypeTranslator func(string) string

func (line *Line) String() string {
	var builder strings.Builder
	builder.WriteString("* " + VehicleTypeTranslator(line.VehicleType) + " " + line.Code + ": " + line.Arrivals.String())
	return builder.String()
}
