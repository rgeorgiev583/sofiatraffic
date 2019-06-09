package schedule

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/rgeorgiev583/sofiatraffic/schedule/l10n"
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

// DoShowOperationMode determines whether info about the urban transit operation mode should be displayed for DetailedTimetable objects.
var DoShowOperationMode bool

// DoShowRoute determines whether info about the urban transit line route should be displayed for DetailedTimetable objects.
var DoShowRoute bool

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

func (line *Line) getTimetableStringDetails(operationModeRoutes *OperationModeRoutes, route *Route, stop *Stop) (timetableDetailsString string) {
	stopTitle := stop.String()
	timetableDetailsString += stopTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(stopTitle)) + "\n"
	if DoShowOperationMode {
		timetableDetailsString += "(" + operationModeRoutes.OperationMode.String() + ")\n"
	}
	timetableDetailsString += "* " + l10n.Translator[line.VehicleType] + " " + line.LineNumber
	if DoShowRoute {
		timetableDetailsString += " - " + l10n.Translator[l10n.OnRoute] + " " + route.Name + "(" + route.Code + ")"
	}
	timetableDetailsString += ": "
	return
}

// GetDetailedTimetableString fetches and returns a detailed string representation of the urban transit stop timetable matching the specified operationModeCode, routeCode and stopCode (annotated with information obtained from the Line object).
func (line *Line) GetDetailedTimetableString(operationModeCode string, routeCode string, stopCode string) (detailedTimetableString string, err error) {
	operationModeRoutes, ok := line.OperationModeRoutesMap[operationModeCode]
	if !ok {
		err = fmt.Errorf("could not find operation mode with code %s in info for line %s of type `%s`", operationModeCode, line.LineNumber, line.VehicleType)
		return
	}

	route, ok := operationModeRoutes.RouteMap[routeCode]
	if !ok {
		err = fmt.Errorf("could not find route with code %s for operation mode %s of line %s of type `%s`", routeCode, operationModeCode, line.LineNumber, line.VehicleType)
		return
	}

	stop, ok := route.StopMap[stopCode]
	if !ok {
		err = fmt.Errorf("could not find stop with code %s for route with code %s for operation mode %s of line %s of type `%s`", stopCode, routeCode, operationModeCode, line.LineNumber, line.VehicleType)
		return
	}

	timetable, err := GetTimetable(operationModeCode, routeCode, stopCode)
	if err != nil {
		return
	}

	detailedTimetableString += line.getTimetableStringDetails(operationModeRoutes, route, stop)
	detailedTimetableString += timetable.String() + "\n"
	return
}

// GetDetailedTimetableStrings fetches and returns a detailed string representation of the urban transit stop timetables for the specified line matching the specified operationModeCode, routeCode and stopCode.
func (line *Line) GetDetailedTimetableStrings(operationModeCode string, routeCode string, stopCode string) (detailedTimetableStrings []string, err error) {
	detailedTimetableStrings = []string{}
	for _, operationModeRoutes := range line.OperationModeRoutesList {
		if operationModeCode == "" || operationModeRoutes.Code == operationModeCode {
			for _, route := range operationModeRoutes.RouteList {
				if routeCode != "" || route.Code == routeCode {
					for _, stop := range route.StopList {
						if stopCode != "" || stop.Code == stopCode {
							timetable, err := GetTimetable(operationModeCode, routeCode, stopCode)
							if err != nil {
								return detailedTimetableStrings, err
							}

							detailedTimetableString := line.getTimetableStringDetails(operationModeRoutes, route, stop) + timetable.String() + "\n"
							detailedTimetableStrings = append(detailedTimetableStrings, detailedTimetableString)
						}
					}
				}
			}
		}
	}
	return
}
