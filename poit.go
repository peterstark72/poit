/*
Package poit implements a search query to the POIT website, poit.bolagsverket.se.

The query is a keyword, e.g. "Malmö", and returns all relevant announcements from
last month. For some reason the POIT search does not allow other time periods.
*/
package poit

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/publicsuffix"
)

//Announcement is an POIT Kungörelse
type Announcement struct {
	ID, Customer, Type, Number, Name, Published string
}

//AnnouncementText is an array of string
type AnnouncementText []string

//PoitURL is the start
const PoitURL = "https://poit.bolagsverket.se/poit/PublikPoitIn.do"

//SearchURL is the Form endpoint for searches
const SearchURL = "https://poit.bolagsverket.se/poit/PublikSokKungorelse.do"

//Client holds the http Client with cookies jar
type Client struct {
	Client *http.Client
}

//NewClient creates a new Poit Client with empty cookie jar
func NewClient() Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return Client{&http.Client{Jar: jar}}
}

//GetAnnouncementText returns details of the announcement
func (poc Client) GetAnnouncementText(ann Announcement) AnnouncementText {

	params := url.Values{}
	params.Set("diarienummer_presentera", ann.ID)
	params.Set("method", "presenteraKungorelse")

	u := fmt.Sprintf("%s?%s", SearchURL, params.Encode())

	resp, err := poc.Client.Get(u)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var s []string
	for _, node := range htmlquery.Find(doc, "//div[@class = 'kungtext']//text()") {
		s = append(s, strings.TrimSpace(htmlquery.InnerText(node)))
	}
	return s
}

//Search returns a channel with announcements
func (poc Client) Search(query string) chan Announcement {

	out := make(chan Announcement)

	go func() {

		var err error

		//We do this just to get a cookie
		_, err = poc.Client.Get(PoitURL)
		if err != nil {
			fmt.Println(err)
		}

		//The search query
		data := url.Values{}
		data.Set("selectedPubliceringsIntervall", "2")
		data.Set("fritext", query)
		data.Set("method", "Sök")

		//For some weird reason, we must make an inital POST and ignore the response
		req, _ := http.NewRequest("POST", SearchURL, bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err = poc.Client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		for {
			//Run forever or until no more pages

			req, _ = http.NewRequest("POST", SearchURL, bytes.NewBufferString(data.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := poc.Client.Do(req)
			if err != nil {
				fmt.Println(err)
			}

			doc, err := htmlquery.Parse(resp.Body)
			if err != nil {
				fmt.Println(err)
			}

			defer resp.Body.Close()

			for _, a := range ParseAnnouncements(doc) {
				out <- a
			}

			if !HasMorePages(doc) {
				break
			}

			//Sets pa
			data = url.Values{}
			data.Set("nextFocus", "movenextTop")
			data.Set("scrollPos", "0,0")
			data.Set("method#button.movenext", ">")
		}
		close(out)
	}()
	return out
}
