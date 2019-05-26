package regular

import "strings"

type linesArrivalsRepresentation struct {
	VehicleType string                           `json:"vehicle_type"`
	Arrivals    []arrivalsArrivalsRepresentation `json:"arrivals"`
	Code        string                           `json:"name"`
}

// Line represents and uniquely identifies an urban transit line.
type Line struct {
	VehicleType string // type of the vehicle (either "bus", "trolley" or "tram")
	Code        string // numerical code of the line
}

// LineArrivalMap represents a map from each urban transit line to the chronologically-ordered list of vehicle arrivals at a specific stop for that line.
type LineArrivalMap map[Line]ArrivalList

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

func (line *Line) String() (str string) {
	str += VehicleTypeTranslator(line.VehicleType) + " " + line.Code
	return
}

func (lineArrivalMap LineArrivalMap) String() string {
	var builder strings.Builder
	for line, arrivals := range lineArrivalMap {
		builder.WriteString("* " + line.String() + ": " + arrivals.String() + "\n")
	}
	return builder.String()
}
