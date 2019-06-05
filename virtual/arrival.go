package virtual

import (
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/virtual/l10n"
)

// Arrival represents the event of arrival of an urban transit vehicle and describes the facilities in the vehicle.
type Arrival struct {
	Time                   string `json:"time"`                     // estimated time of arrival
	HasAirConditioning     bool   `json:"has_air_conditioning"`     // whether the vehicle has air conditioning
	IsWheelchairAccessible bool   `json:"is_wheelchair_accessible"` // whether the vehicle is wheelchair-accessible
}

// ArrivalList represents the list of expected arrivals of vehicles from a specific urban transit line.
type ArrivalList []*Arrival

// LineVehicleArrivalList represents an urban transit line.
type LineVehicleArrivalList struct {
	VehicleType     string      `json:"vehicle_type"` // type of the vehicle (either "bus", "trolley" or "tram")
	Number          string      `json:"name"`         // numerical code of the line
	VehicleArrivals ArrivalList `json:"arrivals"`
}

// DoShowFacilities determines whether info about the available facilities in the vehicles should be displayed for Arrival objects.
var DoShowFacilities bool

func (a *Arrival) String() (str string) {
	str += a.Time
	if DoShowFacilities {
		str += " ("
		if a.HasAirConditioning {
			str += "+"
		} else {
			str += "-"
		}
		str += l10n.Translator[l10n.AirConditioningAbbreviation] + ", "
		if a.IsWheelchairAccessible {
			str += "+"
		} else {
			str += "-"
		}
		str += l10n.Translator[l10n.WheelchairAccessibilityAbbreviation] + ")"
	}
	return
}

func (al ArrivalList) String() string {
	arrivalStrings := make([]string, len(al))
	for i, arrival := range al {
		arrivalStrings[i] = arrival.String()
	}
	return strings.Join(arrivalStrings, ", ")
}

func (la *LineVehicleArrivalList) String() string {
	return "* " + l10n.Translator[la.VehicleType] + " " + la.Number + ": " + la.VehicleArrivals.String()
}
