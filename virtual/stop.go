package virtual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/i18n"
)

// Stop represents an urban transit stop.
type Stop struct {
	Code string `json:"c"` // numerical code of the stop
	Name string `json:"n"` // name of the stop
}

// StopList represents a list of urban transit stops.
type StopList []*Stop

// StopMap represents a map from the code of each urban transit stop to its corresponding Stop object.
type StopMap map[string]*Stop

const (
	apiStopsScheme            = "https"
	apiStopsHostname          = "routes.sofiatraffic.bg"
	apiStopsPath              = "/resources"
	apiStopsEndpointBulgarian = "/stops-bg.json"
	apiStopsEndpointEnglish   = "/stops-en.json"
)

// DoTranslateStopNames determines whether stop names should be translated from Bulgarian to the local language.
var DoTranslateStopNames bool

// GetStops fetches and returns the list of all urban transit stops.
func GetStops() (stops StopList, err error) {
	apiStopsEndpoint := apiStopsEndpointBulgarian
	if DoTranslateStopNames && i18n.Language == i18n.LanguageCodeEnglish {
		apiStopsEndpoint = apiStopsEndpointEnglish
	}
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

// GetStopMap returns a StopMap object containing all stops in the StopList.
func (sl StopList) GetStopMap() (stops StopMap) {
	stops = StopMap{}
	for _, stop := range sl {
		stops[stop.Code] = stop
	}
	return
}

// GetStopsByCodes returns a list containing the stops with the given codes.
func (sm StopMap) GetStopsByCodes(codes []string) (stops StopList, err error) {
	stops = make(StopList, len(codes))
	for i, code := range codes {
		stop, ok := sm[code]
		if !ok {
			err = fmt.Errorf("could not find stop code %s in stop map", code)
			return stops, err
		}

		stops[i] = stop
	}
	return
}

func (s *Stop) String() string {
	return s.Name + " (" + s.Code + ")"
}

func (sl StopList) String() string {
	var builder strings.Builder
	for i, stop := range sl {
		builder.WriteString(strconv.Itoa(i+1) + ". " + stop.String() + "\n")
	}
	return builder.String()
}

func (sl StopList) Len() int {
	return len(sl)
}

func (sl StopList) Less(i, j int) bool {
	return strings.Compare(sl[i].Code, sl[j].Code) == -1
}

func (sl StopList) Swap(i, j int) {
	temp := sl[i]
	sl[i] = sl[j]
	sl[j] = temp
}
