package regular

// Line represents an urban transit line.
type Line struct {
	VehicleType string      `json:"vehicle_type"` // type of the vehicle (either "bus", "trolley" or "tram")
	Code        string      `json:"name"`         // numerical code of the line
	Arrivals    ArrivalList `json:"arrivals"`
}

const (
	// VehicleTypeBus represents a bus.
	VehicleTypeBus = "bus"
	// VehicleTypeTrolleybus represents a trolleybus.
	VehicleTypeTrolleybus = "trolley"
	// VehicleTypeTram represents a tram.
	VehicleTypeTram = "tram"
)

// VehicleTypeTranslator translates vehicle types from English to the local language.
var VehicleTypeTranslator func(string) string

func (l *Line) String() string {
	return "* " + VehicleTypeTranslator(l.VehicleType) + " " + l.Code + ": " + l.Arrivals.String()
}
