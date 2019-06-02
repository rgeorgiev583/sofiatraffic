package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/i18n"
	"github.com/rgeorgiev583/sofiatraffic/stcli/l10n"
	virtual_l10n "github.com/rgeorgiev583/sofiatraffic/virtual/l10n"

	"github.com/rgeorgiev583/sofiatraffic/virtual"
)

func main() {
	i18n.Init()
	virtual_l10n.InitTranslator()
	l10n.InitTranslator()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), l10n.Translator[l10n.Usage], os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var lineNumbersArg string
	flag.StringVar(&lineNumbersArg, l10n.Translator[l10n.LineNumbersFlagName], "", l10n.Translator[l10n.LineNumbersFlagUsage])

	var vehicleTypesArg string
	flag.StringVar(&vehicleTypesArg, l10n.Translator[l10n.VehicleTypesFlagName], "", fmt.Sprintf(l10n.Translator[l10n.VehicleTypesFlagUsage], virtual_l10n.Translator[virtual_l10n.VehicleTypeBus], virtual_l10n.Translator[virtual_l10n.VehicleTypeTrolleybus], virtual_l10n.Translator[virtual_l10n.VehicleTypeTram]))

	var stopCodesArg string
	flag.StringVar(&stopCodesArg, l10n.Translator[l10n.StopCodesFlagName], "", l10n.Translator[l10n.StopCodesFlagUsage])

	flag.BoolVar(&virtual.DoShowGenerationTimeForTimetables, l10n.Translator[l10n.DoShowGenerationTimeForTimetablesFlagName], false, l10n.Translator[l10n.DoShowGenerationTimeForTimetablesFlagUsage])

	flag.BoolVar(&virtual.DoShowFacilities, l10n.Translator[l10n.DoShowFacilitiesFlagName], false, fmt.Sprintf(l10n.Translator[l10n.DoShowFacilitiesFlagUsage], virtual_l10n.Translator[virtual_l10n.AirConditioningAbbreviation], virtual_l10n.Translator[virtual_l10n.WheelchairAccessibilityAbbreviation]))

	var doSortStops bool
	flag.BoolVar(&doSortStops, l10n.Translator[l10n.DoSortStopsFlagName], false, l10n.Translator[l10n.DoSortStopsFlagUsage])

	var doShowStops bool
	flag.BoolVar(&doShowStops, l10n.Translator[l10n.DoShowStopsFlagName], false, l10n.Translator[l10n.DoShowStopsFlagUsage])

	var doShowRoutes bool
	flag.BoolVar(&doShowRoutes, l10n.Translator[l10n.DoShowRoutesFlagName], false, l10n.Translator[l10n.DoShowRoutesFlagUsage])

	flag.BoolVar(&virtual.DoTranslateStopNames, l10n.Translator[l10n.DoTranslateStopNamesFlagName], false, l10n.Translator[l10n.DoTranslateStopNamesFlagUsage])

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

	lineNumbers := strings.Split(lineNumbersArg, ",")
	if lineNumbersArg != "" {
		for i, lineNumber := range lineNumbers {
			lineNumbers[i] = strings.TrimSpace(lineNumber)
		}
	}

	vehicleTypes := strings.Split(vehicleTypesArg, ",")
	if vehicleTypesArg != "" {
		for i, vehicleType := range vehicleTypes {
			vehicleTypes[i] = virtual_l10n.ReverseTranslator[strings.TrimSpace(vehicleType)]
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
