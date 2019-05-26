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
		fmt.Fprintf(flag.CommandLine.Output(), `употреба: %s [-л линия] [-т тип] частично или цяло име на спирка

Програмата извежда виртуалните табла за спирките на градския транспорт в София, чието име частично или изцяло съвпада с
подаденото като аргумент на командния ред.  Ако е зададена `+"`"+`линия`+"`"+` като опционален аргумент, ще бъдат 
изведени само записите за превозните средства от конкретната линия.  Ако е зададен `+"`"+`тип`+"`"+` като опционален
аргумент, ще бъдат изведени само записите за превозните средства от конкретния тип.

Флагове:
`, os.Args[0])
		flag.PrintDefaults()
	}

	var lineCode string
	flag.StringVar(&lineCode, "л", "", `да се изведат само виртуалните табла за превозните средства от конкретната линия`)

	var vehicleType string
	flag.StringVar(&vehicleType, "т", "", `да се изведат само виртуалните табла за превозните средства от конкретния тип ("автобус", "тролейбус" или "трамвай")`)

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "грешка: не е подаден аргумент")
		fmt.Fprintf(os.Stderr, "Изпълнете '%s -h' за указания.\n", os.Args[0])
		os.Exit(1)
	}

	regular.VehicleTypeTranslator = translateVehicleTypeFromEnglishToBulgarian

	stopNamePattern := args[0]
	stopList, err := regular.GetStops()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	line := &regular.Line{
		Code:        lineCode,
		VehicleType: translateVehicleTypeFromBulgarianToEnglish(vehicleType),
	}
	stopArrivalMap, err := stopList.MatchArrivalsByStopNameAndLine(stopNamePattern, line)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Print(stopArrivalMap)
}