package poit

import (
	"regexp"
	"strings"
	"time"
)

//Permit is a real estate Permit, a special kinf of announcements
type Permit struct {
	URL            string
	AnnouncementID string
	Name           string
	Address        string
	Status         string
	Record         string
	Description    string
	Published      time.Time
	text           []string
	Estate         string
}

//SearchPermits published by a specific municipiality (kommun)
func SearchPermits(municipiality string) chan Permit {

	out := make(chan Permit)
	q := Query{}
	q.selectedPubliceringsIntervall = PubliceringsIntervallLastMonth
	q.selectedAmnesomrade = AmnesomradeRealEstate
	q.kundnamn = municipiality
	q.selectedKundTyp = KundTypKommun

	addressRegexp := regexp.MustCompile(`\([A-ZÅÄÖ\s]+\s[A-Z\d]+\)`) //Address
	recordRegexp := regexp.MustCompile(`SBN\s\d{4}-\d{6}`)           //SBN number
	statusRegexp := regexp.MustCompile(`Bygglov[\s\w\,]+`)           //Status
	estateRegexp := regexp.MustCompile(`^[A-ZÅÄÖ\s\d]+`)             //Estate name
	go func() {
		for a := range SearchWithDetails(q) {

			t := a.Text()
			p := Permit{
				URL:            a.URL(),
				Published:      a.Published(),
				AnnouncementID: a.ID(),
				Name:           a.Name(),
				Status:         statusRegexp.FindString(t[1]),
				Record:         recordRegexp.FindString(t[2]),
				Address:        strings.Trim(addressRegexp.FindString(t[0]), "()"),
				Description:    strings.Join(strings.Split(t[0], ",")[1:], ","),
				Estate:         estateRegexp.FindString(t[0]),
				text:           t,
			}

			out <- p
		}
		close(out)
	}()
	return out
}
