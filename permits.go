package poit

import (
	"regexp"
	"strings"
	"time"
)

var addressRegexp = regexp.MustCompile(`\([A-ZÅÄÖ\s]+\s[A-Z\d]+\)`)   //Address
var recordRegexp = regexp.MustCompile(`SBN\s\d{4}-\d{6}`)             //SBN number
var statusRegexp = regexp.MustCompile(`Bygglov[\s\w\,]+`)             //Status
var estateRegexp = regexp.MustCompile(`^[A-ZÅÄÖ\s\d]+`)               //Estate name
var titleRegexp = regexp.MustCompile(`[A-Z][a-zåäö][\wåäöÅÄÖ()é\s]+`) //Title
var streetRegexp = regexp.MustCompile(`[A-ZÅÄÖ\s]+`)                  //Street name

//Permit is a real estate Permit, a special kinf of announcements
type Permit struct {
	URL            string
	AnnouncementID string
	Name           string
	Address        string
	Status         string
	Record         string
	Title          string
	Published      time.Time
	Text           []string
	Estate         string
	Street         string
}

// NewPermit creates a new Permit from an announcement
func NewPermit(a Announcement) Permit {

	t := a.Text()
	address := strings.Trim(addressRegexp.FindString(t[0]), "()")

	return Permit{
		URL:            a.URL(),
		Published:      a.Published(),
		AnnouncementID: a.ID(),
		Name:           a.Name(),
		Status:         t[1],
		Record:         recordRegexp.FindString(t[2]),
		Address:        address,
		Title:          titleRegexp.FindString(t[0]),
		Estate:         estateRegexp.FindString(t[0]),
		Street:         strings.Title(strings.ToLower(streetRegexp.FindString(address))),
	}

}

//SearchPermits published by a specific municipiality (kommun)
func SearchPermits(query Query) chan Permit {
	out := make(chan Permit)
	go func() {
		for a := range SearchWithDetails(query) {
			out <- NewPermit(a)
		}
		close(out)
	}()
	return out
}
