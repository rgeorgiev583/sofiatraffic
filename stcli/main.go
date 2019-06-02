package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/l10n"
	stcli_l10n "github.com/rgeorgiev583/sofiatraffic/stcli/l10n"

	"github.com/rgeorgiev583/sofiatraffic/regular"
)

func main() {
	l10n.InitTranslator()
	stcli_l10n.InitTranslator()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), stcli_l10n.Translator[stcli_l10n.Usage], os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var lineCodesArg string
	flag.StringVar(&lineCodesArg, stcli_l10n.Translator[stcli_l10n.LineCodesFlagName], "", stcli_l10n.Translator[stcli_l10n.LineCodesFlagUsage])

	var vehicleTypesArg string
	flag.StringVar(&vehicleTypesArg, stcli_l10n.Translator[stcli_l10n.VehicleTypesFlagName], "", fmt.Sprintf(stcli_l10n.Translator[stcli_l10n.VehicleTypesFlagUsage], l10n.Translator[l10n.VehicleTypeBus], l10n.Translator[l10n.VehicleTypeTrolleybus], l10n.Translator[l10n.VehicleTypeTram]))

	var stopCodesArg string
	flag.StringVar(&stopCodesArg, stcli_l10n.Translator[stcli_l10n.StopCodesFlagName], "", stcli_l10n.Translator[stcli_l10n.StopCodesFlagUsage])

	flag.BoolVar(&regular.DoShowGenerationTimeForTimetables, stcli_l10n.Translator[stcli_l10n.DoShowGenerationTimeForTimetablesFlagName], false, stcli_l10n.Translator[stcli_l10n.DoShowGenerationTimeForTimetablesFlagUsage])

	flag.BoolVar(&regular.DoShowFacilities, stcli_l10n.Translator[stcli_l10n.DoShowFacilitiesFlagName], false, fmt.Sprintf(stcli_l10n.Translator[stcli_l10n.DoShowFacilitiesFlagUsage], l10n.Translator[l10n.AirConditioningAbbreviation], l10n.Translator[l10n.WheelchairAccessibilityAbbreviation]))

	var doSortStops bool
	flag.BoolVar(&doSortStops, stcli_l10n.Translator[stcli_l10n.DoSortStopsFlagName], false, stcli_l10n.Translator[stcli_l10n.DoSortStopsFlagUsage])

	var doShowStops bool
	flag.BoolVar(&doShowStops, stcli_l10n.Translator[stcli_l10n.DoShowStopsFlagName], false, stcli_l10n.Translator[stcli_l10n.DoShowStopsFlagUsage])

	var doShowRoutes bool
	flag.BoolVar(&doShowRoutes, stcli_l10n.Translator[stcli_l10n.DoShowRoutesFlagName], false, stcli_l10n.Translator[stcli_l10n.DoShowRoutesFlagUsage])

	flag.BoolVar(&regular.DoTranslateStopNames, stcli_l10n.Translator[stcli_l10n.DoTranslateStopNamesFlagName], false, stcli_l10n.Translator[stcli_l10n.DoTranslateStopNamesFlagUsage])

	flag.Parse()

	args := flag.Args()

	if doShowStops && (lineCodesArg != "" || vehicleTypesArg != "" || stopCodesArg != "" || regular.DoShowGenerationTimeForTimetables || regular.DoShowFacilities) ||
		doShowRoutes && (stopCodesArg != "" || regular.DoShowGenerationTimeForTimetables || regular.DoShowFacilities) {
		fmt.Fprintln(os.Stderr, stcli_l10n.Translator[stcli_l10n.IncompatibleFlagsDetected])
		flag.Usage()
		os.Exit(1)
	}

	lineCodes := strings.Split(lineCodesArg, ",")
	if lineCodesArg != "" {
		for i, lineCode := range lineCodes {
			lineCodes[i] = strings.TrimSpace(lineCode)
		}
	}

	vehicleTypes := strings.Split(vehicleTypesArg, ",")
	if vehicleTypesArg != "" {
		for i, vehicleType := range vehicleTypes {
			vehicleTypes[i] = l10n.ReverseTranslator[strings.TrimSpace(vehicleType)]
		}
	}

	stopList, err := regular.GetStops()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if doSortStops {
		sort.Sort(stopList)
	}

	if doShowStops {
		fmt.Print(stopList)
		os.Exit(0)
	}

	if doShowRoutes {
		routes, err := regular.GetRoutes()
		if err != nil {
			log.Fatalln(err.Error())
		}

		stopMap := stopList.GetStopMap()

		forEachLine := func(f func(vehicleType string, lineCode string)) {
			for _, vehicleType := range vehicleTypes {
				for _, lineCode := range lineCodes {
					f(vehicleType, lineCode)
				}
			}
		}

		printRoutesByLine := func(vehicleType string, lineCode string) {
			lineRoutes, err := routes.GetNamedRoutesByLine(vehicleType, lineCode, stopMap)
			if err != nil {
				log.Println(err.Error())
				return
			}

			fmt.Print(lineRoutes)
		}

		forEachLine(printRoutesByLine)
		os.Exit(0)
	}

	forEachLine := func(stopCodeOrName string, f func(stopCodeOrName string, vehicleType string, lineCode string)) {
		for _, vehicleType := range vehicleTypes {
			for _, lineCode := range lineCodes {
				f(stopCodeOrName, vehicleType, lineCode)
			}
		}
	}

	printTimetableByStopCodeAndLine := func(stopCode string, vehicleType string, lineCode string) {
		stopTimetable, err := regular.GetTimetableByStopCodeAndLine(stopCode, vehicleTypesArg, lineCodesArg)
		if err != nil {
			log.Println(err.Error())
			return
		}

		fmt.Print(stopTimetable)
	}

	if stopCodesArg != "" {
		stopCodes := strings.Split(stopCodesArg, ",")
		for i, stopCode := range stopCodes {
			stopCodes[i] = strings.TrimSpace(stopCode)
			forEachLine(stopCodes[i], printTimetableByStopCodeAndLine)
		}
	}

	printTimetablesByStopNameAndLine := func(stopName string, vehicleType string, lineCode string) {
		stopTimetables := stopList.GetTimetablesByStopNameAndLineAsync(stopName, vehicleTypesArg, lineCodesArg, false)
		fmt.Print(stopTimetables)
	}

	if len(args) > 0 {
		for _, stopName := range args {
			forEachLine(stopName, printTimetablesByStopNameAndLine)
		}
	} else {
		forEachLine("", printTimetablesByStopNameAndLine)
	}
}
