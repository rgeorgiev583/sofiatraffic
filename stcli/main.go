package main

import (
	"flag"
	"fmt"
	"os"

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
		usage := "употреба: %s [-л линия] [-т тип] [-покажиВреме] [-покажиУсловия] спирки\n"
		usage += "\n"
		usage += "Програмата извежда виртуалните табла за спирките на градския транспорт в София, чието име частично или изцяло съвпада с някой от подадените позиционни аргументи на командния ред (`спирки`).  Ако е зададена `линия` като опционален аргумент, ще бъдат изведени само записите за превозните средства от конкретната линия.  Ако е зададен `тип` като опционален аргумент, ще бъдат изведени само записите за превозните средства от конкретния тип.\n"
		usage += "\n"
		usage += "Флагове:\n"
		fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])
		flag.PrintDefaults()
	}

	var lineCode string
	flag.StringVar(&lineCode, "л", "", "да се изведат само виртуалните табла за превозните средства от конкретната `линия`")

	var vehicleType string
	flag.StringVar(&vehicleType, "т", "", "да се изведат само виртуалните табла за превозните средства от конкретния `тип` (\"автобус\", \"тролейбус\" или \"трамвай\")")

	flag.BoolVar(&regular.DoShowGenerationTimeForTimetables, "покажиВреме", false, `да се покаже времето на генериране на всяко виртуално табло`)

	flag.BoolVar(&regular.DoShowFacilities, "покажиУсловия", false, `да се покажат подробности за условията в превозните средства (чрез "К" се обозначава дали има климатик в превозното средство, а чрез "И" - дали има рампа за инвалидни колички)`)

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "грешка: трябва да бъде подаден поне един позиционен аргумент")
		fmt.Fprintf(os.Stderr, "Изпълнете '%s -h' за указания.\n", os.Args[0])
		os.Exit(1)
	}

	regular.GenerationTimeLabel = "време на генериране"
	regular.VehicleTypeTranslator = translateVehicleTypeFromEnglishToBulgarian

	stopList, err := regular.GetStops()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	vehicleType = translateVehicleTypeFromBulgarianToEnglish(vehicleType)

	for _, stopName := range args {
		stopTimetables := stopList.GetTimetablesByStopNameAndLineAsync(stopName, vehicleType, lineCode, false)
		fmt.Print(stopTimetables)
	}
}
