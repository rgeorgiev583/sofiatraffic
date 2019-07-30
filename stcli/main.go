/*
Package stcli implements a basic tool with a command-line interface for fetching of information about public transit in Sofia.
*/
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

type commandMode int

const (
	timetablesMode commandMode = iota
	stopsMode
	linesMode
	routesMode
)

type commandContext struct {
	command                                                                                                                   *flag.FlagSet
	lineNumbersArg, vehicleTypesArg, stopCodesArg, routeCodesArg, routeNamesArg, operationModeCodesArg, operationModeNamesArg string
	doSortStops, doTranslateStopNames, doUseSchedule                                                                          bool
	positionalArgs                                                                                                            []string
}

func initCommandContextInMode(mode commandMode, args []string) (context *commandContext, err error) {
	context = &commandContext{}
	switch mode {
	case timetablesMode:
		context.command = flag.NewFlagSet("timetables", flag.ExitOnError)
		context.command.Usage = func() {
			fmt.Fprintf(context.command.Output(), l10n.Translator[l10n.TimetablesSubcommandUsage], os.Args[0])
			context.command.PrintDefaults()
		}
		context.command.StringVar(&context.lineNumbersArg, l10n.Translator[l10n.LineNumbersFlagName], "", l10n.Translator[l10n.LineNumbersFlagUsage])
		context.command.StringVar(&context.vehicleTypesArg, l10n.Translator[l10n.VehicleTypesFlagName], "", fmt.Sprintf(l10n.Translator[l10n.VehicleTypesFlagUsage], l10n.Translator[l10n.VehicleTypeBus], l10n.Translator[l10n.VehicleTypeTrolleybus], l10n.Translator[l10n.VehicleTypeTram]))
		context.command.StringVar(&context.stopCodesArg, l10n.Translator[l10n.StopCodesFlagName], "", l10n.Translator[l10n.StopCodesFlagUsage])
		context.command.StringVar(&context.routeCodesArg, l10n.Translator[l10n.RouteCodesFlagName], "", l10n.Translator[l10n.RouteCodesFlagUsage])
		context.command.StringVar(&context.routeNamesArg, l10n.Translator[l10n.RouteNamesFlagName], "", l10n.Translator[l10n.RouteNamesFlagUsage])
		context.command.StringVar(&context.operationModeCodesArg, l10n.Translator[l10n.OperationModeCodesFlagName], "", l10n.Translator[l10n.OperationModeCodesFlagUsage])
		context.command.StringVar(&context.operationModeNamesArg, l10n.Translator[l10n.OperationModeNamesFlagName], "", l10n.Translator[l10n.OperationModeNamesFlagUsage])
		context.command.BoolVar(&virtual.DoShowGenerationTimeForTimetables, l10n.Translator[l10n.DoShowGenerationTimeForTimetablesFlagName], false, l10n.Translator[l10n.DoShowGenerationTimeForTimetablesFlagUsage])
		context.command.BoolVar(&virtual.DoShowRemainingTimeUntilArrival, l10n.Translator[l10n.DoShowRemainingTimeUntilArrivalFlagName], false, l10n.Translator[l10n.DoShowRemainingTimeUntilArrivalFlagUsage])
		context.command.BoolVar(&virtual.DoShowFacilities, l10n.Translator[l10n.DoShowFacilitiesFlagName], false, fmt.Sprintf(l10n.Translator[l10n.DoShowFacilitiesFlagUsage], l10n.Translator[l10n.AirConditioningAbbreviation], l10n.Translator[l10n.WheelchairAccessibilityAbbreviation]))
		context.command.BoolVar(&schedule.DoShowOperationMode, l10n.Translator[l10n.DoShowOperationModeFlagName], false, l10n.Translator[l10n.DoShowOperationModeFlagUsage])
		context.command.BoolVar(&schedule.DoShowRoute, l10n.Translator[l10n.DoShowRouteFlagName], false, l10n.Translator[l10n.DoShowRouteFlagUsage])
		context.command.BoolVar(&context.doSortStops, l10n.Translator[l10n.DoSortStopsFlagName], false, l10n.Translator[l10n.DoSortStopsFlagUsage])
		context.command.BoolVar(&context.doTranslateStopNames, l10n.Translator[l10n.DoTranslateStopNamesFlagName], false, l10n.Translator[l10n.DoTranslateStopNamesFlagUsage])
		context.command.BoolVar(&context.doUseSchedule, l10n.Translator[l10n.DoUseScheduleFlagName], false, l10n.Translator[l10n.DoUseScheduleFlagUsage])

	case stopsMode:
		context.command = flag.NewFlagSet("stops", flag.ExitOnError)
		context.command.Usage = func() {
			fmt.Fprintf(context.command.Output(), l10n.Translator[l10n.StopsSubcommandUsage], os.Args[0])
			context.command.PrintDefaults()
		}
		context.command.BoolVar(&context.doSortStops, l10n.Translator[l10n.DoSortStopsFlagName], false, l10n.Translator[l10n.DoSortStopsFlagUsage])
		context.command.BoolVar(&context.doTranslateStopNames, l10n.Translator[l10n.DoTranslateStopNamesFlagName], false, l10n.Translator[l10n.DoTranslateStopNamesFlagUsage])

	case linesMode:
		context.command = flag.NewFlagSet("lines", flag.ExitOnError)
		context.command.Usage = func() {
			fmt.Fprintf(context.command.Output(), l10n.Translator[l10n.LinesSubcommandUsage], os.Args[0])
			context.command.PrintDefaults()
		}
		context.doUseSchedule = true

	case routesMode:
		context.command = flag.NewFlagSet("routes", flag.ExitOnError)
		context.command.Usage = func() {
			fmt.Fprintf(context.command.Output(), l10n.Translator[l10n.RoutesSubcommandUsage], os.Args[0])
			context.command.PrintDefaults()
		}
		context.command.StringVar(&context.lineNumbersArg, l10n.Translator[l10n.LineNumbersFlagName], "", l10n.Translator[l10n.LineNumbersFlagUsage])
		context.command.StringVar(&context.vehicleTypesArg, l10n.Translator[l10n.VehicleTypesFlagName], "", fmt.Sprintf(l10n.Translator[l10n.VehicleTypesFlagUsage], l10n.Translator[l10n.VehicleTypeBus], l10n.Translator[l10n.VehicleTypeTrolleybus], l10n.Translator[l10n.VehicleTypeTram]))
		context.command.BoolVar(&context.doSortStops, l10n.Translator[l10n.DoSortStopsFlagName], false, l10n.Translator[l10n.DoSortStopsFlagUsage])
		context.command.BoolVar(&context.doTranslateStopNames, l10n.Translator[l10n.DoTranslateStopNamesFlagName], false, l10n.Translator[l10n.DoTranslateStopNamesFlagUsage])
		context.command.BoolVar(&context.doUseSchedule, l10n.Translator[l10n.DoUseScheduleFlagName], false, l10n.Translator[l10n.DoUseScheduleFlagUsage])
	}

	err = context.command.Parse(args)
	if err != nil {
		return
	}

	virtual.DoTranslateStopNames = context.doTranslateStopNames
	schedule.DoTranslateStopNames = context.doTranslateStopNames
	context.positionalArgs = context.command.Args()
	return
}

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

func parseList(listString string) (list []string) {
	list = strings.Split(listString, ",")
	list = uniq(list)
	for i, element := range list {
		list[i] = strings.TrimSpace(element)
	}
	return
}

func initStopNameTranslatorIfNecessary() {
	if schedule.DoTranslateStopNames && i18n.Language == i18n.LanguageCodeEnglish {
		stopsInBulgarian, err := virtual.GetStopsInLanguage(i18n.LanguageCodeBulgarian)
		if err != nil {
			log.Fatalln(err.Error())
		}

		stopsInEnglish, err := virtual.GetStopsInLanguage(i18n.LanguageCodeEnglish)
		if err != nil {
			log.Fatalln(err.Error())
		}

		schedule.StopNameTranslator = map[string]string{}
		for i := 0; i < len(stopsInBulgarian) && i < len(stopsInEnglish); i++ {
			schedule.StopNameTranslator[stopsInBulgarian[i].Name] = stopsInEnglish[i].Name
		}
	}
}

func main() {
	i18n.Init()
	l10n.InitTranslator()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), l10n.Translator[l10n.Usage], os.Args[0], os.Args[0], os.Args[0])
	}

	var mode commandMode
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case l10n.Translator[l10n.TimetablesSubcommandName]:
			mode = timetablesMode

		case l10n.Translator[l10n.StopsSubcommandName]:
			mode = stopsMode

		case l10n.Translator[l10n.LinesSubcommandName]:
			mode = linesMode

		case l10n.Translator[l10n.RoutesSubcommandName]:
			mode = routesMode

		default:
			flag.Parse()

			fmt.Fprintln(os.Stderr, l10n.Translator[l10n.InvalidSubcommandName])
			flag.Usage()
			os.Exit(1)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
	context, err := initCommandContextInMode(mode, os.Args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if context.doUseSchedule && (virtual.DoShowGenerationTimeForTimetables || virtual.DoShowFacilities || context.doSortStops) {
		fmt.Fprintln(os.Stderr, l10n.Translator[l10n.IncompatibleFlagsDetected])
		context.command.Usage()
		os.Exit(1)
	}

	var libraryReverseTranslator map[string]string
	if context.doUseSchedule {
		schedule_l10n.InitTranslator()
		libraryReverseTranslator = schedule_l10n.ReverseTranslator
	} else {
		virtual_l10n.InitTranslator()
		libraryReverseTranslator = virtual_l10n.ReverseTranslator
	}

	lineNumbers := parseList(context.lineNumbersArg)
	vehicleTypes := parseList(context.vehicleTypesArg)
	stopCodes := parseList(context.stopCodesArg)
	routeCodes := parseList(context.routeCodesArg)
	routeNames := parseList(context.routeNamesArg)
	operationModeCodes := parseList(context.operationModeCodesArg)
	operationModeNames := parseList(context.operationModeNamesArg)

	stopNames := context.positionalArgs

	for i, vehicleType := range vehicleTypes {
		vehicleTypes[i] = libraryReverseTranslator[vehicleType]
	}

	forEachLine := func(f func(vehicleType string, lineNumber string)) {
		for _, vehicleType := range vehicleTypes {
			for _, lineNumber := range lineNumbers {
				f(vehicleType, lineNumber)
			}
		}
	}

	if context.doUseSchedule {
		switch mode {
		case linesMode:
			lines, err := schedule.GetLines()
			if err != nil {
				log.Fatalln(err.Error())
			}

			fmt.Println(lines)

		case routesMode:
			noVehicleTypesAreSpecified := len(vehicleTypes) == 1 && vehicleTypes[0] == ""
			noLineNumbersAreSpecified := len(lineNumbers) == 1 && lineNumbers[0] == ""
			if noVehicleTypesAreSpecified || noLineNumbersAreSpecified {
				detailsList := []string{}
				if noVehicleTypesAreSpecified {
					detailsList = append(detailsList, l10n.Translator[l10n.VehicleTypes])
				}
				if noLineNumbersAreSpecified {
					detailsList = append(detailsList, l10n.Translator[l10n.LineNumbers])
				}
				log.Fatalln(l10n.Translator[l10n.NotEnoughDetailsSpecified] + ": " + strings.Join(detailsList, ", "))
			}

			initStopNameTranslatorIfNecessary()
			printRoutesByLine := func(vehicleType string, lineNumber string) {
				lineRoutes, err := schedule.GetLine(vehicleType, lineNumber)
				if err != nil {
					log.Println(err.Error())
					return
				}

				fmt.Print(lineRoutes)
			}
			forEachLine(printRoutesByLine)

		case timetablesMode:
			forEachRouteByStop := func(stopCode string, f func(stopCode string, operationModeCode string, routeCode string)) {
				for _, operationModeCode := range operationModeCodes {
					for _, routeCode := range routeCodes {
						f(stopCode, operationModeCode, routeCode)
					}
				}
			}
			if len(vehicleTypes) > 0 && vehicleTypes[0] != "" && len(lineNumbers) > 0 && lineNumbers[0] != "" {
				initStopNameTranslatorIfNecessary()
				var printTimetableByLineStopCodeAndRoute func(line *schedule.Line, stopCode string, operationModeCode string, routeCode string)
				if len(stopCodes) == 1 && stopCodes[0] == "" || len(operationModeCodes) == 1 && operationModeCodes[0] == "" || len(routeCodes) == 1 && routeCodes[0] == "" {
					printTimetableByLineStopCodeAndRoute = func(line *schedule.Line, stopCode string, operationModeCode string, routeCode string) {
						stopTimetables, err := line.GetDetailedTimetableStrings(operationModeCode, routeCode, stopCode)
						if err != nil {
							log.Println(err.Error())
							return
						}

						for stopTimetable := range stopTimetables {
							fmt.Print(stopTimetable)
						}
					}
				} else {
					printTimetableByLineStopCodeAndRoute = func(line *schedule.Line, stopCode string, operationModeCode string, routeCode string) {
						stopTimetable, err := line.GetDetailedTimetableString(operationModeCode, routeCode, stopCode)
						if err != nil {
							log.Println(err.Error())
							return
						}

						fmt.Print(stopTimetable)
					}
				}
				lines := []*schedule.Line{}
				forEachLine(func(vehicleType string, lineNumber string) {
					line, err := schedule.GetLine(vehicleType, lineNumber)
					if err != nil {
						log.Println(err.Error())
						return
					}

					lines = append(lines, line)
				})
				if len(operationModeNames) > 0 && operationModeNames[0] != "" && len(routeNames) > 0 && routeNames[0] != "" && len(stopNames) > 0 && stopNames[0] != "" {
					forEachName := func(names []string, codes []string, nameToCodeTranslator func(string) string, customAction func(string)) {
						for _, name := range names {
							if name == "" {
								continue
							}

							translatedCode := nameToCodeTranslator(name)
							if translatedCode == "" {
								continue
							}

							isCodeInList := false
							for _, code := range codes {
								if code == translatedCode {
									isCodeInList = true
									break
								}
							}
							if !isCodeInList {
								codes = append(codes, translatedCode)
							}

							customAction(translatedCode)
						}
					}
					for _, line := range lines {
						operationModeMap := line.GetOperationModeMap()
						operationModeNameToCodeTranslator := func(operationModeName string) string {
							operationMode, ok := operationModeMap[operationModeName]
							if !ok {
								log.Printf("could not find operation mode code for name %s\n", operationModeName)
							}
							return operationMode.Code
						}

						forEachName(operationModeNames, operationModeCodes, operationModeNameToCodeTranslator, func(operationModeCode string) {
							routeMap := line.GetRouteMapByOperationMode(operationModeCode)
							routeNameToCodeTranslator := func(routeName string) string {
								route, ok := routeMap[routeName]
								if !ok {
									log.Printf("could not find route code for name %s\n", routeName)
								}
								return route.Code
							}

							forEachName(routeNames, routeCodes, routeNameToCodeTranslator, func(routeCode string) {
								stopMap := line.GetStopMapByOperationModeAndRoute(operationModeCode, routeCode)
								stopNameToCodeTranslator := func(stopName string) string {
									stop, ok := stopMap[stopName]
									if !ok {
										log.Printf("could not find stop code for name %s\n", stopName)
									}
									return stop.Code
								}

								forEachName(stopNames, stopCodes, stopNameToCodeTranslator, func(stopCode string) {})
							})
						})
					}
				}
				for _, line := range lines {
					for _, stopCode := range stopCodes {
						forEachRouteByStop(stopCode, func(stopCode string, operationModeCode string, routeCode string) {
							printTimetableByLineStopCodeAndRoute(line, stopCode, operationModeCode, routeCode)
						})
					}
				}
			} else {
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
		}
	} else {
		stopList, err := virtual.GetStops()
		if err != nil {
			log.Fatalln(err.Error())
		}

		if context.doSortStops {
			sort.Sort(stopList)
		}

		switch mode {
		case stopsMode:
			fmt.Print(stopList)

		case routesMode:
			routes, err := virtual.GetRoutes()
			if err != nil {
				log.Fatalln(err.Error())
			}

			stopMap := stopList.GetStopMap()
			if len(vehicleTypes) == 1 && len(lineNumbers) == 1 {
				lineRouteListList, err := routes.GetNamedRoutesByLine(vehicleTypes[0], lineNumbers[0], stopMap)
				if err != nil {
					log.Fatalln(err.Error())
				}

				fmt.Print(lineRouteListList)
			} else {
				routeMap, err := routes.GetRouteMap(stopMap)
				if err != nil {
					log.Fatalln(err.Error())
				}

				printRoutesByLine := func(vehicleType string, lineNumber string) {
					lineRoutes, ok := routeMap[virtual.Line{VehicleType: vehicleType, LineNumber: lineNumber}]
					if !ok {
						log.Printf("could not find line with vehicle type %s and number %s in the route map\n", vehicleType, lineNumber)
						return
					}

					fmt.Print(lineRoutes)
				}
				forEachLine(printRoutesByLine)
			}

		case timetablesMode:
			forEachLineByStop := func(stopCodeOrName string, f func(stopCodeOrName string, vehicleType string, lineNumber string)) {
				for _, vehicleType := range vehicleTypes {
					for _, lineNumber := range lineNumbers {
						f(stopCodeOrName, vehicleType, lineNumber)
					}
				}
			}
			printTimetableByStopCodeAndLine := func(stopCode string, vehicleType string, lineNumber string) {
				stopTimetable, err := virtual.GetTimetableByStopCodeAndLine(stopCode, context.vehicleTypesArg, context.lineNumbersArg)
				if err != nil {
					log.Println(err.Error())
					return
				}

				fmt.Print(stopTimetable)
			}
			if len(stopCodes) > 0 && stopCodes[0] != "" {
				for _, stopCode := range stopCodes {
					forEachLineByStop(stopCode, printTimetableByStopCodeAndLine)
				}
			}
			printTimetablesByStopNameAndLine := func(stopName string, vehicleType string, lineNumber string) {
				stopTimetables := stopList.GetTimetablesByStopNameAndLineAsync(stopName, context.vehicleTypesArg, context.lineNumbersArg, false)
				fmt.Print(stopTimetables)
			}
			if len(context.positionalArgs) > 0 {
				for _, stopName := range stopNames {
					forEachLineByStop(stopName, printTimetablesByStopNameAndLine)
				}
			} else {
				forEachLineByStop("", printTimetablesByStopNameAndLine)
			}
		}
	}
}
