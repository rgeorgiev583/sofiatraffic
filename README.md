# sofiatraffic
Tools for fetching information about public transit in Sofia from the Urban Mobility Centre Web APIs.

## Installation of CLI tool

        go get github.com/rgeorgiev583/sofiatraffic/stcli

## Usage of CLI tool

        stcli -h

## TODO

* scraping of information from schedules in addition to virtual timetables (i.e. timetables, stops, routes, etc.)
  * regular transport
  * metro
* remaining time until arrival
* trip guru integration
* listing of nearest stops via geolocation
* GUI
  * nearest stops: map service integration
* listing of route change history
* local caching and exporting of timetable/stop/route data
* reading timetable/stop/route data from file