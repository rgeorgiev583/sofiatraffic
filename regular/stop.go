package regular

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Language determines the language to be used for stop names.
type Language int

type stopArrivalsRepresentation struct {
	Code  string
	Lines []linesArrivalsRepresentation
	Time  string
	Name  string
}

// Stop represents an urban transit stop.
type Stop struct {
	Name string  `json:"n"` // name of the stop
	Code string  `json:"c"` // numerical code of the stop
	X    float64 `json:"x"` // geographic latitude of the stop (in degrees)
	Y    float64 `json:"y"` // geographic longitude of the stop (in degrees)
}

// StopList represents the list of all urban transit stops.
type StopList []*Stop

// StopArrivalContext represents a LineArrivalMap together with its corresponding Stop object.
type StopArrivalContext struct {
	LineArrivalMap
	*Stop
}

// StopArrivalMap represents a map from each stop code to its corresponding StopArrivalContext.
type StopArrivalMap map[string]*StopArrivalContext

const (
	// Bulgarian language.
	Bulgarian Language = iota
	// English language.
	English
)

const (
	apiStopsScheme   = "https"
	apiStopsHostname = "routes.sofiatraffic.bg"
	apiStopsPath     = "/resources"
	apiStopsEndpoint = "/stops-bg.json"
)

func getStopRepresentation(code string, line *Line) (stopRepresentation *stopArrivalsRepresentation, err error) {
	apiArrivalsEndpointURL := &url.URL{
		Scheme: apiArrivalsScheme,
		Host:   apiArrivalsHostname,
		Path:   apiArrivalsPath + apiArrivalsEndpoint + "/" + code + "/",
	}
	query := url.Values{}
	if line != nil {
		if line.Code != "" {
			query.Set("line", line.Code)
		}
		if line.VehicleType != "" {
			query.Set("type", line.VehicleType)
		}
	}
	apiArrivalsEndpointURL.RawQuery = query.Encode()
	response, err := http.Get(apiArrivalsEndpointURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the API endpoint: %s", err.Error())
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	stopRepresentation = &stopArrivalsRepresentation{}
	err = decoder.Decode(stopRepresentation)
	if err != nil {
		stopRepresentation = nil
		err = fmt.Errorf("could not decode JSON data returned by the API endpoint: %s", err.Error())
		return
	}

	return
}

// GetStops fetches the list of bus/trolleybus/tram stops from the API endpoint.
func GetStops() (stops StopList, err error) {
	apiStopsEndpointURL := &url.URL{
		Scheme: apiStopsScheme,
		Host:   apiStopsHostname,
		Path:   apiStopsPath + apiStopsEndpoint,
	}
	response, err := http.Get(apiStopsEndpointURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the API endpoint: %s", err.Error())
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&stops)
	if err != nil {
		err = fmt.Errorf("could not decode JSON data returned by the API endpoint: %s", err.Error())
		return
	}

	return
}

// GetArrivalsByStopCodeAndLine returns a list containing all expected vehicle arrivals at the stop with the given code. If the line argument is non-nil, only arrivals of vehicles from the given type OR from the given line will be listed.
func GetArrivalsByStopCodeAndLine(stopCode string, line *Line) (lineArrivalMap LineArrivalMap, err error) {
	stopRepresentation, err := getStopRepresentation(stopCode, nil)
	if err != nil {
		err = fmt.Errorf("could not decode the representation of the arrivals at the stop returned by the API endpoint: %s", err.Error())
		return
	}

	lineArrivalMap = LineArrivalMap{}
	for _, linesRepresentation := range stopRepresentation.Lines {
		arrivals := make(ArrivalList, len(linesRepresentation.Arrivals))
		for i, arrivalsRepresentation := range linesRepresentation.Arrivals {
			arrival := &Arrival{
				Time: arrivalsRepresentation.Time,
				Vehicle: &VehicleFacilities{
					HasAirConditioning:     arrivalsRepresentation.HasAirConditioning,
					IsWheelchairAccessible: arrivalsRepresentation.IsWheelchairAccessible,
				},
			}
			arrivals[i] = arrival
		}
		line := Line{
			VehicleType: linesRepresentation.VehicleType,
			Code:        linesRepresentation.Code,
		}
		lineArrivalMap[line] = arrivals
	}

	return
}

// GetArrivalsByStopNameAndLine is like GetArrivalsByStopCodeAndLine but the stop is determined by its name instead of its code.
func (stopList StopList) GetArrivalsByStopNameAndLine(stopName string, line *Line) (stopArrivalMap StopArrivalMap, err error) {
	stopArrivalMap = StopArrivalMap{}
	for _, stop := range stopList {
		if stop.Name == stopName {
			lineArrivalMap, err := GetArrivalsByStopCodeAndLine(stop.Code, line)
			if err != nil {
				break
			}

			stopArrivalContext := &StopArrivalContext{}
			stopArrivalContext.LineArrivalMap = lineArrivalMap
			stopArrivalContext.Stop = stop
			stopArrivalMap[stop.Code] = stopArrivalContext
		}
	}

	return
}

// MatchArrivalsByStopNameAndLine is like GetArrivalsByStopNameAndLine but it performs case-insensitive matching by stopNamePattern on the stop name.
func (stopList StopList) MatchArrivalsByStopNameAndLine(stopNamePattern string, line *Line) (stopArrivalMap StopArrivalMap, err error) {
	stopNamePattern = strings.ToUpper(stopNamePattern)
	stopArrivalMap = StopArrivalMap{}
	for _, stop := range stopList {
		if strings.Contains(stop.Name, stopNamePattern) {
			lineArrivalMap, err := GetArrivalsByStopCodeAndLine(stop.Code, line)
			if err != nil {
				break
			}

			stopArrivalContext := &StopArrivalContext{}
			stopArrivalContext.LineArrivalMap = lineArrivalMap
			stopArrivalContext.Stop = stop
			stopArrivalMap[stop.Code] = stopArrivalContext
		}
	}

	return
}

func (stopArrivalMap StopArrivalMap) String() (str string) {
	var builder strings.Builder
	for _, stopArrivalContext := range stopArrivalMap {
		if len(stopArrivalContext.LineArrivalMap) == 0 {
			continue
		}

		stopTitle := stopArrivalContext.Stop.Name + " (" + stopArrivalContext.Stop.Code + "):"
		builder.WriteString(stopTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(stopTitle)) + "\n")
		builder.WriteString(stopArrivalContext.LineArrivalMap.String() + "\n")
	}
	return builder.String()
}