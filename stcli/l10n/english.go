package l10n

// EnglishTranslator maps names of terms in the reference language (i.e. English) to their translation in English.
var EnglishTranslator = map[string]string{
	VehicleTypeBus:        "bus",
	VehicleTypeTrolleybus: "trolleybus",
	VehicleTypeTram:       "tram",
	VehicleTypeMetro:      "metro",

	AirConditioningAbbreviation:         "A",
	WheelchairAccessibilityAbbreviation: "W",

	Usage: "usage: %s [-l line numbers] [-t vehicle types] [-s stop codes] [-r route codes] [-o operation mode codes] [-showTime] [-showFacilities] [-showRoute] [-showOperationMode] [-useSchedule] [-sortStops] [-translateStopNames] [stop names]\n" +
		"       %s -showStops [-sortStops] [-translateStopNames]\n" +
		"       %s -showLines [-useSchedule]\n" +
		"       %s -showRoutes -l line numbers [-t vehicle types] [-useSchedule] [-sortStops] [-translateStopNames]\n" +
		"\n" +
		"The program outputs the timetables for the Sofia urban transit stops whose name partially or exactly matches one of the `stop names` passed as positional arguments.  In addition, it outputs the timetables for stops whose numerical code matches one of the `stop codes` passed as an optional argument.  If there are no positional arguments, timetables will be shown for all stops.  If `line numbers` are passed as an optional argument, only entries for the respective lines will be shown.  If `vehicle types` are passed as an optional argument, only entries for the respective vehicle types will be shown.\n" +
		"If the `-showStops` optional argument is passed, the program would just output a list of all stops and exit.\n" +
		"If the `-showLines` optional argument is passed, the program would just output a list of all line numbers (grouped by vehicle type) and exit.\n" +
		"If the `-showRoutes` optional argument is passed, the program would just output a list of all routes and exit.  If `line numbers` are passed as an optional argument, only routes for the respective lines will be shown.  If `vehicle types` are passed as an optional argument, only routes for the respective vehicle types will be shown.\n" +
		"\n" +
		"Flags:\n",

	LineNumbersFlagName:                        "l",
	LineNumbersFlagUsage:                       "only output timetables for vehicles with the specified comma-separated `line numbers`",
	VehicleTypesFlagName:                       "t",
	VehicleTypesFlagUsage:                      "only output timetables for vehicles of the specified comma-separated `vehicle types` (\"%s\", \"%s\" or \"%s\")",
	StopCodesFlagName:                          "s",
	StopCodesFlagUsage:                         "only output timetables for stops with the specified comma-separated `stop codes` (in addition to stops passed as positional arguments)",
	RouteCodesFlagName:                         "r",
	RouteCodesFlagUsage:                        "only output timetables for routes with the specified comma-separated `route codes`",
	OperationModeCodesFlagName:                 "o",
	OperationModeCodesFlagUsage:                "only output timetables for modes of operation with the specified comma-separated `operation mode codes`",
	DoShowGenerationTimeForTimetablesFlagName:  "showTime",
	DoShowGenerationTimeForTimetablesFlagUsage: "show the generation time for each timetable",
	DoShowFacilitiesFlagName:                   "showFacilities",
	DoShowFacilitiesFlagUsage:                  `show detailed information about the facilities available in each vehicle ("%s" stands for "air conditioning" and "%s" stands for "wheelchair ramp slope")`,
	DoShowRouteFlagName:                        "showRoute",
	DoShowRouteFlagUsage:                       "show the name and code of the route for schedule timetables",
	DoShowOperationModeFlagName:                "showOperationMode",
	DoShowOperationModeFlagUsage:               "show the name and code of the operation mode for schedule timetables",
	DoUseScheduleFlagName:                      "useSchedule",
	DoUseScheduleFlagUsage:                     "fetch information from the schedule pages on the Urban Mobility Centre website instead of the REST API for virtual timetables",
	DoSortStopsFlagName:                        "sortStops",
	DoSortStopsFlagUsage:                       "sort list of stops by code internally",
	DoTranslateStopNamesFlagName:               "translateStopNames",
	DoTranslateStopNamesFlagUsage:              "translate names of stops from Bulgarian to the local language",
	DoShowStopsFlagName:                        "showStops",
	DoShowStopsFlagUsage:                       "instead of outputting timetables, show the code and name of each stop",
	DoShowLinesFlagName:                        "showLines",
	DoShowLinesFlagUsage:                       "instead of outputting timetables, show the numbers of all lines grouped by vehicle type (implies `-useSchedule`)",
	DoShowRoutesFlagName:                       "showRoutes",
	DoShowRoutesFlagUsage:                      "instead of outputting timetables, show the routes for the specified line",

	IncompatibleFlagsDetected: "incompatible flags detected",
	NoLineSpecified:           "no line specified",
}
