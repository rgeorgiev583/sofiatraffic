package l10n

// EnglishTranslator maps names of terms in the reference language (i.e. English) to their translation in English.
var EnglishTranslator = map[string]string{
	VehicleTypeBus:        "bus",
	VehicleTypeTrolleybus: "trolleybus",
	VehicleTypeTram:       "tram",
	VehicleTypeMetro:      "metro",

	AirConditioningAbbreviation:         "A",
	WheelchairAccessibilityAbbreviation: "W",

	Usage: "%s is a tool for fetching information about public transit in Sofia.\n" +
		"\n" +
		"Usage:\n" +
		"\n" +
		"        %s <command> [arguments]\n" +
		"\n" +
		"The commands are:\n" +
		"\n" +
		"        timetables    show urban transit timetables\n" +
		"        stops         show urban transit stops\n" +
		"        lines         show urban transit lines\n" +
		"        routes        show urban transit routes\n" +
		"\n" +
		"Use \"%s <command> -h\" for more information about a command.\n",

	TimetablesSubcommandName: "timetables",
	TimetablesSubcommandUsage: "usage: %s timetables [-l line numbers] [-t vehicle types] [-s stop codes] [-r route codes] [-o operation mode codes] [-showTime] [-showFacilities] [-showRoute] [-showOperationMode] [-useSchedule] [-sortStops] [-translateStopNames] [stop names]\n" +
		"\n" +
		"Timetables shows the timetables for Sofia urban transit stops whose name partially or exactly matches one of the `stop names` passed as positional arguments. In addition, it shows the timetables for stops whose numerical code matches one of the `stop codes` passed as an optional argument.\n" +
		"If there are no positional arguments, timetables will be shown for all stops. If `line numbers` are passed as an optional argument, only entries for the respective lines will be shown. If `vehicle types` are passed as an optional argument, only entries for the respective vehicle types will be shown.\n" +
		"\n" +
		"Flags:\n",
	StopsSubcommandName: "stops",
	StopsSubcommandUsage: "usage: %s stops [-sortStops] [-translateStopNames]\n" +
		"\n" +
		"Stops shows a list containing the code and name of each stop.\n" +
		"\n" +
		"Flags:\n",
	LinesSubcommandName: "lines",
	LinesSubcommandUsage: "usage: %s lines\n" +
		"\n" +
		"Lines shows a list containing the numbers of all lines grouped by vehicle type.\n",
	RoutesSubcommandName: "routes",
	RoutesSubcommandUsage: "usage: %s routes -l line numbers [-t vehicle types] [-useSchedule] [-sortStops] [-translateStopNames]\n" +
		"\n" +
		"Routes shows the routes for each line. If `line numbers` are passed as an optional argument, only routes for the respective lines will be shown. If `vehicle types` are passed as an optional argument, only routes for the respective vehicle types will be shown.\n" +
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
	DoShowGenerationTimeForTimetablesFlagName:  "showGenerationTime",
	DoShowGenerationTimeForTimetablesFlagUsage: "show the generation time for each timetable",
	DoShowRemainingTimeUntilArrivalFlagName:    "showRemainingTime",
	DoShowRemainingTimeUntilArrivalFlagUsage:   "show the remaining time until the arrival of each vehicle instead of the specific time of arrival",
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

	InvalidSubcommandName:     "invalid command name",
	IncompatibleFlagsDetected: "incompatible flags detected",
	NoLineSpecified:           "no line specified",
}
