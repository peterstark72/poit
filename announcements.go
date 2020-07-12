package poit

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//Announcement is an POIT Kungörelse
type Announcement struct {
	ID, Customer, Type, Number, Name, Published string
}

//Permit is a real estate permit
type Permit struct {
	Announcement
	RealEstateName, Address, RecordNumber string
}

//GetDetails returns details of the announcement
func (a Announcement) GetDetails(client *http.Client) []string {

	params := url.Values{}
	params.Set("diarienummer_presentera", a.ID)

	u := fmt.Sprintf("https://poit.bolagsverket.se/poit/PublikSokKungorelse.do?method=presenteraKungorelse&%s", params.Encode())

	resp, err := client.Get(u)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var s []string
	for _, node := range htmlquery.Find(doc, "//div[@class = 'kungtext']//text()") {
		s = append(s, htmlquery.InnerText(node))
		/*
			diarienummer := regexp.MustCompile(`SBN\s\d{4}-\d{5}`)
			pos := diarienummer.FindStringIndex(s)
			if pos != nil {
				fmt.Println(s[pos[0] : pos[1]+1])
			}
		*/
	}
	return s
}

//ParseAnnouncements from resultpage
func ParseAnnouncements(doc *html.Node) []Announcement {

	getcol := func(col *html.Node, expr string) string {
		return strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(col, expr)))
	}

	var announcements []Announcement
	for _, row := range htmlquery.Find(doc, "//table[@class = 'result']/tbody/tr") {

		a := Announcement{
			getcol(row, "td[1]/a/text()"),
			getcol(row, "td[2]"),
			getcol(row, "td[3]"),
			getcol(row, "td[4]"),
			getcol(row, "td[5]"),
			getcol(row, "td[6]"),
		}
		announcements = append(announcements, a)
	}

	return announcements
}

//HasMorePages returns true if there are more pages
func HasMorePages(node *html.Node) bool {
	pageButtons := htmlquery.FindOne(node, "//em[@class = 'gotopagebuttons']")
	if pageButtons == nil {
		return false
	}
	paginationRegexp := regexp.MustCompile(`\d+`)
	if pageButtons != nil {
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
