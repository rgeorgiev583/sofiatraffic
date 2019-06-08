package virtual

import "github.com/rgeorgiev583/sofiatraffic/virtual/l10n"

const (
	// VehicleTypeBus represents a bus.
	VehicleTypeBus = "bus"
	// VehicleTypeTrolleybus represents a trolleybus.
	VehicleTypeTrolleybus = "trolley"
	// VehicleTypeTram represents a tram.
	VehicleTypeTram = "tram"
)

// Line represents an urban transit line.
type Line struct {
	VehicleType, LineNumber string
}

func (l *Line) String() string {
	return l10n.Translator[l.VehicleType] + " " + l.LineNumber
}
