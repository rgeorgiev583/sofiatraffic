package regular

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/rgeorgiev583/sofiatraffic/l10n"
)

// StopTimetable represents the list of all expected urban transit vehicle arrivals at a specific stop.
type StopTimetable struct {
	Code  string  `json:"code"`                 // numerical code of the stop
	Name  string  `json:"name"`                 // name of the stop
	Lines []*Line `json:"lines"`                // list of lines which stop at the stop
	Time  string  `json:"timestamp_calculated"` // time at which the timetable was generated
}

// StopTimetableList represents a list of urban transit stop timetables.
type StopTimetableList []*StopTimetable

// StopTimetableFetchResult represents the result of attempting to fetch an urban transit stop timetable.
type StopTimetableFetchResult struct {
	Timetable *StopTimetable
	Err       error
}

// StopTimetableChannel is to be used for asynchronous processing of fetched urban transit stop timetables.
type StopTimetableChannel <-chan *StopTimetableFetchResult

const (
	apiArrivalsScheme   = "https"
	apiArrivalsHostname = "api-arrivals.sofiatraffic.bg"
	apiArrivalsPath     = "/api/v1"
	apiArrivalsEndpoint = "/arrivals"
)

// DoShowGenerationTimeForTimetables determines whether the generation time of an urban transit stop timetable should be included in its display representation.
var DoShowGenerationTimeForTimetables bool

// GetTimetableByStopCodeAndLine fetches and returns the timetable for the urban transit stop with the specified code. If the vehicleType argument is non-empty, only arrivals of vehicles of the specified type will be listed. If the lineCode argument is non-empty, only arrivals of vehicles from the line with the specified code will be listed.
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
		err = fmt.Errorf("could not decode JSON data returned by the API endpoint: %s", err.Error())
		return
	}

	return
}

// GetTimetablesByStopNameAndLine fetches and returns a list containing all timetables for urban transit stops with the specified name. The vehicleType and lineCode arguments behave as in GetTimetableByStopCodeAndLine. The isExactMatch argument determines whether the specified stopName should be matched exactly or as a substring.
func (sl StopList) GetTimetablesByStopNameAndLine(stopName string, vehicleType string, lineCode string, isExactMatch bool) (timetables StopTimetableList, err error) {
	if !isExactMatch {
		stopName = strings.ToUpper(stopName)
	}
	timetables = StopTimetableList{}
	for _, stop := range sl {
		if isExactMatch && stop.Name == stopName || !isExactMatch && strings.Contains(stop.Name, stopName) {
			timetable, err := GetTimetableByStopCodeAndLine(stop.Code, vehicleType, lineCode)
			if err != nil {
				return timetables, err
			}

			timetables = append(timetables, timetable)
		}
	}
	return
}

// GetTimetablesByStopNameAndLineAsync is the asynchronous version of GetTimetableByStopCodeAndLine.
func (sl StopList) GetTimetablesByStopNameAndLineAsync(stopName string, vehicleType string, lineCode string, isExactMatch bool) (timetables StopTimetableChannel) {
	if !isExactMatch {
		stopName = strings.ToUpper(stopName)
	}
	fetchResults := make(chan *StopTimetableFetchResult)
	timetables = fetchResults
	var timetableFetchers sync.WaitGroup
	for _, stop := range sl {
		if isExactMatch && stop.Name == stopName || !isExactMatch && strings.Contains(stop.Name, stopName) {
			timetableFetchers.Add(1)
			go func(stop *Stop) {
				timetable, err := GetTimetableByStopCodeAndLine(stop.Code, vehicleType, lineCode)
				fetchResults <- &StopTimetableFetchResult{Timetable: timetable, Err: err}
				timetableFetchers.Done()
			}(stop)
		}
	}
	go func() {
		timetableFetchers.Wait()
		close(fetchResults)
	}()
	return
}

func (t *StopTimetable) String() string {
	var builder strings.Builder
	stopTitle := t.Name + " (" + t.Code + ")"
	builder.WriteString(stopTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(stopTitle)) + "\n")
	if DoShowGenerationTimeForTimetables {
		builder.WriteString("(" + l10n.Translator[l10n.GenerationTime] + ": " + t.Time + ")\n")
	}
	for _, line := range t.Lines {
		builder.WriteString(line.String() + "\n")
	}
	return builder.String()
}

func (tl StopTimetableList) String() string {
	var builder strings.Builder
	for _, timetable := range tl {
		if len(timetable.Lines) == 0 {
			continue
		}

		builder.WriteString(timetable.String() + "\n")
	}
	return builder.String()
}

func (tc StopTimetableChannel) String() string {
	var builder strings.Builder
	for fetchResult := range tc {
		if fetchResult.Err != nil {
			log.Println(fetchResult.Err.Error())
			continue
		}

		if len(fetchResult.Timetable.Lines) == 0 {
			continue
		}

		builder.WriteString(fetchResult.Timetable.String() + "\n")
	}
	return builder.String()
}
