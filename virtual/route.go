package virtual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/rgeorgiev583/sofiatraffic/virtual/l10n"
)

// Route represents the list of stops where an urban transit line stops in a specific direction.
type Route struct {
	StopCodes []string `json:"codes"` // numerical codes of the stops
}

// LineNumberRoutes represents the pair of routes for the two directions of an urban transit line with the specified code.
type LineNumberRoutes struct {
	LineNumber string   `json:"name"`   // numerical code of the line
	Routes     []*Route `json:"routes"` // list of routes for the line; should have exactly two elements
}

// VehicleTypeLineNumberRoutes represents the list of LineNumberRoutes objects for urban transit vehicles of the specified type.
type VehicleTypeLineNumberRoutes struct {
	VehicleType string              `json:"type"` // type of the vehicle
	Lines       []*LineNumberRoutes `json:"lines"`
}

// VehicleTypeLineNumberRoutesList represents the list of all VehicleTypeRoutes objects.
type VehicleTypeLineNumberRoutesList []*VehicleTypeLineNumberRoutes // should have as many elements as there are vehicle types (i.e. three)

// NamedRoute represents a route with a name.
type NamedRoute struct {
	Name  string // name of the route
	Stops StopList
}

// LineNamedRoutes represents the list of routes for the urban transit line with the specified VehicleType and LineNumber.
type LineNamedRoutes struct {
	VehicleType string // type of the vehicle
	LineNumber  string // numerical code of the line
	Routes      []*NamedRoute
}

const (
	apiRoutesScheme   = "https"
	apiRoutesHostname = "routes.sofiatraffic.bg"
	apiRoutesPath     = "/resources"
	apiRoutesEndpoint = "/routes.json"
)

// GetName returns the name of a route. The name of a route consists of the names of its first and last stops (which are determined using the StopMap passed as argument).
func (r *Route) GetName(stops StopMap) (name string, err error) {
	if len(r.StopCodes) < 2 {
		err = fmt.Errorf("route should have at least two stops")
		return
	}

	firstStopCode := r.StopCodes[0]
	firstStop, ok := stops[firstStopCode]
	if !ok {
		err = fmt.Errorf("could not determine name for stop with code %s", firstStopCode)
		return
	}

	lastStopCode := r.StopCodes[len(r.StopCodes)-1]
	lastStop, ok := stops[lastStopCode]
	if !ok {
		err = fmt.Errorf("could not determine name for stop with code %s", lastStopCode)
		return
	}

	name = firstStop.Name + " - " + lastStop.Name
	return
}

// GetRoutes fetches and returns the list of all urban transit routes.
func GetRoutes() (routes VehicleTypeLineNumberRoutesList, err error) {
	apiRoutesEndpointURL := &url.URL{
		Scheme: apiRoutesScheme,
		Host:   apiRoutesHostname,
		Path:   apiRoutesPath + apiRoutesEndpoint,
	}
	response, err := http.Get(apiRoutesEndpointURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the API endpoint: %s", err.Error())
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&routes)
	if err != nil {
		err = fmt.Errorf("could not decode JSON data returned by the API endpoint: %s", err.Error())
		return
	}

	return
}

// GetNamedRoutesByLine returns the list of named routes for the urban transit line with the specified vehicleType and lineNumber. The stops argument is used to determine the names of the stops.
func (rl VehicleTypeLineNumberRoutesList) GetNamedRoutesByLine(vehicleType string, lineNumber string, stops StopMap) (namedRoutes *LineNamedRoutes, err error) {
	for _, vehicleTypeRoutes := range rl {
		if vehicleType == "" || vehicleTypeRoutes.VehicleType == vehicleType {
			for _, lineNumberRoutes := range vehicleTypeRoutes.Lines {
				if lineNumber == "" || lineNumberRoutes.LineNumber == lineNumber {
					namedRoutes = &LineNamedRoutes{
						VehicleType: vehicleTypeRoutes.VehicleType,
						LineNumber:  lineNumberRoutes.LineNumber,
						Routes:      make([]*NamedRoute, len(lineNumberRoutes.Routes)),
					}

					for _, route := range lineNumberRoutes.Routes {
						routeName, err := route.GetName(stops)
						if err != nil {
							return namedRoutes, err
						}

						routeStops, err := stops.GetStopsByCodes(route.StopCodes)
						if err != nil {
							return namedRoutes, err
						}

						namedRoutes.Routes = append(namedRoutes.Routes, &NamedRoute{Name: routeName, Stops: routeStops})
					}
					return namedRoutes, err
				}
			}
		}
	}

	return
}

func (rs *LineNamedRoutes) String() string {
	var builder strings.Builder
	lineTitle := l10n.Translator[rs.VehicleType] + " " + rs.LineNumber
	builder.WriteString(lineTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(lineTitle)) + "\n\n")
	for _, route := range rs.Routes {
		builder.WriteString(route.Name + "\n" + strings.Repeat("-", utf8.RuneCountInString(route.Name)) + "\n" + route.Stops.String() + "\n")
	}
	return builder.String()
}
