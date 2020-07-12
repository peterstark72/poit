package poit

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//RecordNumberRegexp is regexp for a record number
var RecordNumberRegexp = regexp.MustCompile(`SBN\s\d{4}-\d{5}`)

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
