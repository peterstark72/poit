/*
Package poit implements a search query to the POIT website, poit.bolagsverket.se.
*/
package poit

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/publicsuffix"
)

//PoitURL is the start
const PoitURL = "https://poit.bolagsverket.se/poit/PublikPoitIn.do"

//SearchURL is the Form endpoint for searches
const SearchURL = "https://poit.bolagsverket.se/poit/PublikSokKungorelse.do"

//Search returns a channel with announcements
func Search(query string) chan Announcement {

	out := make(chan Announcement)

	go func() {

		var err error

		jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		client := &http.Client{Jar: jar}

		//We do this just to get a cookie
		client.Get(PoitURL)

		//The search query
		data := url.Values{}
		data.Set("selectedPubliceringsIntervall", "2")
		data.Set("fritext", query)
		data.Set("method", "Sök")

		req, _ := http.NewRequest("POST", SearchURL, bytes.NewBufferString(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		for {

			req, _ = http.NewRequest("POST", SearchURL, bytes.NewBufferString(data.Encode()))
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
