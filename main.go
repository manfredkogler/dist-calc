package main

import (
	"dist-calc/requests"
)

func main() {
	distance := requests.Distance(requests.Locations[requests.Lambertgasse], requests.Locations[requests.Schottenring])
	distance = requests.Distance(requests.Locations[requests.Schottenring], requests.Locations[requests.Lambertgasse])
	_ = distance

	location := requests.ForwardGeocode("Schottenring 1 Wien")
	location = requests.ForwardGeocode("Lambertgasse 4 Wien")
	_ = location
}
