package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/i18n"
	stcli_l10n "github.com/rgeorgiev583/sofiatraffic/stcli/l10n"
	"github.com/rgeorgiev583/sofiatraffic/virtual/l10n"

	"github.com/rgeorgiev583/sofiatraffic/virtual"
)

func main() {
	i18n.Init()
	l10n.InitTranslator()
	stcli_l10n.InitTranslator()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), stcli_l10n.Translator[stcli_l10n.Usage], os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var lineNumbersArg string
	flag.StringVar(&lineNumbersArg, stcli_l10n.Translator[stcli_l10n.LineNumbersFlagName], "", stcli_l10n.Translator[stcli_l10n.LineNumbersFlagUsage])

	var vehicleTypesArg string
	flag.StringVar(&vehicleTypesArg, stcli_l10n.Translator[stcli_l10n.VehicleTypesFlagName], "", fmt.Sprintf(stcli_l10n.Translator[stcli_l10n.VehicleTypesFlagUsage], l10n.Translator[l10n.VehicleTypeBus], l10n.Translator[l10n.VehicleTypeTrolleybus], l10n.Translator[l10n.VehicleTypeTram]))

	var stopCodesArg string
	flag.StringVar(&stopCodesArg, stcli_l10n.Translator[stcli_l10n.StopCodesFlagName], "", stcli_l10n.Translator[stcli_l10n.StopCodesFlagUsage])

	flag.BoolVar(&virtual.DoShowGenerationTimeForTimetables, stcli_l10n.Translator[stcli_l10n.DoShowGenerationTimeForTimetablesFlagName], false, stcli_l10n.Translator[stcli_l10n.DoShowGenerationTimeForTimetablesFlagUsage])

	flag.BoolVar(&virtual.DoShowFacilities, stcli_l10n.Translator[stcli_l10n.DoShowFacilitiesFlagName], false, fmt.Sprintf(stcli_l10n.Translator[stcli_l10n.DoShowFacilitiesFlagUsage], l10n.Translator[l10n.AirConditioningAbbreviation], l10n.Translator[l10n.WheelchairAccessibilityAbbreviation]))

	var doSortStops bool
	flag.BoolVar(&doSortStops, stcli_l10n.Translator[stcli_l10n.DoSortStopsFlagName], false, stcli_l10n.Translator[stcli_l10n.DoSortStopsFlagUsage])

	var doShowStops bool
	flag.BoolVar(&doShowStops, stcli_l10n.Translator[stcli_l10n.DoShowStopsFlagName], false, stcli_l10n.Translator[stcli_l10n.DoShowStopsFlagUsage])

	var doShowRoutes bool
	flag.BoolVar(&doShowRoutes, stcli_l10n.Translator[stcli_l10n.DoShowRoutesFlagName], false, stcli_l10n.Translator[stcli_l10n.DoShowRoutesFlagUsage])

	flag.BoolVar(&virtual.DoTranslateStopNames, stcli_l10n.Translator[stcli_l10n.DoTranslateStopNamesFlagName], false, stcli_l10n.Translator[stcli_l10n.DoTranslateStopNamesFlagUsage])

	flag.Parse()

	args := flag.Args()

	if doShowStops && (lineNumbersArg != "" || vehicleTypesArg != "" || stopCodesArg != "" || virtual.DoShowGenerationTimeForTimetables || virtual.DoShowFacilities) ||
		doShowRoutes && (stopCodesArg != "" || virtual.DoShowGenerationTimeForTimetables || virtual.DoShowFacilities) {
		fmt.Fprintln(os.Stderr, stcli_l10n.Translator[stcli_l10n.IncompatibleFlagsDetected])
		flag.Usage()
		os.Exit(1)
	}

	if doShowRoutes && lineNumbersArg == "" {
		fmt.Fprintln(os.Stderr, stcli_l10n.Translator[stcli_l10n.NoLineSpecified])
		flag.Usage()
		os.Exit(1)
	}

	lineNumbers := strings.Split(lineNumbersArg, ",")
	if lineNumbersArg != "" {
		for i, lineNumber := range lineNumbers {
			lineNumbers[i] = strings.TrimSpace(lineNumber)
		}
	}

	vehicleTypes := strings.Split(vehicleTypesArg, ",")
	if vehicleTypesArg != "" {
		for i, vehicleType := range vehicleTypes {
			vehicleTypes[i] = l10n.ReverseTranslator[strings.TrimSpace(vehicleType)]
		}
	}

	stopList, err := virtual.GetStops()
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
		routes, err := virtual.GetRoutes()
		if err != nil {
			log.Fatalln(err.Error())
		}

		stopMap := stopList.GetStopMap()

		forEachLine := func(f func(vehicleType string, lineNumber string)) {
			for _, vehicleType := range vehicleTypes {
				for _, lineNumber := range lineNumbers {
					f(vehicleType, lineNumber)
				}
			}
		}

		printRoutesByLine := func(vehicleType string, lineNumber string) {
			lineRoutes, err := routes.GetNamedRoutesByLine(vehicleType, lineNumber, stopMap)
			if err != nil {
				log.Println(err.Error())
				return
			}

			fmt.Print(lineRoutes)
		}

		forEachLine(printRoutesByLine)
		os.Exit(0)
	}

	forEachLine := func(stopCodeOrName string, f func(stopCodeOrName string, vehicleType string, lineNumber string)) {
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

	if stopCodesArg != "" {
		stopCodes := strings.Split(stopCodesArg, ",")
		for i, stopCode := range stopCodes {
			stopCodes[i] = strings.TrimSpace(stopCode)
			forEachLine(stopCodes[i], printTimetableByStopCodeAndLine)
		}
	}

	printTimetablesByStopNameAndLine := func(stopName string, vehicleType string, lineNumber string) {
		stopTimetables := stopList.GetTimetablesByStopNameAndLineAsync(stopName, vehicleTypesArg, lineNumbersArg, false)
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
