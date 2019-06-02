package l10n

// BulgarianTranslator maps names of terms in the reference language (i.e. English) to their translation in Bulgarian.
var BulgarianTranslator = map[string]string{
	Usage: "употреба: %s [-л номера на линии] [-т типове превозни средства] [-с кодове на спирки] [-покажиВреме] [-покажиУсловия] [-сортирайСпирки] [имена на спирки]\n" +
		"          %s -покажиСпирки [-сортирайСпирки]\n" +
		"          %s -покажиМаршрути -л номера на линии [-т типове превозни средства] [-сортирайСпирки]\n" +
		"\n" +
		"Програмата извежда виртуалните табла за спирките на градския транспорт в София, чието име частично или изцяло съвпада с някое от всичките `имена на спирки`, подадени като позиционни аргументи на командния ред, или чийто код съвпада с някой от зададените чрез опционален аргумент `кодове на спирки`.  Ако не са подадени позиционни аргументи, ще бъдат показани виртуалните табла за всички спирки.  Ако са зададени `кодове на линии` чрез опционален аргумент, ще бъдат изведени само записите за конкретните линии.  Ако са зададени `типове превозни средства` чрез опционален аргумент, ще бъдат изведени само записите за превозните средства от конкретните типове.\n" +
		"Ако е подаден опционалният аргумент `-покажиСпирки`, вместо това програмата ще изведе списък със всички спирки и ще приключи.\n" +
		"Ако е подаден опционалният аргумент `-покажиМаршрути`, вместо това програмата ще изведе списък със всички маршрути и ще приключи.  Ако са зададени `номера на линии` чрез опционален аргумент, ще бъдат изведени маршрутите на конкретните линии.  Ако са зададени `типове превозни средства` чрез опционален аргумент, ще бъдат изведени само маршрутите на превозните средства от конкретните типове.\n" +
		"\n" +
		"Опционални аргументи:\n",

	LineNumbersFlagName:                        "л",
	LineNumbersFlagUsage:                       "да се изведат само виртуалните табла за превозните средства от конкретните `линии`, разделени със запетая",
	VehicleTypesFlagName:                       "т",
	VehicleTypesFlagUsage:                      "да се изведат само виртуалните табла за превозните средства от конкретните `типове` (\"%s\", \"%s\" или \"%s\"), разделени със запетая",
	StopCodesFlagName:                          "с",
	StopCodesFlagUsage:                         "да се изведат виртуалните табла за спирките със зададените `кодове`, разделени със запетая (в допълнение към спирките, зададени чрез позиционни аргументи)",
	DoShowGenerationTimeForTimetablesFlagName:  "покажиВреме",
	DoShowGenerationTimeForTimetablesFlagUsage: "да се покаже времето на генериране на всяко виртуално табло",
	DoShowFacilitiesFlagName:                   "покажиУсловия",
	DoShowFacilitiesFlagUsage:                  `да се покажат подробности за условията в превозните средства (чрез "%s" се обозначава дали има климатик в превозното средство, а чрез "%s" - дали има рампа за инвалидни колички)`,
	DoSortStopsFlagName:                        "сортирайСпирки",
	DoSortStopsFlagUsage:                       "да се подредят спирките по код",
	DoShowStopsFlagName:                        "покажиСпирки",
	DoShowStopsFlagUsage:                       "вместо да се извеждат виртуалните табла, да се изведат по двойки кодовете и имената на всички спирки",
	DoShowRoutesFlagName:                       "покажиМаршрути",
	DoShowRoutesFlagUsage:                      "вместо да се извеждат виртуалните табла, да се изведат двата маршрута на зададената линия",
	DoTranslateStopNamesFlagName:               "преведиИменатаНаСпирките",
	DoTranslateStopNamesFlagUsage:              "да се преведат имената на спирките от български на локалния език",

	IncompatibleFlagsDetected: "подадени са несъвместими опционални аргументи",
	NoLineSpecified:           "не е зададена линия",
}