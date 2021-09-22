package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gocarina/gocsv"
	"github.com/olivere/elastic/v7"
)

func main() {
	elasticClient, err := elastic.NewSimpleClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}

	f, err := os.Open("../data/20210921VAERSData.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	events := []*VaersEvent{}

	err = gocsv.UnmarshalFile(f, &events)
	if err != nil {
		panic(err)
	}

	fmt.Println("COUNT:", len(events))
	for _, e := range events {
		j, _ := json.Marshal(e)

		if len(e.VAERSID) > 0 {
			layout, _ := dateparse.ParseFormat(e.RECVDATE)
			e.RECVTime, _ = time.Parse(e.RECVDATE, layout)

			layout, _ = dateparse.ParseFormat(e.RPTDATE)
			e.RPTTime, _ = time.Parse(e.RPTDATE, layout)

			layout, _ = dateparse.ParseFormat(e.DATEDIED)
			e.DATEDIEDTime, _ = time.Parse(e.DATEDIED, layout)

			layout, _ = dateparse.ParseFormat(e.ONSETDATE)
			e.ONSETDATETime, _ = time.Parse(e.ONSETDATE, layout)

			layout, _ = dateparse.ParseFormat(e.TODAYSDATE)
			e.TODAYSDATETime, _ = time.Parse(e.TODAYSDATE, layout)

			layout, _ = dateparse.ParseFormat(e.VAXDATE)
			e.VAXDATETime, _ = time.Parse(e.VAXDATE, layout)

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
	VAERSID  string `csv:"VAERS_ID"`
	RECVDATE string `csv:"RECVDATE"`
	RECVTime time.Time

	STATE   string `csv:"STATE"`
	AGEYRS  int    `csv:"AGE_YRS"`
	CAGEYR  int    `csv:"CAGE_YR"`
	CAGEMO  string `csv:"CAGE_MO"`
	SEX     string `csv:"SEX"`
	RPTDATE string `csv:"RPT_DATE"`
	RPTTime time.Time

	SYMPTOMTEXT  string `csv:"SYMPTOM_TEXT"`
	DIED         string `csv:"DIED"`
	DATEDIED     string `csv:"DATEDIED"`
	DATEDIEDTime time.Time

	LTHREAT     string `csv:"L_THREAT"`
	ERVISIT     string `csv:"ER_VISIT"`
	HOSPITAL    string `csv:"HOSPITAL"`
	HOSPDAYS    string `csv:"HOSPDAYS"`
	XSTAY       string `csv:"X_STAY"`
	DISABLE     string `csv:"DISABLE"`
	RECOVD      string `csv:"RECOVD"`
	VAXDATE     string `csv:"VAX_DATE"`
	VAXDATETime time.Time

	ONSETDATE     string `csv:"ONSET_DATE"`
	ONSETDATETime time.Time

	NUMDAYS        int    `csv:"NUMDAYS"`
	LABDATA        string `csv:"LAB_DATA"`
	VADMINBY       string `csv:"V_ADMINBY"`
	VFUNDBY        string `csv:"V_FUNDBY"`
	OTHERMEDS      string `csv:"OTHER_MEDS"`
	CURILL         string `csv:"CUR_ILL"`
	HISTORY        string `csv:"HISTORY"`
	PRIORVAX       string `csv:"PRIOR_VAX"`
	SPLTTYPE       string `csv:"SPLTTYPE"`
	FORMVERS       int    `csv:"FORM_VERS"`
	TODAYSDATE     string `csv:"TODAYS_DATE"`
	TODAYSDATETime time.Time

	BIRTHDEFECT string `csv:"BIRTH_DEFECT"`
	OFCVISIT    string `csv:"OFC_VISIT"`
	EREDVISIT   string `csv:"ER_ED_VISIT"`
	ALLERGIES   string `csv:"ALLERGIES"`
}
