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

// VehicleTypeRoutes represents the list of LineNumberRoutes objects for urban transit vehicles of the specified type.
type VehicleTypeRoutes struct {
	VehicleType string              `json:"type"` // type of the vehicle
	Lines       []*LineNumberRoutes `json:"lines"`
}

// VehicleTypeRoutesList represents the list of all VehicleTypeRoutes objects.
type VehicleTypeRoutesList []*VehicleTypeRoutes // should have as many elements as there are vehicle types (i.e. three)

// NamedRoute represents a route with a name.
type NamedRoute struct {
	Name  string // name of the route
	Stops StopList
}

// LineNamedRoutes represents the pair of routes for the two directions of the urban transit line with the specified VehicleType and LineNumber.
type LineNamedRoutes struct {
	VehicleType                               string // type of the vehicle
	LineNumber                                string // numerical code of the line
	FirstDirectionRoute, SecondDirectionRoute *NamedRoute
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
func GetRoutes() (routes VehicleTypeRoutesList, err error) {
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
func (rl VehicleTypeRoutesList) GetNamedRoutesByLine(vehicleType string, lineNumber string, stops StopMap) (routes *LineNamedRoutes, err error) {
	for _, vehicleTypeRoutes := range rl {
		if vehicleType == "" || vehicleTypeRoutes.VehicleType == vehicleType {
			for _, lineNumberRoutes := range vehicleTypeRoutes.Lines {
				if lineNumber == "" || lineNumberRoutes.LineNumber == lineNumber {
					if len(lineNumberRoutes.Routes) != 2 {
						err = fmt.Errorf("there should be exactly two routes for a line")
						return
					}

					firstDirectionRoute := lineNumberRoutes.Routes[0]
					secondDirectionRoute := lineNumberRoutes.Routes[1]

					firstDirectionRouteName, err := firstDirectionRoute.GetName(stops)
					if err != nil {
						return routes, err
					}

					secondDirectionRouteName, err := secondDirectionRoute.GetName(stops)
					if err != nil {
						return routes, err
					}

					firstDirectionRouteStops, err := stops.GetStopsByCodes(firstDirectionRoute.StopCodes)
					if err != nil {
						return routes, err
					}

					secondDirectionRouteStops, err := stops.GetStopsByCodes(secondDirectionRoute.StopCodes)
					if err != nil {
						return routes, err
					}

					routes = &LineNamedRoutes{
						VehicleType:          vehicleTypeRoutes.VehicleType,
						LineNumber:           lineNumberRoutes.LineNumber,
						FirstDirectionRoute:  &NamedRoute{Name: firstDirectionRouteName, Stops: firstDirectionRouteStops},
						SecondDirectionRoute: &NamedRoute{Name: secondDirectionRouteName, Stops: secondDirectionRouteStops},
					}
					return routes, err
				}
			}
		}
	}

	return
}

func (rs *LineNamedRoutes) String() (str string) {
	lineTitle := l10n.Translator[rs.VehicleType] + " " + rs.LineNumber
	str += lineTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(lineTitle)) + "\n\n"
	str += rs.FirstDirectionRoute.Name + "\n" + strings.Repeat("-", utf8.RuneCountInString(rs.FirstDirectionRoute.Name)) + "\n" + rs.FirstDirectionRoute.Stops.String() + "\n"
	str += rs.SecondDirectionRoute.Name + "\n" + strings.Repeat("-", utf8.RuneCountInString(rs.SecondDirectionRoute.Name)) + "\n" + rs.SecondDirectionRoute.Stops.String() + "\n"
	return
}
