package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"strings"
)

func main() {
	elasticClient, err := elastic.NewSimpleClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile("../vaers.json")
	if err != nil {
		panic(err)
	}

	events := []VaersEvent{}
	err = json.Unmarshal(b, &events)
	if err != nil {
		panic(err)
	}

	for _, e := range events {
		j, _ := json.Marshal(e)

		if len(e.VAERSID) > 0 {
			_, err = elasticClient.Index().Index("vaers").Type("event").Id(strings.ToLower(e.VAERSID)).BodyJson(string(j)).Do(context.Background())
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println(string(j))
		}
	}

	defer elasticClient.Stop()

}

type VaersEvent struct {
	Notes                               string `json:"Notes"`
	Vaccine                             string `json:"Vaccine"`
	VaccineCode                         string `json:"Vaccine Code"`
	VAERSID                             string `json:"VAERS ID"`
	VAERSIDCode                         string `json:"VAERS ID Code"`
	VaccineManufacturer                 string `json:"Vaccine Manufacturer"`
	VaccineManufacturerCode             string `json:"Vaccine Manufacturer Code"`
	VaccineLot                          string `json:"Vaccine Lot"`
	VaccineLotCode                      string `json:"Vaccine Lot Code"`
	AdverseEventDescription             string `json:"Adverse Event Description"`
	LabData                             string `json:"Lab Data"`
	CurrentIllness                      string `json:"Current Illness"`
	AdverseEventsAfterPriorVaccinations string `json:"Adverse Events After Prior Vaccinations"`
	MedicationsAtTimeOfVaccination      string `json:"Medications At Time Of Vaccination"`
	HistoryAllergies                    string `json:"History/Allergies"`
}
