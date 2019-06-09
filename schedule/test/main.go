package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rgeorgiev583/sofiatraffic/schedule"
)

func main() {
	lines, err := schedule.GetLines()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Printf("%#v", lines)
	if len(os.Args) >= 4 {
		timetable, err := schedule.GetTimetable(os.Args[1], os.Args[2], os.Args[3])
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Printf("%#v", timetable)
	} else if len(os.Args) == 3 {
		line, err := schedule.GetLine(os.Args[1], os.Args[2])
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Printf("%#v", line)
	}
}
