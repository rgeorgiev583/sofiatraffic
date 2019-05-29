package regular

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// StopTimetable represents the list of all expected urban transit vehicle arrivals at a specific stop.
type StopTimetable struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Lines []*Line `json:"lines"`
	Time  string  `json:"timestamp_calculated"`
}

// StopTimetableList represents the list of all urban transit stop timetables.
type StopTimetableList []*StopTimetable

const (
	apiArrivalsScheme   = "https"
	apiArrivalsHostname = "api-arrivals.sofiatraffic.bg"
	apiArrivalsPath     = "/api/v1"
	apiArrivalsEndpoint = "/arrivals"
)

// DoShowGenerationTimeForTimetables determines whether the generation Time of a StopTimetable object should be included in its display representation.
var DoShowGenerationTimeForTimetables bool

// GenerationTimeLabel determines the label which should be displayed for the generation Time of a StopTimetable object.
var GenerationTimeLabel string

// GetTimetableByStopCodeAndLine returns the StopTimetable for the stop with the given code. If the vehicleType argument is non-empty, only arrivals of vehicles of the given type will be listed. If the lineCode argument is non-empty, only arrivals of vehicles from the line with the given code will be listed.
func GetTimetableByStopCodeAndLine(stopCode string, vehicleType string, lineCode string) (stopTimetable *StopTimetable, err error) {
	apiArrivalsEndpointURL := &url.URL{
		Scheme: apiArrivalsScheme,
		Host:   apiArrivalsHostname,
		Path:   apiArrivalsPath + apiArrivalsEndpoint + "/" + stopCode + "/",
	}
	query := url.Values{}
	if lineCode != "" {
		query.Set("line", lineCode)
	}
	if vehicleType != "" {
		query.Set("type", vehicleType)
	}
	apiArrivalsEndpointURL.RawQuery = query.Encode()
	response, err := http.Get(apiArrivalsEndpointURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the API endpoint: %s", err.Error())
		return
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	stopTimetable = &StopTimetable{}
	err = decoder.Decode(stopTimetable)
	if err != nil {
		stopTimetable = nil
		err = fmt.Errorf("could not decode JSON data returned by the API endpoint: %s", err.Error())
		return
	}

	return
}

// GetTimetablesByStopNameAndLine returns a StopTimetableList containing all timetables for stops with the given name. The vehicleType and lineCode arguments behave as in GetTimetableByStopCodeAndLine.
func (sl StopList) GetTimetablesByStopNameAndLine(stopName string, vehicleType string, lineCode string, isExactMatch bool) (timetables StopTimetableList, err error) {
	if !isExactMatch {
		stopName = strings.ToUpper(stopName)
	}
	timetables = StopTimetableList{}
	for _, stop := range sl {
		if isExactMatch && stop.Name == stopName || !isExactMatch && strings.Contains(stop.Name, stopName) {
			timetable, err := GetTimetableByStopCodeAndLine(stop.Code, vehicleType, lineCode)
			if err != nil {
				break
			}

			timetables = append(timetables, timetable)
		}
	}
	return
}

func (st *StopTimetable) String() string {
	var builder strings.Builder
	stopTitle := st.Name + " (" + st.Code + ")"
	builder.WriteString(stopTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(stopTitle)) + "\n")
	if DoShowGenerationTimeForTimetables {
		builder.WriteString("(" + GenerationTimeLabel + ": " + st.Time + ")\n")
	}
	for _, line := range st.Lines {
		builder.WriteString(line.String() + "\n")
	}
	return builder.String()
}

func (stl StopTimetableList) String() string {
	var builder strings.Builder
	for _, timetable := range stl {
		if len(timetable.Lines) == 0 {
			continue
		}

		builder.WriteString(timetable.String() + "\n")
	}
	return builder.String()
}
