package virtual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Route represents the list of stops where an urban transit line stops when traveling in a specific direction.
type Route struct {
	StopCodes []string `json:"codes"` // numerical codes of the stops
}

// RouteList represents a list of routes.
type RouteList []*Route

// LineNumberRouteList represents the list of routes for an urban transit line with the specified number.
type LineNumberRouteList struct {
	LineNumber string `json:"name"` // number of the line
	RouteList  `json:"routes"`
}

// LineNumberRouteListList represents a list of LineNumberRouteList objects.
type LineNumberRouteListList []*LineNumberRouteList

// VehicleTypeLineNumberRouteListList represents the list of LineNumberRouteList objects for urban transit vehicles of the specified type.
type VehicleTypeLineNumberRouteListList struct {
	VehicleType             string `json:"type"` // type of the vehicle
	LineNumberRouteListList `json:"lines"`
}

// VehicleTypeLineNumberRouteListListList represents the list of all VehicleTypeLineNumberRouteListList objects.
type VehicleTypeLineNumberRouteListListList []*VehicleTypeLineNumberRouteListList // should have as many elements as there are vehicle types (i.e. three)

// NamedRoute represents a route with a name.
type NamedRoute struct {
	Name string
	StopList
}

// NamedRouteList represents a list of NamedRoute objects.
type NamedRouteList []*NamedRoute

// LineNamedRouteList represents the list of named routes for the urban transit line with the specified VehicleType and LineNumber.
type LineNamedRouteList struct {
	*Line
	NamedRouteList
}

// LineNamedRouteListList represents a list of LineNamedRouteList objects.
type LineNamedRouteListList []*LineNamedRouteList

// LineNamedRouteListMap represents a map from an urban transit line to a list of named routes.
type LineNamedRouteListMap map[Line]NamedRouteList

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

// GetRoutes fetches and returns the VehicleTypeLineNumberRouteListListList of all urban transit routes.
func GetRoutes() (routes VehicleTypeLineNumberRouteListListList, err error) {
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

// GetNamedRoutesByLine returns the list of named routes for the urban transit line with the specified vehicleType and lineNumber wrapped in a list (or, alternatively, for all lines matching the other criterion if one of them is empty; or for all lines if both are empty). The stops argument is used to determine the names of the stops.
func (rl VehicleTypeLineNumberRouteListListList) GetNamedRoutesByLine(vehicleType string, lineNumber string, stops StopMap) (namedRouteListList LineNamedRouteListList, err error) {
	namedRouteListList = LineNamedRouteListList{}
	for _, vehicleTypeRoutes := range rl {
		if vehicleType == "" || vehicleTypeRoutes.VehicleType == vehicleType {
			for _, lineNumberRoutes := range vehicleTypeRoutes.LineNumberRouteListList {
				if lineNumber == "" || lineNumberRoutes.LineNumber == lineNumber {
					namedRouteList := &LineNamedRouteList{
						Line:           &Line{VehicleType: vehicleTypeRoutes.VehicleType, LineNumber: lineNumberRoutes.LineNumber},
						NamedRouteList: make([]*NamedRoute, len(lineNumberRoutes.RouteList)),
					}
					for i, route := range lineNumberRoutes.RouteList {
						routeName, err := route.GetName(stops)
						if err != nil {
							return namedRouteListList, err
						}

						routeStops, err := stops.GetStopsByCodes(route.StopCodes)
						if err != nil {
							return namedRouteListList, err
						}

						namedRouteList.NamedRouteList[i] = &NamedRoute{Name: routeName, StopList: routeStops}
					}
					namedRouteListList = append(namedRouteListList, namedRouteList)
				}
			}
		}
	}
	return
}

// GetRouteMap returns a LineNamedRouteList map containing all routes in the VehicleTypeLineNumberRouteListListList object.
func (rl VehicleTypeLineNumberRouteListListList) GetRouteMap(stops StopMap) (routes LineNamedRouteListMap, err error) {
	routes = LineNamedRouteListMap{}
	for _, vehicleTypeRoutes := range rl {
		for _, lineNumberRoutes := range vehicleTypeRoutes.LineNumberRouteListList {
			line := Line{VehicleType: vehicleTypeRoutes.VehicleType, LineNumber: lineNumberRoutes.LineNumber}
			namedRouteList := make([]*NamedRoute, len(lineNumberRoutes.RouteList))
			for _, route := range lineNumberRoutes.RouteList {
				routeName, err := route.GetName(stops)
				if err != nil {
					return routes, err
				}

				routeStops, err := stops.GetStopsByCodes(route.StopCodes)
				if err != nil {
					return routes, err
				}

				namedRouteList = append(namedRouteList, &NamedRoute{Name: routeName, StopList: routeStops})
			}
			routes[line] = namedRouteList
		}
	}
	return
}

func (nr *NamedRoute) String() string {
	return nr.Name + "\n" + strings.Repeat("-", utf8.RuneCountInString(nr.Name)) + "\n" + nr.StopList.String()
}

func (nrl NamedRouteList) String() string {
	var builder strings.Builder
	for _, namedRoute := range nrl {
		builder.WriteString(namedRoute.String() + "\n")
	}
	return builder.String()
}

func (lnrl *LineNamedRouteList) String() (str string) {
	lineTitle := lnrl.Line.String()
	str += lineTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(lineTitle)) + "\n\n"
	str += lnrl.NamedRouteList.String()
	return
}

func (lnrll LineNamedRouteListList) String() string {
	var builder strings.Builder
	for _, namedRouteList := range lnrll {
		builder.WriteString(namedRouteList.String() + "\n")
	}
	return builder.String()
}
