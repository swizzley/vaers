package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/olivere/elastic/v7"
)

func main() {
	elasticClient, err := elastic.NewSimpleClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}

	for yr := 1990; yr <= time.Now().Year(); yr++ {
		// for yr := 2021; yr <= time.Now().Year(); yr++ {
		loadData(elasticClient, fmt.Sprintf("%v", yr))
	}
	loadData(elasticClient, "NonDomestic")

	defer elasticClient.Stop()

}

func loadData(elasticClient *elastic.Client, yr string) {
	fmt.Println("LOADING...", yr)

	fd, err := os.Open(fmt.Sprintf("../data/%vVAERSDATA.csv", yr))
	if err != nil {
		panic(err)
	}

	events := []VaersEvent{}
	err = gocsv.UnmarshalFile(fd, &events)
	if err != nil {
		panic(err)
	}

	fs, err := os.Open(fmt.Sprintf("../data/%vVAERSSYMPTOMS.csv", yr))
	if err != nil {
		panic(err)
	}

	symp := []Symptoms{}
	err = gocsv.UnmarshalFile(fs, &symp)
	if err != nil {
		panic(err)
	}

	fv, err := os.Open(fmt.Sprintf("../data/%vVAERSVAX.csv", yr))
	if err != nil {
		panic(err)
	}

	vax := []Vax{}
	err = gocsv.UnmarshalFile(fv, &vax)
	if err != nil {
		panic(err)
	}

	for i, e := range events {

		var j []byte
		layout := "01/02/2006"

		if len(e.VAERSID) > 0 {

		S:
			for _, s := range symp {
				if s.VaersID == e.VAERSID {
					e.Symptoms = s
					break S
				}
			}
		V:
			for _, v := range vax {
				if v.VaersID == e.VAERSID {
					e.Vax = v
					break V
				}
			}

			if len(e.RECVDATE) > 0 {
				t, err := time.Parse(layout, e.RECVDATE)
				if err != nil {
					panic(err)
				}
				e.RECVDATETime = t
			}

			if len(e.RPTDATE) > 0 {
				e.RPTDATETime, err = time.Parse(layout, e.RPTDATE)
				if err != nil {
					panic(err)
				}
			}

			if len(e.DATEDIED) > 0 {
				e.DATEDIEDTime, err = time.Parse(layout, e.DATEDIED)
				if err != nil {
					panic(err)
				}
			}

			if len(e.ONSETDATE) > 0 {
				e.ONSETDATETime, err = time.Parse(layout, e.ONSETDATE)
				if err != nil {
					panic(err)
				}
			}

			if len(e.TODAYSDATE) > 0 {
				e.TODAYSDATETime, err = time.Parse(layout, e.TODAYSDATE)
				if err != nil {
					panic(err)
				}
			}

			if len(e.VAXDATE) > 0 {
				e.VAXDATETime, err = time.Parse(layout, e.VAXDATE)
				if err != nil {
					panic(err)
				}
			}

			e.indexTime()

			if e.DATEDIEDTime.After(time.Time{}) ||
				e.RECOVD == "N" ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom1), "death") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom2), "death") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom3), "death") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom4), "death") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom5), "death") ||
				strings.Contains(strings.ToLower(e.SYMPTOMTEXT), "death") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom1), "died") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom2), "died") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom3), "died") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom4), "died") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom5), "died") ||
				strings.Contains(strings.ToLower(e.SYMPTOMTEXT), "died") {
				e.DIED = "Y"
			}
			if strings.Contains(strings.ToLower(e.Symptoms.Symptom1), "abortion") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom2), "abortion") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom3), "abortion") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom4), "abortion") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom5), "abortion") ||
				strings.Contains(strings.ToLower(e.SYMPTOMTEXT), "abortion") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom1), "miscarriage") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom2), "miscarriage") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom3), "miscarriage") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom4), "miscarriage") ||
				strings.Contains(strings.ToLower(e.Symptoms.Symptom5), "miscarriage") ||
				strings.Contains(strings.ToLower(e.SYMPTOMTEXT), "miscarriage") {
				e.BIRTHDEFECT = "Y"
				e.DIED = "Y"
			}

			j, _ = json.Marshal(e)
			body := string(j)

			// fmt.Println(string(j))

			rslt, err := elasticClient.Index().Index("vaers").Type("event").Id(strings.ToLower(e.VAERSID)).BodyJson(body).Do(context.Background())
			if err != nil {
				panic(err)
			}
			fmt.Println(rslt.Result, rslt.Id, i, "of", len(events))

		}
	}

	fd.Close()
	fs.Close()
	fv.Close()
}

func (e *VaersEvent) indexTime() {
	now := time.Now()
	e.IndexTime = time.Now()

	if e.RECVDATETime.Before(e.IndexTime) {
		e.IndexTime = e.RECVDATETime
	}
	if e.RECVDATETime.Before(e.RPTDATETime) {
		e.IndexTime = e.RPTDATETime
	}
	if e.RECVDATETime.Before(e.DATEDIEDTime) {
		e.IndexTime = e.DATEDIEDTime
	}
	if e.RECVDATETime.Before(e.ONSETDATETime) {
		e.IndexTime = e.ONSETDATETime
	}
	if e.RECVDATETime.Before(e.TODAYSDATETime) {
		e.IndexTime = e.TODAYSDATETime
	}
	if e.RECVDATETime.Before(e.VAXDATETime) {
		e.IndexTime = e.VAXDATETime
	}

	if e.IndexTime == now {
		e.IndexTime = time.Date(2021, 1, 20, 1, 1, 1, 1, time.Local)
	}
}

type VaersEvent struct {
	IndexTime      time.Time
	VAERSID        string `csv:"VAERS_ID"`
	RECVDATE       string `csv:"RECVDATE"`
	RECVDATETime   time.Time
	STATE          string `csv:"STATE"`
	AGEYRS         int    `csv:"AGE_YRS"`
	CAGEYR         int    `csv:"CAGE_YR"`
	CAGEMO         string `csv:"CAGE_MO"`
	SEX            string `csv:"SEX"`
	RPTDATE        string `csv:"RPT_DATE"`
	RPTDATETime    time.Time
	SYMPTOMTEXT    string `csv:"SYMPTOM_TEXT"`
	DIED           string `csv:"DIED"`
	DATEDIED       string `csv:"DATEDIED"`
	DATEDIEDTime   time.Time
	LTHREAT        string `csv:"L_THREAT"`
	ERVISIT        string `csv:"ER_VISIT"`
	HOSPITAL       string `csv:"HOSPITAL"`
	HOSPDAYS       int    `csv:"HOSPDAYS"`
	XSTAY          string `csv:"X_STAY"`
	DISABLE        string `csv:"DISABLE"`
	RECOVD         string `csv:"RECOVD"`
	VAXDATE        string `csv:"VAX_DATE"`
	VAXDATETime    time.Time
	ONSETDATE      string `csv:"ONSET_DATE"`
	ONSETDATETime  time.Time
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
	BIRTHDEFECT    string `csv:"BIRTH_DEFECT"`
	OFCVISIT       string `csv:"OFC_VISIT"`
	EREDVISIT      string `csv:"ER_ED_VISIT"`
	ALLERGIES      string `csv:"ALLERGIES"`
	Symptoms       Symptoms
	Vax            Vax
}

type Symptoms struct {
	VaersID         string  `csv:"VAERS_ID"`
	Symptom1        string  `csv:"SYMPTOM1"`
	Symptomversion1 float64 `csv:"SYMPTOMVERSION1"`
	Symptom2        string  `csv:"SYMPTOM2"`
	Symptomversion2 float64 `csv:"SYMPTOMVERSION2"`
	Symptom3        string  `csv:"SYMPTOM3"`
	Symptomversion3 float64 `csv:"SYMPTOMVERSION3"`
	Symptom4        string  `csv:"SYMPTOM4"`
	Symptomversion4 float64 `csv:"SYMPTOMVERSION4"`
	Symptom5        string  `csv:"SYMPTOM5"`
	Symptomversion5 float64 `csv:"SYMPTOMVERSION5"`
}

type Vax struct {
	VaersID       string `csv:"VAERS_ID"`
	VaxType       string `csv:"VAX_TYPE"`
	VaxManu       string `csv:"VAX_MANU"`
	VaxLot        string `csv:"VAX_LOT"`
	VaxDoseSeries string `csv:"VAX_DOSE_SERIES"`
	VaxRoute      string `csv:"VAX_ROUTE"`
	VaxSite       string `csv:"VAX_SITE"`
	VaxName       string `csv:"VAX_NAME"`
}
