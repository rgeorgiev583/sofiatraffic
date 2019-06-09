package l10n

// EnglishTranslator maps names of terms in the reference language (i.e. English) to their translation in Bulgarian.
var EnglishTranslator = map[string]string{
	VehicleTypeBus:        "bus",
	VehicleTypeTrolleybus: "trolleybus",
	VehicleTypeTram:       "tram",
	VehicleTypeMetro:      "metro",

	AirConditioningAbbreviation:         "A",
	WheelchairAccessibilityAbbreviation: "W",

	Usage: "usage: %s [-l line numbers] [-t vehicle types] [-s stop codes] [-showTime] [-showFacilities] [-sortStops] [-translateStopNames] [stop names]\n" +
		"       %s -showStops [-sortStops] [-translateStopNames]\n" +
		"       %s -showRoutes -l line numbers [-t vehicle types] [-sortStops] [-translateStopNames]\n" +
		"\n" +
		"The program outputs the timetables for the Sofia urban transit stops whose name partially or exactly matches one of the `stop names` passed as positional arguments or whose numerical code matches one of the `stop codes` passed as an optional argument.  If there are no positional arguments, timetables will be shown for all stops.  If `line numbers` are passed as an optional argument, only entries for the respective lines will be shown.  If `vehicle types` are passed as an optional argument, only entries for the respective vehicle types will be shown.\n" +
		"If the `-showStops` optional argument is passed, the program would just output a list of all stops and exit.\n" +
		"If the `-showRoutes` optional argument is passed, the program would just output a list of all routes and exit.  If `line numbers` are passed as an optional argument, only routes for the respective lines will be shown.  If `vehicle types` are passed as an optional argument, only routes for the respective vehicle types will be shown.\n" +
		"\n" +
		"Flags:\n",

	LineNumbersFlagName:                        "l",
	LineNumbersFlagUsage:                       "only output timetables for vehicles with the specified comma-separated `line numbers`",
	VehicleTypesFlagName:                       "t",
	VehicleTypesFlagUsage:                      "only output timetables for vehicles of the specified comma-separated `vehicle types` (\"%s\", \"%s\" or \"%s\")",
	StopCodesFlagName:                          "s",
	StopCodesFlagUsage:                         "output timetables for stops with the specified comma-separated `stop codes` (in addition to stops passed as positional arguments)",
	RouteCodesFlagName:                         "r",
	RouteCodesFlagUsage:                        "output timetables for the routes with the specified comma-separated `route codes`",
	OperationModeCodesFlagName:                 "o",
	OperationModeCodesFlagUsage:                "output timetables for the modes of operation with the specified comma-separated `operation mode codes`",
	DoShowGenerationTimeForTimetablesFlagName:  "showTime",
	DoShowGenerationTimeForTimetablesFlagUsage: "show generation time for each timetable",
	DoShowFacilitiesFlagName:                   "showFacilities",
	DoShowFacilitiesFlagUsage:                  `show detailed information about the facilities available in each vehicle ("%s" stands for "air conditioning" and "%s" stands for "wheelchair ramp slope")`,
	DoShowOperationModeFlagName:                "showOperationMode",
	DoShowOperationModeFlagUsage:               "show the name and code of the operation mode for schedule timetables",
	DoShowRouteFlagName:                        "showRoute",
	DoShowRouteFlagUsage:                       "show the name and code of the route for schedule timetables",
	DoUseScheduleFlagName:                      "useSchedule",
	DoUseScheduleFlagUsage:                     "fetch information from the schedule pages on the Urban Mobility Centre website instead of the REST API for the virtual tables",
	DoSortStopsFlagName:                        "sortStops",
	DoSortStopsFlagUsage:                       "sort stops by code",
	DoTranslateStopNamesFlagName:               "translateStopNames",
	DoTranslateStopNamesFlagUsage:              "translate names of stops from Bulgarian to the local language",
	DoShowStopsFlagName:                        "showStops",
	DoShowStopsFlagUsage:                       "instead of outputting timetables, show the code and name of each stop",
	DoShowLinesFlagName:                        "showLines",
	DoShowLinesFlagUsage:                       "instead of outputting timetables, show the numbers of all lines",
	DoShowRoutesFlagName:                       "showRoutes",
	DoShowRoutesFlagUsage:                      "instead of outputting timetables, show the routes for the specified line",

	IncompatibleFlagsDetected: "incompatible flags detected",
	NoLineSpecified:           "no line specified",
}
