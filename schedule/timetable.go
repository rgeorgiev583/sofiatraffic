package schedule

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html/atom"

	"golang.org/x/net/html"
)

// Timetable represents a list of urban transit vehicle arrival times.
type Timetable []string

type timetableScannerState int

const (
	timetablePagePath = "/server/html/schedule_load"

	timetableScannerNotInsideRelevantElement timetableScannerState = iota
	timetableScannerInsideHoursCellDiv
	timetableScannerInsideHoursCellAnchor
)

// GetTimetable fetches and returns the urban transit stop timetable matching the specified operationModeCode, routeCode and stopCode.
func GetTimetable(operationModeCode string, routeCode string, stopCode string) (timetable Timetable, err error) {
	timetablePageURL := &url.URL{
		Scheme: Scheme,
		Host:   Hostname,
		Path:   timetablePagePath + "/" + operationModeCode + "/" + routeCode + "/" + stopCode,
	}
	response, err := http.Get(timetablePageURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the schedule timetable page: %s", err.Error())
		return
	}
	defer response.Body.Close()

	var state timetableScannerState
	for tokenizer := html.NewTokenizer(response.Body); tokenizer.Next() != html.ErrorToken; {
		token := tokenizer.Token()
		switch token.Type {
		case html.StartTagToken:
			switch token.DataAtom {
			case atom.Div:
				var class string
				for _, attr := range token.Attr {
					if atom.Lookup([]byte(attr.Key)) == atom.Class {
						class = attr.Val
						break
					}
				}
				if strings.Contains(class, "hours_cell") {
					state = timetableScannerInsideHoursCellDiv
				}

			case atom.A:
				if state == timetableScannerInsideHoursCellDiv {
					state = timetableScannerInsideHoursCellAnchor
				}
			}

		case html.EndTagToken:
			switch token.DataAtom {
			case atom.A:
				if state == timetableScannerInsideHoursCellAnchor {
					state = timetableScannerInsideHoursCellDiv
				}

			case atom.Div:
				state = timetableScannerNotInsideRelevantElement
			}

		case html.TextToken:
			if state == timetableScannerInsideHoursCellAnchor {
				timetable = append(timetable, token.Data)
			}
		}
	}
	return
}

func (t Timetable) String() string {
	return strings.Join(t, ", ")
}
