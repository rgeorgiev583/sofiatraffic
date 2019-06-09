package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/schedule"

	"github.com/rgeorgiev583/sofiatraffic/i18n"
	schedule_l10n "github.com/rgeorgiev583/sofiatraffic/schedule/l10n"
	"github.com/rgeorgiev583/sofiatraffic/stcli/l10n"
	virtual_l10n "github.com/rgeorgiev583/sofiatraffic/virtual/l10n"

	"github.com/rgeorgiev583/sofiatraffic/virtual"
)

// uniq returns a slice containing the sequentially unique elements of list (i.e. the ones not repeated in a row).
func uniq(list []string) (uniqueItems []string) {
	uniqueItems = list[:0]
	for _, item := range list {
		if len(uniqueItems) == 0 || uniqueItems[len(uniqueItems)-1] != item {
			uniqueItems = append(uniqueItems, item)
		}
	}
	return
}

func main() {
	i18n.Init()
	l10n.InitTranslator()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), l10n.Translator[l10n.Usage], os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var lineNumbersArg string
	flag.StringVar(&lineNumbersArg, l10n.Translator[l10n.LineNumbersFlagName], "", l10n.Translator[l10n.LineNumbersFlagUsage])

	var vehicleTypesArg string
	flag.StringVar(&vehicleTypesArg, l10n.Translator[l10n.VehicleTypesFlagName], "", fmt.Sprintf(l10n.Translator[l10n.VehicleTypesFlagUsage], l10n.Translator[l10n.VehicleTypeBus], l10n.Translator[l10n.VehicleTypeTrolleybus], l10n.Translator[l10n.VehicleTypeTram]))

	var stopCodesArg string
	flag.StringVar(&stopCodesArg, l10n.Translator[l10n.StopCodesFlagName], "", l10n.Translator[l10n.StopCodesFlagUsage])

	var routeCodesArg string
	flag.StringVar(&routeCodesArg, l10n.Translator[l10n.RouteCodesFlagName], "", l10n.Translator[l10n.RouteCodesFlagUsage])

	var operationModeCodesArg string
	flag.StringVar(&operationModeCodesArg, l10n.Translator[l10n.OperationModeCodesFlagName], "", l10n.Translator[l10n.OperationModeCodesFlagUsage])

	flag.BoolVar(&virtual.DoShowGenerationTimeForTimetables, l10n.Translator[l10n.DoShowGenerationTimeForTimetablesFlagName], false, l10n.Translator[l10n.DoShowGenerationTimeForTimetablesFlagUsage])

	flag.BoolVar(&virtual.DoShowFacilities, l10n.Translator[l10n.DoShowFacilitiesFlagName], false, fmt.Sprintf(l10n.Translator[l10n.DoShowFacilitiesFlagUsage], l10n.Translator[l10n.AirConditioningAbbreviation], l10n.Translator[l10n.WheelchairAccessibilityAbbreviation]))

	var doSortStops bool
	flag.BoolVar(&doSortStops, l10n.Translator[l10n.DoSortStopsFlagName], false, l10n.Translator[l10n.DoSortStopsFlagUsage])

	var doShowStops bool
	flag.BoolVar(&doShowStops, l10n.Translator[l10n.DoShowStopsFlagName], false, l10n.Translator[l10n.DoShowStopsFlagUsage])

	var doShowLines bool
	flag.BoolVar(&doShowLines, l10n.Translator[l10n.DoShowLinesFlagName], false, l10n.Translator[l10n.DoShowLinesFlagUsage])

	var doShowRoutes bool
	flag.BoolVar(&doShowRoutes, l10n.Translator[l10n.DoShowRoutesFlagName], false, l10n.Translator[l10n.DoShowRoutesFlagUsage])

	flag.BoolVar(&virtual.DoTranslateStopNames, l10n.Translator[l10n.DoTranslateStopNamesFlagName], false, l10n.Translator[l10n.DoTranslateStopNamesFlagUsage])

	var doUseSchedule bool
	flag.BoolVar(&doUseSchedule, l10n.Translator[l10n.DoUseScheduleFlagName], false, l10n.Translator[l10n.DoUseScheduleFlagUsage])

	flag.Parse()
	args := flag.Args()

	if doShowStops && (lineNumbersArg != "" || vehicleTypesArg != "" || stopCodesArg != "" || virtual.DoShowGenerationTimeForTimetables || virtual.DoShowFacilities) ||
		doShowRoutes && (stopCodesArg != "" || virtual.DoShowGenerationTimeForTimetables || virtual.DoShowFacilities) {
		fmt.Fprintln(os.Stderr, l10n.Translator[l10n.IncompatibleFlagsDetected])
		flag.Usage()
		os.Exit(1)
	}

	if doShowRoutes && lineNumbersArg == "" {
		fmt.Fprintln(os.Stderr, l10n.Translator[l10n.NoLineSpecified])
		flag.Usage()
		os.Exit(1)
	}

	var libraryReverseTranslator map[string]string
	if doUseSchedule {
		schedule_l10n.InitTranslator()
		libraryReverseTranslator = schedule_l10n.ReverseTranslator
	} else {
		virtual_l10n.InitTranslator()
		libraryReverseTranslator = virtual_l10n.ReverseTranslator
	}

	lineNumbers := strings.Split(lineNumbersArg, ",")
	lineNumbers = uniq(lineNumbers)
	for i, lineNumber := range lineNumbers {
		lineNumbers[i] = strings.TrimSpace(lineNumber)
	}

	vehicleTypes := strings.Split(vehicleTypesArg, ",")
	vehicleTypes = uniq(vehicleTypes)
	for i, vehicleType := range vehicleTypes {
		vehicleTypes[i] = libraryReverseTranslator[strings.TrimSpace(vehicleType)]
	}

	stopCodes := strings.Split(stopCodesArg, ",")
	stopCodes = uniq(stopCodes)
	for i, stopCode := range stopCodes {
		stopCodes[i] = strings.TrimSpace(stopCode)
	}

	routeCodes := strings.Split(routeCodesArg, ",")
	routeCodes = uniq(routeCodes)
	for i, routeCode := range routeCodes {
		routeCodes[i] = strings.TrimSpace(routeCode)
	}

	operationModeCodes := strings.Split(operationModeCodesArg, ",")
	operationModeCodes = uniq(operationModeCodes)
	for i, operationModeCode := range operationModeCodes {
		operationModeCodes[i] = strings.TrimSpace(operationModeCode)
	}

	forEachLine := func(f func(vehicleType string, lineNumber string)) {
		for _, vehicleType := range vehicleTypes {
			for _, lineNumber := range lineNumbers {
				f(vehicleType, lineNumber)
			}
		}
	}

	if doUseSchedule {
		if doShowLines {
			lines, err := schedule.GetLines()
			if err != nil {
				log.Fatalln(err.Error())
			}

			fmt.Println(lines)
		} else if doShowRoutes {
			printRoutesByLine := func(vehicleType string, lineNumber string) {
				lineRoutes, err := schedule.GetLine(vehicleType, lineNumber)
				if err != nil {
					log.Println(err.Error())
					return
				}

				fmt.Print(lineRoutes)
			}
			forEachLine(printRoutesByLine)
		} else {
			forEachRouteByStop := func(stopCode string, f func(stopCode string, operationModeCode string, routeCode string)) {
				for _, operationModeCode := range operationModeCodes {
					for _, routeCode := range routeCodes {
						f(stopCode, operationModeCode, routeCode)
					}
				}
			}
			printTimetableByStopCodeAndRoute := func(stopCode string, operationModeCode string, routeCode string) {
				stopTimetable, err := schedule.GetTimetable(operationModeCode, routeCode, stopCode)
				if err != nil {
					log.Println(err.Error())
					return
				}

				fmt.Print(stopTimetable)
			}
			for _, stopCode := range stopCodes {
				forEachRouteByStop(stopCode, printTimetableByStopCodeAndRoute)
			}
		}
	} else {
		stopList, err := virtual.GetStops()
		if err != nil {
			log.Fatalln(err.Error())
		}

		if doSortStops {
			sort.Sort(stopList)
		}

		if doShowStops {
			fmt.Print(stopList)
		} else if doShowRoutes {
			routes, err := virtual.GetRoutes()
			if err != nil {
				log.Fatalln(err.Error())
			}

			stopMap := stopList.GetStopMap()
			if len(vehicleTypes) == 1 && len(lineNumbers) == 1 {
				lineRoutes, err := routes.GetNamedRoutesByLine(vehicleTypes[0], lineNumbers[0], stopMap)
				if err != nil {
					log.Fatalln(err.Error())
				}

				fmt.Print(lineRoutes)
			} else {
				routeMap, err := routes.GetRouteMap(stopMap)
				if err != nil {
					log.Fatalln(err.Error())
				}

				printRoutesByLine := func(vehicleType string, lineNumber string) {
					lineRoutes, ok := routeMap[virtual.Line{VehicleType: vehicleType, LineNumber: lineNumber}]
					if !ok {
						log.Fatalf("could not find line with vehicle type %s and number %s in the route map\n", vehicleType, lineNumber)
						return
					}

					fmt.Print(lineRoutes)
				}
				forEachLine(printRoutesByLine)
			}
		} else {
			forEachLineByStop := func(stopCodeOrName string, f func(stopCodeOrName string, vehicleType string, lineNumber string)) {
				for _, vehicleType := range vehicleTypes {
					for _, lineNumber := range lineNumbers {
						f(stopCodeOrName, vehicleType, lineNumber)
					}
				}
			}
			printTimetableByStopCodeAndLine := func(stopCode string, vehicleType string, lineNumber string) {
				stopTimetable, err := virtual.GetTimetableByStopCodeAndLine(stopCode, vehicleTypesArg, lineNumbersArg)
				if err != nil {
					log.Println(err.Error())
					return
				}

				fmt.Print(stopTimetable)
			}
			for _, stopCode := range stopCodes {
				forEachLineByStop(stopCode, printTimetableByStopCodeAndLine)
			}
			printTimetablesByStopNameAndLine := func(stopName string, vehicleType string, lineNumber string) {
				stopTimetables := stopList.GetTimetablesByStopNameAndLineAsync(stopName, vehicleTypesArg, lineNumbersArg, false)
				fmt.Print(stopTimetables)
			}
			if len(args) > 0 {
				for _, stopName := range args {
					forEachLineByStop(stopName, printTimetablesByStopNameAndLine)
				}
			} else {
				forEachLineByStop("", printTimetablesByStopNameAndLine)
			}
		}
	}
}
