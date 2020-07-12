package poit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//RecordNumberRegexp is regexp for a record number
var RecordNumberRegexp = regexp.MustCompile(`SBN\s\d{4}-\d{5}`)

//ParseAnnouncementText parses announcement text field
func ParseAnnouncementText(doc *html.Node) AnnouncementText {
	var s []string
	for _, node := range htmlquery.Find(doc, "//div[@class = 'kungtext']//text()") {
		s = append(s, strings.TrimSpace(htmlquery.InnerText(node)))
	}
	return s
}

//ParseAnnouncements from resultpage
func ParseAnnouncements(doc *html.Node) []Announcement {

	var announcements []Announcement
	for _, row := range htmlquery.Find(doc, "//table[@class = 'result']/tbody/tr") {

		var cols []string
		for _, col := range htmlquery.Find(row, "td") {
			cols = append(cols, strings.TrimSpace(htmlquery.InnerText(col)))
		}

		pubDate, err := time.Parse("2006-01-02", cols[5])
		if err != nil {
			fmt.Println("Could not read published date.")
		}
		announcements = append(announcements, Announcement{cols[0], cols[1], cols[2], cols[3], cols[4], pubDate})
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
