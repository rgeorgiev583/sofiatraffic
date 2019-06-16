package schedule

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rgeorgiev583/sofiatraffic/i18n"
	"github.com/rgeorgiev583/sofiatraffic/schedule/l10n"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// OperationMode represents a type of urban transit line schedule classified by the frequency of vehicle arrivals (which is determined by the current day being a workday or a holiday).
type OperationMode struct {
	Code, Name string
}

// Stop represents an urban transit stop.
type Stop struct {
	Code, Name string
}

// StopList represents a list of urban transit stops.
type StopList []*Stop

// StopMap represents a map from the code of each urban transit stop to its corresponding Stop object.
type StopMap map[string]*Stop

// Route represents the sequence of stops where an urban transit line stops when traveling in a specific direction.
type Route struct {
	Code, Name string
	StopList
	StopMap
}

// RouteList represents a list of routes.
type RouteList []*Route

// RouteMap represents a map from the code of each urban transit line route to its corresponding Route object.
type RouteMap map[string]*Route

// OperationModeRoutes represents the routes for a specific urban transit line operation mode.
type OperationModeRoutes struct {
	*OperationMode
	RouteList
	RouteMap
}

// OperationModeRoutesList represents a list of OperationModeRoutes objects.
type OperationModeRoutesList []*OperationModeRoutes

// OperationModeRoutesMap represents a map from the code of each urban transit line operation mode to its corresponding OperationModeRoutes object.
type OperationModeRoutesMap map[string]*OperationModeRoutes

// Line represents an urban transit line.
type Line struct {
	VehicleType, LineNumber string
	OperationModeRoutesList
	OperationModeRoutesMap
}

type lineScannerState int

const (
	// VehicleTypeBus represents a bus.
	VehicleTypeBus = "autobus"
	// VehicleTypeTrolleybus represents a trolleybus.
	VehicleTypeTrolleybus = "trolleybus"
	// VehicleTypeTram represents a tram.
	VehicleTypeTram = "tramway"
	// VehicleTypeMetro represents the metro.
	VehicleTypeMetro = "metro"

	lineScannerNotInsideRelevantElement lineScannerState = iota
	lineScannerInsideDirectionAnchor
	lineScannerInsideDirectionSpan
	lineScannerInsideOperationAnchor
	lineScannerInsideOperationSpan
	lineScannerInsideStopAnchor
)

// DoTranslateStopNames determines whether stop names should be translated from Bulgarian to the local language.
var DoTranslateStopNames bool

// StopNameTranslator maps names of stops in Bulgarian to their translation in the local language.
var StopNameTranslator map[string]string

func translateOperationModeName(name string) string {
	name = strings.ReplaceAll(name, l10n.BulgarianTranslator[l10n.OperationModeWeekday], l10n.Translator[l10n.OperationModeWeekday])
	name = strings.ReplaceAll(name, l10n.BulgarianTranslator[l10n.OperationModePreHoliday], l10n.Translator[l10n.OperationModePreHoliday])
	name = strings.ReplaceAll(name, l10n.BulgarianTranslator[l10n.OperationModeHoliday], l10n.Translator[l10n.OperationModeHoliday])
	return name
}

// GetLine returns the urban transit line with the specified vehicleType and lineNumber.
func GetLine(vehicleType string, lineNumber string) (line *Line, err error) {
	linePageURL := &url.URL{
		Scheme: Scheme,
		Host:   Hostname,
		Path:   "/" + vehicleType + "/" + lineNumber,
	}
	response, err := http.Get(linePageURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the schedule line page: %s", err.Error())
		return
	}
	defer response.Body.Close()

	line = &Line{
		VehicleType:            vehicleType,
		LineNumber:             lineNumber,
		OperationModeRoutesMap: map[string]*OperationModeRoutes{},
	}
	var schedule *OperationModeRoutes
	var route *Route
	var stop *Stop
	var state lineScannerState
	for tokenizer := html.NewTokenizer(response.Body); tokenizer.Next() != html.ErrorToken; {
		token := tokenizer.Token()
		switch token.Type {
		case html.StartTagToken:
			switch token.DataAtom {
			case atom.A:
				var id, class string
				for _, attr := range token.Attr {
					if id != "" && class != "" {
						break
					}

					switch atom.Lookup([]byte(attr.Key)) {
					case atom.Id:
						id = attr.Val

					case atom.Class:
						class = attr.Val
					}
				}
				if id == "" || class == "" {
					continue
				}

				if strings.Contains(class, "schedule_active_list_tab") {
					state = lineScannerInsideOperationAnchor
					idComponents := strings.Split(id, "_")
					if len(idComponents) != 3 || idComponents[0] != "schedule" || idComponents[2] != "button" {
						err = fmt.Errorf("the value of the `id` attribute of the `a` element is not valid: %s", id)
						return line, err
					}

					operationModeCode := idComponents[1]
					schedule = &OperationModeRoutes{
						OperationMode: &OperationMode{Code: operationModeCode},
						RouteList:     RouteList{},
						RouteMap:      RouteMap{},
					}
					line.OperationModeRoutesList = append(line.OperationModeRoutesList, schedule)
					line.OperationModeRoutesMap[operationModeCode] = schedule
				} else if strings.Contains(class, "schedule_view_direction_tab") {
					state = lineScannerInsideDirectionAnchor
					idComponents := strings.Split(id, "_")
					if len(idComponents) != 5 || idComponents[0] != "schedule" || idComponents[1] != "direction" || idComponents[4] != "button" {
						err = fmt.Errorf("the value of the `id` attribute of the `a` element is not valid: %s", id)
						return line, err
					}

					operationModeCode := idComponents[2]
					schedule, ok := line.OperationModeRoutesMap[operationModeCode]
					if !ok {
						err = fmt.Errorf("invalid operation mode code: %s", operationModeCode)
						return line, err
					}

					routeCode := idComponents[3]
					route = &Route{
						Code:     routeCode,
						StopList: StopList{},
						StopMap:  StopMap{},
					}
					schedule.RouteList = append(schedule.RouteList, route)
					schedule.RouteMap[routeCode] = route
				} else if strings.Contains(class, "stop_change") {
					state = lineScannerInsideStopAnchor
					idComponents := strings.Split(id, "_")
					if len(idComponents) != 6 || idComponents[0] != "schedule" || idComponents[2] != "direction" || idComponents[4] != "sign" {
						err = fmt.Errorf("the value of the `id` attribute of the `a` element is not valid: %s", id)
						return line, err
					}

					operationModeCode := idComponents[1]
					schedule, ok := line.OperationModeRoutesMap[operationModeCode]
					if !ok {
						err = fmt.Errorf("invalid operation mode code: %s", operationModeCode)
						return line, err
					}

					routeCode := idComponents[3]
					route, ok := schedule.RouteMap[routeCode]
					if !ok {
						err = fmt.Errorf("invalid route code: %s", operationModeCode)
						return line, err
					}

					stopCode := idComponents[5]
					stop = &Stop{Code: stopCode}
					route.StopList = append(route.StopList, stop)
					route.StopMap[stopCode] = stop
				}

			case atom.Span:
				switch state {
				case lineScannerInsideOperationAnchor:
					state = lineScannerInsideOperationSpan

				case lineScannerInsideDirectionAnchor:
					state = lineScannerInsideDirectionSpan
				}
			}

		case html.EndTagToken:
			switch token.DataAtom {
			case atom.A, atom.Span:
				state = lineScannerNotInsideRelevantElement
			}

		case html.TextToken:
			switch state {
			case lineScannerInsideOperationSpan:
				schedule.OperationMode.Name = token.Data

			case lineScannerInsideDirectionSpan:
				route.Name = token.Data

			case lineScannerInsideStopAnchor:
				stop.Name = token.Data
			}
		}
	}
	return
}

func (om *OperationMode) String() string {
	return l10n.Translator[l10n.OperationMode] + ": " + translateOperationModeName(om.Name) + " (" + om.Code + ")"
}

func (s *Stop) String() string {
	var translatedStopName string
	if DoTranslateStopNames && i18n.Language == i18n.LanguageCodeEnglish {
		translatedStopName = StopNameTranslator[s.Name]
	} else {
		translatedStopName = s.Name
	}
	return translatedStopName + " (" + s.Code + ")"
}

func (sl StopList) String() string {
	var builder strings.Builder
	for i, stop := range sl {
		builder.WriteString(strconv.Itoa(i+1) + ". " + stop.String() + "\n")
	}
	return builder.String()
}

func (r *Route) String() string {
	return "### " + r.Name + " (" + r.Code + ")\n" + r.StopList.String()
}

func (nrl RouteList) String() string {
	var builder strings.Builder
	for _, route := range nrl {
		builder.WriteString(route.String() + "\n")
	}
	return builder.String()
}

func (omr *OperationModeRoutes) String() string {
	operationModeTitle := omr.OperationMode.String()
	return operationModeTitle + "\n" + strings.Repeat("-", utf8.RuneCountInString(operationModeTitle)) + "\n" + omr.RouteList.String()
}

func (omrl OperationModeRoutesList) String() string {
	var builder strings.Builder
	for _, operationModeRoutes := range omrl {
		builder.WriteString(operationModeRoutes.String() + "\n")
	}
	return builder.String()
}

func (l *Line) String() string {
	lineTitle := l10n.Translator[l.VehicleType] + " " + l.LineNumber
	return lineTitle + "\n" + strings.Repeat("=", utf8.RuneCountInString(lineTitle)) + "\n\n" + l.OperationModeRoutesList.String()
}
