//Package poit implements a search query to the POIT website, poit.bolagsverket.se.
package poit

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

//BaseURL is the where POIT is located
const BaseURL = "https://poit.bolagsverket.se"

//PoitURL is where search starts
const PoitURL = "https://poit.bolagsverket.se/poit/PublikPoitIn.do"

//SearchURL is Search App
const SearchURL = "https://poit.bolagsverket.se/poit/PublikSokKungorelse.do"

//MaxNumberOfSessions is the max number of parallell sessions
const MaxNumberOfSessions = 20

//PubliceringsIntervall
const (
	PubliceringsIntervallLastWeek     = "1"
	PubliceringsIntervallLastMonth    = "2"
	PubliceringsIntervallLastQuarter  = "3"
	PubliceringsIntervallLastHalfYear = "5"
	PubliceringsIntervallLastYear     = "6"
)

//Amnesomrade
const (
	AmnesomradeSamtliga        = "-1"
	AmnesomradeKallelser       = "1"
	AmnesomradeBolagsverket    = "2"
	AmnesomradeKonkurser       = "3"
	AmnesomradeFamiljeratt     = "4"
	AmnesomradeSkuldsaneringar = "5"
	AmnesomradeRealEstate      = "8"
)

//KundTyp
const (
	KundTypSamtliga = "-1" //Alla
	KundTypKommun   = "33" //Kommun
)

//Query is a POIT search query
type Query struct {
	selectedPubliceringsIntervall, selectedAmnesomrade, kundnamn, fritext, selectedKundTyp string
}

//AsValues returns query as url.Values
func (q Query) AsValues() url.Values {

	data := url.Values{}
	data.Set("selectedPubliceringsIntervall", q.selectedPubliceringsIntervall)
	data.Set("selectedAmnesomrade", q.selectedAmnesomrade)
	data.Set("kundnamn", q.kundnamn)
	data.Set("fritext", q.fritext)
	data.Set("selectedKundTyp", q.selectedKundTyp)
	data.Set("method", "Sök")

	return data
}

//createSession creates new Poit search session
func createSession(data url.Values) *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client := http.Client{Jar: jar}

	var err error

	//Get the cookie
	_, err = client.Get(PoitURL)
	if err != nil {
		fmt.Println(err)
	}

	//Initiate search app
	req, _ := http.NewRequest("POST", SearchURL, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return &client
}

//hasMorePages returns true if the search result doc has more pages
func hasMorePages(doc *html.Node) bool {
	pageButtons := htmlquery.FindOne(doc, "//em[@class = 'gotopagebuttons']")
	if pageButtons != nil {
		paginationRegexp := regexp.MustCompile(`\d+`)
		breadcrumbs := paginationRegexp.FindAllString(htmlquery.InnerText(pageButtons), 2)
		if len(breadcrumbs) != 2 {
			return false
		}

		curr, _ := strconv.Atoi(breadcrumbs[0])
		tot, _ := strconv.Atoi(breadcrumbs[1])

		return curr < tot
	}

	return false
}

//Announcement is a search result item
type Announcement map[string]string

//ID returns the announcement ID or empty string
func (a Announcement) ID() string {
	return a["id"]
}

//URL returns the URL to announcement text
func (a Announcement) URL() string {
	host, _ := url.Parse(BaseURL)
	u, _ := host.Parse(a["path"])
	return u.String()
}

//Name returns the announcement name or empty string
func (a Announcement) Name() string {
	return a["h-personorgnamn"]
}

//Published returns the published date
func (a Announcement) Published() time.Time {
	d, _ := time.Parse("2006-01-02", a["h-publicerad"])
	return d
}

//Text returns the full text
func (a Announcement) Text() []string {
	return strings.Split(a["kungtext"], "\n")
}

//Search returns a channel with announcements
func Search(q Query) chan Announcement {

	out := make(chan Announcement)

	//Create a POIT search session
	data := q.AsValues()
	client := createSession(data)

	go func() {

		for {
			//Run forever or until no more pages
			req, _ := http.NewRequest("POST", SearchURL, bytes.NewBufferString(data.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}

			doc, err := htmlquery.Parse(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

			for _, row := range htmlquery.Find(doc, "//table[@class = 'result']/tbody/tr") {

				cells := htmlquery.Find(row, "td")

				a := make(map[string]string)
				a["path"] = TrimInnerText(htmlquery.FindOne(cells[0], "/a/@href"))
				a["id"] = TrimInnerText(htmlquery.FindOne(cells[0], "/a/text()"))

				for i := 1; i < len(cells); i++ {
					header := TrimInnerText(htmlquery.FindOne(cells[i], "/@headers"))
					value := TrimInnerText(htmlquery.FindOne(cells[i], "//text()"))
					a[header] = value
				}

				out <- a
			}

			if !hasMorePages(doc) {
				break
			}

			//Sets pagination parameters
			data = url.Values{}
			data.Set("nextFocus", "movenextTop")
			data.Set("scrollPos", "0,0")
			data.Set("method#button.movenext", ">")
		}
		close(out)
	}()
	return out
}

//SearchWithDetails is same as Search but is also loading
//separarete document texts
func SearchWithDetails(q Query) chan Announcement {

	announcements := Search(q)
	out := make(chan Announcement)

	var wg sync.WaitGroup
	worker := func() {

		defer wg.Done()

		client := createSession(q.AsValues())

		for {

			a, ok := <-announcements
			if !ok {
				break
			}

			resp, err := client.Get(a.URL())
			if err != nil {
				fmt.Println(err)
			}

			doc, _ := htmlquery.Parse(resp.Body)
			resp.Body.Close()

			var s []string
			for _, node := range htmlquery.Find(doc, "//div[@class = 'kungtext']//text()") {
				s = append(s, strings.TrimSpace(htmlquery.InnerText(node)))
			}
			a["kungtext"] = strings.Join(s, "\n")
			out <- a
		}
	}

	wg.Add(MaxNumberOfSessions)
	for i := 0; i < MaxNumberOfSessions; i++ {
		go worker()
	}

	//Closer
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
