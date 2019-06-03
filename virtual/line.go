package virtual

import "github.com/rgeorgiev583/sofiatraffic/virtual/l10n"

// Line represents an urban transit line.
type Line struct {
	VehicleType     string      `json:"vehicle_type"` // type of the vehicle (either "bus", "trolley" or "tram")
	Number          string      `json:"name"`         // numerical code of the line
	VehicleArrivals ArrivalList `json:"arrivals"`
}

const (
	// VehicleTypeBus represents a bus.
	VehicleTypeBus = "bus"
	// VehicleTypeTrolleybus represents a trolleybus.
	VehicleTypeTrolleybus = "trolley"
	// VehicleTypeTram represents a tram.
	VehicleTypeTram = "tram"
)

func (l *Line) String() string {
	return "* " + l10n.Translator[l.VehicleType] + " " + l.Number + ": " + l.VehicleArrivals.String()
}
