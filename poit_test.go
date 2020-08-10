package poit

import (
	"fmt"
	"testing"
)

func TestSearchWithDetails(t *testing.T) {

	q := Query{}
	q.selectedPubliceringsIntervall = PubliceringsIntervallLastMonth
	q.selectedAmnesomrade = AmnesomradeSamtliga
	q.selectedKundTyp = KundTypSamtliga
	q.fritext = "Tygelsjö"

	for a := range SearchWithDetails(q) {
		fmt.Println(a.ID(), a.Name(), a.Published(), a.Text())
	}

}

var q = Query{
	selectedPubliceringsIntervall: PubliceringsIntervallLastMonth,
	selectedAmnesomrade:           AmnesomradeRealEstate,
	kundnamn:                      "Malmö kommun",
	selectedKundTyp:               KundTypKommun,
}

func TestSearchMany(t *testing.T) {

	for a := range Search(q) {
		fmt.Println(a.ID(), a.Name(), a.Published())
	}

}

func TestSearchPermits(t *testing.T) {
	for p := range SearchPermits(q) {
		//fmt.Printf("%#v", p)
		fmt.Printf("%s på %s\n", p.Title, p.Street)
	}
}
