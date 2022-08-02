package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const JsonFile = "data.json"

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func main() {
	var departureStation string
	fmt.Println("Введи станцію відїзду.")
	fmt.Scanf("%s\n", &departureStation)

	var arrivalStation string
	fmt.Println("введи станцію приїзду.")
	fmt.Scanf("%s\n", &arrivalStation)

	var criteria string
	fmt.Println("Введи критерій price | arrival-time | departure-time")
	fmt.Scanf("%s\n", &criteria)

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(result); i++ {
		fmt.Printf("ID %v:  Price: %.1f   Arrivaltime: %s    Deptime: %s \n",
			result[i].TrainID,
			result[i].Price,
			result[i].ArrivalTime.Format("15:04:05"),
			result[i].DepartureTime.Format("15:04:05"))
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	AllTrains := ParseJson(JsonFile)
	arrivnum, errarr := strconv.Atoi(arrivalStation)
	deperturenum, errdep := strconv.Atoi(departureStation)

	if departureStation == "" {
		return nil, errors.New("empty departure station")
	}
	if arrivalStation == "" {
		return nil, errors.New("empty arrival station")
	}

	if errarr != nil || arrivnum < 1 {
		return nil, errors.New("bad arrival station input")
	}
	if errdep != nil || deperturenum < 1 {
		return nil, errors.New("bad departure station input")
	}

	if criteria != "price" && criteria != "arrival-time" && criteria != "departure-time" {
		return nil, errors.New("unsupported criteria")
	}
	RigthTrains := Trains{}
	for _, tinfo := range AllTrains {
		if tinfo.DepartureStationID == deperturenum && tinfo.ArrivalStationID == arrivnum {
			RigthTrains = append(RigthTrains, tinfo)
		}
	}
	switch criteria {
	case "price":
		sort.SliceStable(RigthTrains, func(i, j int) bool {
			return RigthTrains[i].Price < RigthTrains[j].Price
		})
	case "arrival-time":
		sort.SliceStable(RigthTrains, func(i, j int) bool {
			return RigthTrains[i].ArrivalTime.Before(RigthTrains[j].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(RigthTrains, func(i, j int) bool {
			return RigthTrains[i].DepartureTime.Before(RigthTrains[j].DepartureTime)
		})
	}
	if len(RigthTrains) == 0 {
		return nil, nil
	} else if len(RigthTrains) >= 3 {
		return RigthTrains[:3], nil
	}
	return RigthTrains, nil

}
func ParseJson(filename string) []Train {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//custom structure to parse json
	type parsetrain struct {
		TrainID            int     `json:"trainId"`
		DepartureStationID int     `json:"departureStationId"`
		ArrivalStationID   int     `json:"arrivalStationId"`
		Price              float32 `json:"price"`
		ArrivalTime        string  `json:"arrivalTime"`
		DepartureTime      string  `json:"departureTime"`
	}

	var UnmarshalTrains []parsetrain
	var AllTrains []Train

	_ = json.Unmarshal(byteValue, &UnmarshalTrains)
	//rewriting to right structure
	for i := 0; i < len(UnmarshalTrains); i++ {
		var newtrain Train
		newtrain.TrainID = UnmarshalTrains[i].TrainID
		newtrain.DepartureStationID = UnmarshalTrains[i].DepartureStationID
		newtrain.ArrivalStationID = UnmarshalTrains[i].ArrivalStationID
		newtrain.Price = UnmarshalTrains[i].Price
		arrtime := strings.Split(UnmarshalTrains[i].ArrivalTime, ":")
		deptime := strings.Split(UnmarshalTrains[i].DepartureTime, ":")
		newtrain.DepartureTime, newtrain.ArrivalTime = setTime(arrtime, deptime)
		AllTrains = append(AllTrains, newtrain)
	}
	return AllTrains
}
func setTime(arrtime, deptime []string) (time.Time, time.Time) {

	dephour, _ := strconv.Atoi(deptime[0])
	depmin, _ := strconv.Atoi(deptime[1])
	depsec, _ := strconv.Atoi(deptime[2])
	//set deptime
	traindeptime := time.Date(0, time.January, 1, dephour, depmin, depsec, 0, time.UTC)
	arrhour, _ := strconv.Atoi(arrtime[0])
	arrmin, _ := strconv.Atoi(arrtime[1])
	arrsec, _ := strconv.Atoi(arrtime[2])
	//setarrtime
	trainarr := time.Date(0, time.January, 1, arrhour, arrmin, arrsec, 0, time.UTC)
	return traindeptime, trainarr
}
