package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/regular"
)

func translateVehicleTypeFromBulgarianToEnglish(vehicleType string) string {
	switch vehicleType {
	case "автобус":
		return "bus"

	case "тролейбус":
		return "trolley"

	case "трамвай":
		return "tram"
	}

	return ""
}

func translateVehicleTypeFromEnglishToBulgarian(vehicleType string) string {
	switch vehicleType {
	case "bus":
		return "автобус"

	case "trolley":
		return "тролейбус"

	case "tram":
		return "трамвай"
	}

	return ""
}

func main() {
	flag.Usage = func() {
		usage := "употреба: %s [-л линии] [-т типове] [-с кодове на спирки] [-покажиВреме] [-покажиУсловия] [-сортирайСпирки] [спирки]\n"
		usage += "          %s -покажиСпирки [-сортирайСпирки]\n"
		usage += "\n"
		usage += "Програмата извежда виртуалните табла за спирките на градския транспорт в София, чието име частично или изцяло съвпада с някой от подадените позиционни аргументи на командния ред (`спирки`).  Ако не са подадени позиционни аргументи, ще бъдат показани виртуалните табла за всички спирки.  Ако е зададена `линия` като опционален аргумент, ще бъдат изведени само записите за превозните средства от конкретната линия.  Ако е зададен `тип` като опционален аргумент, ще бъдат изведени само записите за превозните средства от конкретния тип.\n"
		usage += "\n"
		usage += "Флагове:\n"
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var lineCodesArg string
	flag.StringVar(&lineCodesArg, "л", "", "да се изведат само виртуалните табла за превозните средства от конкретните `линии`, разделени със запетая")

	var vehicleTypesArg string
	flag.StringVar(&vehicleTypesArg, "т", "", "да се изведат само виртуалните табла за превозните средства от конкретните `типове` (\"автобус\", \"тролейбус\" или \"трамвай\"), разделени със запетая")

	var stopCodesArg string
	flag.StringVar(&stopCodesArg, "с", "", "да се изведат виртуалните табла за спирките със зададените `кодове`, разделени със запетая (в допълнение към спирките, зададени чрез позиционни аргументи)")

	flag.BoolVar(&regular.DoShowGenerationTimeForTimetables, "покажиВреме", false, `да се покаже времето на генериране на всяко виртуално табло`)

	flag.BoolVar(&regular.DoShowFacilities, "покажиУсловия", false, `да се покажат подробности за условията в превозните средства (чрез "К" се обозначава дали има климатик в превозното средство, а чрез "И" - дали има рампа за инвалидни колички)`)

	var doSortStops bool
	flag.BoolVar(&doSortStops, "сортирайСпирки", false, "да се подредят спирките по код")

	var doShowStops bool
	flag.BoolVar(&doShowStops, "покажиСпирки", false, "вместо да се извеждат виртуалните табла, да се изведат по двойки кодовете и имената на всички спирки")

	flag.Parse()

	args := flag.Args()

	regular.GenerationTimeLabel = "време на генериране"
	regular.VehicleTypeTranslator = translateVehicleTypeFromEnglishToBulgarian

	lineCodes := strings.Split(lineCodesArg, ",")
	if lineCodesArg != "" {
		for i, lineCode := range lineCodes {
			lineCodes[i] = strings.TrimSpace(lineCode)
		}
	}

	vehicleTypes := strings.Split(vehicleTypesArg, ",")
	if vehicleTypesArg != "" {
		for i, vehicleType := range vehicleTypes {
			vehicleTypes[i] = translateVehicleTypeFromBulgarianToEnglish(strings.TrimSpace(vehicleType))
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
