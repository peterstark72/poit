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

func TestSearchMany(t *testing.T) {

	q := Query{}
	q.selectedPubliceringsIntervall = PubliceringsIntervallLastMonth
	q.selectedAmnesomrade = AmnesomradeRealEstate
	q.kundnamn = "Malmö kommun"
	q.selectedKundTyp = KundTypKommun

	for a := range Search(q) {
		fmt.Println(a.ID(), a.Name(), a.Published())
	}

}

func TestSearchPermits(t *testing.T) {

	for p := range SearchPermits("Malmö kommun") {
		fmt.Printf("%s / %s: %s[%s] -- %s\n", p.AnnouncementID, p.Record, p.Estate, p.Address, p.Description)
	}

}
