package l10n

// BulgarianTranslator maps names of terms in the reference language (i.e. English) to their translation in Bulgarian.
var BulgarianTranslator = map[string]string{
	VehicleTypeBus:        "автобус",
	VehicleTypeTrolleybus: "тролейбус",
	VehicleTypeTram:       "трамвай",
	VehicleTypeMetro:      "метро",

	AirConditioningAbbreviation:         "К",
	WheelchairAccessibilityAbbreviation: "И",

	Usage: "%s е инструмент за достъпване на информация за обществения транспорт в София.\n" +
		"\n" +
		"Употреба:\n" +
		"\n" +
		"        %s <команда> [аргументи]\n" +
		"\n" +
		"Командите са:\n" +
		"\n" +
		"        табла       показва времената на пристигане на градския транспорт\n" +
		"        спирки      показва спирките на градския транспорт\n" +
		"        линии       показва линиите на градския транспорт\n" +
		"        маршрути    показва маршрутите на градския транспорт\n" +
		"\n" +
		"Използвайте \"%s <команда> -h\" за повече информация за дадената команда.\n",

	TimetablesSubcommandName: "табла",
	TimetablesSubcommandUsage: "употреба: %s табла [-л номера на линии] [-т типове превозни средства] [-с кодове на спирки] [-м кодове на маршрути] [-р кодове на режими] [-покажиВремеНаГенериране] [-покажиОставащоВреме] [-покажиУсловия] [-покажиМаршрут] [-покажиРежим] [-използвайРазписание] [-сортирайСпирки] [-преведиИменаНаСпирки] [имена на спирки]\n" +
		"\n" +
		"Табла извежда времената на пристигане за спирките на градския транспорт в София, чието име частично или изцяло съвпада с някое от всичките `имена на спирки`, подадени като позиционни аргументи на командния ред. Освен това тя извежда времената на пристигане за спирките, чийто код съвпада с някой от зададените чрез опционален аргумент `кодове на спирки`.\n" +
		"Ако не са подадени позиционни аргументи, ще бъдат показани времената на пристигане за всички спирки. Ако са зададени `номера на линии` чрез опционален аргумент, ще бъдат изведени само записите за конкретните линии. Ако са зададени `типове превозни средства` чрез опционален аргумент, ще бъдат изведени само записите за превозните средства от конкретните типове.\n" +
		"\n" +
		"Опционални аргументи:\n",
	StopsSubcommandName: "спирки",
	StopsSubcommandUsage: "употреба: %s спирки [-сортирайСпирки] [-преведиИменаНаСпирки]\n" +
		"\n" +
		"Спирки показва списък, съдържащ кодовете и имената на всички спирки.\n" +
		"\n" +
		"Опционални аргументи:\n",
	LinesSubcommandName: "линии",
	LinesSubcommandUsage: "употреба: %s линии\n" +
		"\n" +
		"Линии показва списък, съдържащ номерата на всички линии, групирани по тип на превозното средство.\n",
	RoutesSubcommandName: "маршрути",
	RoutesSubcommandUsage: "употреба: %s маршрути -л номера на линии [-т типове превозни средства] [-използвайРазписание] [-сортирайСпирки] [-преведиИменаНаСпирки]\n" +
		"\n" +
		"Маршрути показва маршрутите за всяка линия. Ако е извикана подкомандата `маршрути`, програмата просто ще изведе списък, съдържащ маршрутите на всички линии, и ще приключи. Ако са зададени `номера на линии` чрез опционален аргумент, ще бъдат изведени само маршрутите на конкретните линии. Ако са зададени `типове превозни средства` чрез опционален аргумент, ще бъдат изведени само маршрутите на превозните средства от конкретните типове.\n" +
		"\n" +
		"Опционални аргументи:\n",

	LineNumbersFlagName:                        "л",
	LineNumbersFlagUsage:                       "да се изведат времената на пристигане само за превозни средства със зададените `номера на линии`, разделени със запетая",
	VehicleTypesFlagName:                       "т",
	VehicleTypesFlagUsage:                      "да се изведат времената на пристигане само за превозни средства от зададените `типове превозни средства` (\"%s\", \"%s\" или \"%s\"), разделени със запетая",
	StopCodesFlagName:                          "с",
	StopCodesFlagUsage:                         "да се изведат времената на пристигане само за спирки със зададените `кодове на спирки`, разделени със запетая (в допълнение към спирките, зададени чрез позиционни аргументи)",
	RouteCodesFlagName:                         "м",
	RouteCodesFlagUsage:                        "да се изведат времената на пристигане само за превозни средства, минаващи по маршрути със зададените `кодове на маршрути`, разделени със запетая",
	OperationModeCodesFlagName:                 "р",
	OperationModeCodesFlagUsage:                "да се изведат времената на пристигане само по време на режими със зададените `кодове на режими`, разделени със запетая",
	DoShowGenerationTimeForTimetablesFlagName:  "покажиВремеНаГенериране",
	DoShowGenerationTimeForTimetablesFlagUsage: "да се покаже времето на генериране на всяко виртуално табло",
	DoShowRemainingTimeUntilArrivalFlagName:    "покажиОставащоВреме",
	DoShowRemainingTimeUntilArrivalFlagUsage:   "покажи оставащото време до пристигането на всяко превозно средство вместо конкретното време на пристигане",
	DoShowFacilitiesFlagName:                   "покажиУсловия",
	DoShowFacilitiesFlagUsage:                  `да се покажат подробности за условията във всяко превозно средство (чрез "%s" се обозначава дали има климатик, а чрез "%s" - дали има рампа за инвалидни колички)`,
	DoShowRouteFlagName:                        "покажиМаршрут",
	DoShowRouteFlagUsage:                       "да се покажат името и кодът на маршрута за всяко разписание",
	DoShowOperationModeFlagName:                "покажиРежим",
	DoShowOperationModeFlagUsage:               "да се покажат името и кодът на режима на всяко разписание",
	DoUseScheduleFlagName:                      "използвайРазписание",
	DoUseScheduleFlagUsage:                     "да се използва информацията от страниците с разписанието на сайта на Центъра за градска мобилност вместо от REST API-то за виртуалните табла",
	DoSortStopsFlagName:                        "сортирайСпирки",
	DoSortStopsFlagUsage:                       "да се подредят вътрешно спирките по код",
	DoTranslateStopNamesFlagName:               "преведиИменаНаСпирки",
	DoTranslateStopNamesFlagUsage:              "да се преведат имената на спирките от български на локалния език",

	InvalidSubcommandName:     "невалидно име на команда",
	IncompatibleFlagsDetected: "подадени са несъвместими опционални аргументи",
	NoLineSpecified:           "не е зададена линия",
}
