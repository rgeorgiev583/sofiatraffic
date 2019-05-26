package regular

import "strings"

// Arrival represents the event of arrival of an urban transit vehicle.
type Arrival struct {
	Time                   string `json:"time"`                     // estimated time of arrival
	HasAirConditioning     bool   `json:"has_air_conditioning"`     // whether the vehicle has air conditioning
	IsWheelchairAccessible bool   `json:"is_wheelchair_accessible"` // whether the vehicle is wheelchair-accessible
}

// ArrivalList represents the list of expected arrivals of vehicles from a given urban transit line.
type ArrivalList []*Arrival

// DoShowFacilities determines whether info about the available facilities in the vehicles should be displayed for Arrival objects.
var DoShowFacilities bool

func (arrival *Arrival) String() (str string) {
	str = arrival.Time
	if DoShowFacilities {
		var airConditioningStateRepresentation string
		if arrival.HasAirConditioning {
			airConditioningStateRepresentation = "+"
		} else {
			airConditioningStateRepresentation = "-"
		}

		var wheelchairAccessibilityStateRepresentation string
		if arrival.IsWheelchairAccessible {
			wheelchairAccessibilityStateRepresentation = "+"
		} else {
			wheelchairAccessibilityStateRepresentation = "-"
		}

		str += " (" + airConditioningStateRepresentation + "К, " + wheelchairAccessibilityStateRepresentation + "И)"
	}
	return
}

func (arrivalList ArrivalList) String() string {
	arrivalStrings := make([]string, len(arrivalList))
	for i, arrival := range arrivalList {
		arrivalStrings[i] = arrival.String()
	}
	return strings.Join(arrivalStrings, ", ")
}
