package regular

import "strings"

type arrivalsArrivalsRepresentation struct {
	Time                   string `json:"time"`
	HasAirConditioning     bool   `json:"has_air_conditioning"`
	IsWheelchairAccessible bool   `json:"is_wheelchair_accessible"`
}

// VehicleFacilities represents the set of facilities available in a vehicle.
type VehicleFacilities struct {
	HasAirConditioning     bool // whether the vehicle is air-conditioned
	IsWheelchairAccessible bool // whether the vehicle is wheelchair-accessible
}

// Arrival represents the event of arrival of an urban transit vehicle.
type Arrival struct {
	Time string // estimated time of arrival
	*VehicleFacilities
}

// ArrivalList represents the list of arrivals of an urban transit vehicle.
type ArrivalList []*Arrival

const (
	apiArrivalsScheme   = "https"
	apiArrivalsHostname = "api-arrivals.sofiatraffic.bg"
	apiArrivalsPath     = "/api/v1"
	apiArrivalsEndpoint = "/arrivals"
)

func (arrival *Arrival) String() string {
	return arrival.Time
}

func (arrivalList ArrivalList) String() string {
	arrivalStrings := make([]string, len(arrivalList))
	for i, arrival := range arrivalList {
		arrivalStrings[i] = arrival.String()
	}
	return strings.Join(arrivalStrings, ", ")
}
