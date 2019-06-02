package l10n

// EnglishTranslator maps names of terms in the reference language (i.e. English) to their translation in Bulgarian.
var EnglishTranslator = map[string]string{
	Usage: "usage: %s [-l lines] [-t types] [-s stop codes] [-showTime] [-showFacilities] [-sortStops] [stops]\n" +
		"          %s -showStops [-sortStops]\n" +
		"          %s -showRoutes [-l lines] [-t types] [-sortStops]\n" +
		"\n" +
		"The program outputs the timetables for the Sofia urban transit stops whose name partially or exactly matches one of the positional arguments (`stops`) or whose numerical code matches one of the `stop codes` passed as an optional argument.  If there are no positional arguments, timetables will be shown for all stops.  If `lines` are passed as an optional argument, only entries for the respective lines will be shown.  If `types` are passed as an optional argument, only entries for the respective vehicle types will be shown.\n" +
		"If the `-showStops` optional argument is passed, the program would just output a list of all stops and exit.\n" +
		"If the `-showRoutes` optional argument is passed, the program would just output a list of all routes and exit.  If `lines` are passed as an optional argument, only routes for the respective lines will be shown.  If `types` are passed as an optional argument, only routes for the respective vehicle types will be shown.\n" +
		"\n" +
		"Flags:\n",

	LineCodesFlagName:                          "l",
	LineCodesFlagUsage:                         "only output timetables for the specified comma-separated `lines`",
	VehicleTypesFlagName:                       "t",
	VehicleTypesFlagUsage:                      "only output timetables for the specified comma-separated vehicle `types` (\"%s\", \"%s\" or \"%s\")",
	StopCodesFlagName:                          "s",
	StopCodesFlagUsage:                         "output timetables for the stop with the specified comma-separated `stop codes` (in addition to stops passed as positional arguments)",
	DoShowGenerationTimeForTimetablesFlagName:  "showTime",
	DoShowGenerationTimeForTimetablesFlagUsage: "show generation time for each timetable",
	DoShowFacilitiesFlagName:                   "showFacilities",
	DoShowFacilitiesFlagUsage:                  `show detailed information about the facilities available in each vehicle ("%s" stands for "air conditioning" and "%s" stands for "wheelchair ramp slope")`,
	DoSortStopsFlagName:                        "sortStops",
	DoSortStopsFlagUsage:                       "sort stops by code",
	DoShowStopsFlagName:                        "showStops",
	DoShowStopsFlagUsage:                       "instead of outputting timetables, show the code and name of each stop",
	DoShowRoutesFlagName:                       "showRoutes",
	DoShowRoutesFlagUsage:                      "instead of outputting timetables, show the routes for all lines (or only for the ones passed as an argument of `-l`)",

	IncompatibleFlagsDetected: "incompatible flags detected",
}
