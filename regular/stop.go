package regular

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Stop represents an urban transit stop.
type Stop struct {
	Code string `json:"c"` // numerical code of the stop
	Name string `json:"n"` // name of the stop
}

// StopList represents the list of all urban transit stops.
type StopList []*Stop

const (
	apiStopsScheme   = "https"
	apiStopsHostname = "routes.sofiatraffic.bg"
	apiStopsPath     = "/resources"
	apiStopsEndpoint = "/stops-bg.json"
)

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

func (sl StopList) String() string {
	var builder strings.Builder
	for i, stop := range sl {
		builder.WriteString(strconv.Itoa(i+1) + ". " + stop.Name + " (" + stop.Code + ")\n")
	}
	return builder.String()
}
